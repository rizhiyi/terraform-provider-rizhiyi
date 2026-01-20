package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-rizhiyi/yottaweb"
)

func resourceDashboards() *schema.Resource {
	return &schema.Resource{
		Create: resourceDashboardsCreate,
		Read:   resourceDashboardsRead,
		Update: resourceDashboardsUpdate,
		Delete: resourceDashboardsDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name for Dashboard resource",
			},
			"rt_names": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource group name to which the Dashboard resource belongs.",
			},
			"app_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Associated app ID for the Dashboard resource.",
			},
			"data_user": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "viewer",
				Description: "The user's role in accessing the Dashboard，the optional parameters are 'viewer' and 'creator'. (default value viewer)",
			},
			"export": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "local",
				Description: "Resource scope: local (visible within the app) or system (globally visible).",
			},
			"default_display": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Default display setting.",
			},
			"sequences": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sequences configuration.",
			},
			"active_tab": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the active tab.",
			},
			"manage_tabs": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether Terraform should manage dashboard tabs. If false, tabs are read-only and ignored in diffs.",
			},
			"tabs": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the tab.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the tab.",
						},
						"content": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							DiffSuppressFunc: suppressEquivalentJSON,
							Description: "Content of the tab (JSON string).",
						},
						"uuid": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UUID of the tab.",
						},
						"creator_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creator ID of the tab.",
						},
					},
				},
			},
		},
	}
}

func suppressEquivalentJSON(k, old, new string, d *schema.ResourceData) bool {
	var jo interface{}
	var jn interface{}
	_ = json.Unmarshal([]byte(old), &jo)
	_ = json.Unmarshal([]byte(new), &jn)
	return reflect.DeepEqual(jo, jn)
}

func toInt(v interface{}) (int, bool) {
	if v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	case string:
		if t == "" {
			return 0, false
		}
		if i, e := strconv.Atoi(t); e == nil {
			return i, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func resourceDashboardsCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	
	requestBody := map[string]interface{}{
		"name":            name,
		"rt_names":        d.Get("rt_names").(string),
		"app_id":          d.Get("app_id").(int),
		"data_user":       d.Get("data_user").(string),
		"export":          d.Get("export").(string),
		"default_display": d.Get("default_display").(int),
		"sequences":       d.Get("sequences").(string),
		"active_tab":      d.Get("active_tab").(int),
	}

	// Use v3 API
	endpoint := c.BuildRizhiyiURL(nil, "..", "v3", "dashboards")
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response to get ID
	bodyBytes, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return fmt.Errorf("failed to unmarshal create response: %s", err)
	}

	var dashboardIDStr string
	if obj, ok := result["object"].(map[string]interface{}); ok {
		if idVal, ok := obj["id"].(float64); ok {
			dashboardIDStr = strconv.Itoa(int(idVal))
		}
	}

	if dashboardIDStr == "" {
		// Fallback to GetResourceIdByName if ID not found in response
		// Note: GetResourceIdByName uses v2, might fail if v3 items are different, but usually works.
		// However, let's try to avoid it if possible. If we are here, it means response format was unexpected.
		// Let's log or error out? Or try fallback.
		// Fallback:
		var errFallback error
		dashboardIDStr, errFallback = c.GetResourceIdByName(name, "..", "v3", "dashboards")
		if errFallback != nil {
			return fmt.Errorf("failed to get dashboard ID from response and fallback: %s (response: %s)", errFallback, string(bodyBytes))
		}
	}
	
	d.SetId(dashboardIDStr)

	// Create Tabs
	if d.Get("manage_tabs").(bool) {
		if v, ok := d.GetOk("tabs"); ok {
		tabs := v.([]interface{})
		for _, t := range tabs {
			tab := t.(map[string]interface{})
			tabName := tab["name"].(string)
			tabContent := tab["content"].(string)

			tabBody := map[string]interface{}{
				"name":    tabName,
				"content": tabContent,
			}
			// POST /api/v3/dashboards/{dashboard_id}/tabs/
			tabEndpoint := c.BuildRizhiyiURL(nil, "..", "v3", "dashboards", dashboardIDStr, "tabs")
			tabResp, err := c.Post(tabEndpoint, tabBody)
			if err != nil {
				return fmt.Errorf("failed to create tab %s: %s", tabName, err)
			}
			tabResp.Body.Close()
		}
		}
	}

	return resourceDashboardsRead(d, m)
}

func resourceDashboardsRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	id := d.Id()
	if id == "" {
		return nil
	}

	// GET /api/v3/dashboards/{id}/
	endpoint := c.BuildRizhiyiURL(nil, "..", "v3", "dashboards", id)
	resp, err := c.Get(endpoint)
	if err != nil {
		// If 404, remove from state
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return err
	}

	// Check for API level error in 200 OK response
	if res, ok := result["result"].(bool); ok && !res {
		if errObj, ok := result["error"].(map[string]interface{}); ok {
			msg, _ := errObj["message"].(string)
			code, _ := errObj["code"].(string)
			// If permission denied (1604) or not found, we might consider it gone if the ID looks wrong?
			// But for now, return a clear error.
			return fmt.Errorf("API error reading dashboard %s: code=%v, message=%v", id, code, msg)
		}
		return fmt.Errorf("API returned result: false but no error detail for dashboard %s", id)
	}

	var dashboard map[string]interface{}
	if obj, ok := result["object"].(map[string]interface{}); ok {
		dashboard = obj
	} else {
		// If structure is flat? Unlikely based on schema.
		// If not found?
		return fmt.Errorf("unexpected response format for dashboard read: %s", string(bodyBytes))
	}

	d.Set("name", dashboard["name"])
	if v, ok := dashboard["rt_names"]; ok {
		d.Set("rt_names", v)
	}
	if v, ok := dashboard["app_id"]; ok {
		if iv, ok2 := toInt(v); ok2 {
			d.Set("app_id", iv)
		}
	}
	if v, ok := dashboard["data_user"]; ok {
		d.Set("data_user", v)
	}
	if v, ok := dashboard["export"]; ok {
		d.Set("export", v)
	}
	if v, ok := dashboard["default_display"]; ok {
		if iv, ok2 := toInt(v); ok2 {
			d.Set("default_display", iv)
		}
	}
	if v, ok := dashboard["sequences"]; ok {
		d.Set("sequences", v)
	}
	if v, ok := dashboard["active_tab"]; ok {
		if iv, ok2 := toInt(v); ok2 {
			d.Set("active_tab", iv)
		}
	}

	// Process Tabs from dashboard object
	var tabsList []interface{}
	if v, ok := dashboard["tabs"].([]interface{}); ok {
		tabsList = v
	}

	if d.Get("manage_tabs").(bool) {
		tfTabs := make([]map[string]interface{}, 0, len(tabsList))
		for _, t := range tabsList {
			tabMap := t.(map[string]interface{})
			tfTab := map[string]interface{}{
				"id":         int(tabMap["id"].(float64)),
				"name":       tabMap["name"].(string),
				"content":    tabMap["content"].(string),
				"uuid":       tabMap["uuid"].(string),
				"creator_id": int(tabMap["creator_id"].(float64)),
			}
			tfTabs = append(tfTabs, tfTab)
		}
		d.Set("tabs", tfTabs)
	}

	return nil
}

func resourceDashboardsUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	id := d.Id()
	
	requestBody := map[string]interface{}{
		"name":            d.Get("name").(string),
		"rt_names":        d.Get("rt_names").(string),
		"app_id":          d.Get("app_id").(int),
		"data_user":       d.Get("data_user").(string),
		"export":          d.Get("export").(string),
		"default_display": d.Get("default_display").(int),
		"sequences":       d.Get("sequences").(string),
		"active_tab":      d.Get("active_tab").(int),
	}

	endpoint := c.BuildRizhiyiURL(nil, "..", "v3", "dashboards", id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}
	resp.Body.Close()

	// Update Tabs
	if d.Get("manage_tabs").(bool) && d.HasChange("tabs") {
		// 1. Get current state (we can re-read or trust Read was called)
		// Better to fetch current tabs to be safe
		readEndpoint := c.BuildRizhiyiURL(nil, "..", "v3", "dashboards", id)
		readResp, err := c.Get(readEndpoint)
		if err != nil {
			return err
		}
		defer readResp.Body.Close()
		
		readBytes, _ := io.ReadAll(readResp.Body)
		var readResult map[string]interface{}
		json.Unmarshal(readBytes, &readResult)
		
		type tabInfo struct {
			id   int
			name string
		}
		existingTabs := make([]tabInfo, 0)
		if obj, ok := readResult["object"].(map[string]interface{}); ok {
			if tabs, ok := obj["tabs"].([]interface{}); ok {
				for _, t := range tabs {
					tm := t.(map[string]interface{})
					existingTabs = append(existingTabs, tabInfo{
						id:   int(tm["id"].(float64)),
						name: tm["name"].(string),
					})
				}
			}
		}

		// 2. Process desired tabs
		newTabs := d.Get("tabs").([]interface{})
		matchedExistingIDs := make(map[int]bool)

		for _, t := range newTabs {
			tab := t.(map[string]interface{})
			name := tab["name"].(string)
			content := tab["content"].(string)
			desiredID := 0
			if v, ok := tab["id"]; ok {
				switch dv := v.(type) {
				case int:
					desiredID = dv
				case float64:
					desiredID = int(dv)
				}
			}

			tabBody := map[string]interface{}{
				"name":    name,
				"content": content,
			}

			var targetID int
			if desiredID > 0 {
				targetID = desiredID
			} else {
				// find first unmatched existing by name
				for _, et := range existingTabs {
					if et.name == name && !matchedExistingIDs[et.id] {
						targetID = et.id
						break
					}
				}
			}

			if targetID > 0 {
				// Update
				// PUT /api/v3/dashboards/{did}/tabs/{tid}/
				tabIDStr := strconv.Itoa(targetID)
				updEndpoint := c.BuildRizhiyiURL(nil, "..", "v3", "dashboards", id, "tabs", tabIDStr)
				if r, e := c.Put(updEndpoint, tabBody); e == nil {
					r.Body.Close()
				} else {
					return fmt.Errorf("failed to update tab %s: %s", name, e)
				}
				matchedExistingIDs[targetID] = true
			} else {
				// Create
				createEndpoint := c.BuildRizhiyiURL(nil, "..", "v3", "dashboards", id, "tabs")
				if r, e := c.Post(createEndpoint, tabBody); e == nil {
					r.Body.Close()
				} else {
					return fmt.Errorf("failed to create tab %s: %s", name, e)
				}
			}
		}

		// 3. Delete removed tabs
		for _, et := range existingTabs {
			if !matchedExistingIDs[et.id] {
				// Delete
				tabIDStr := strconv.Itoa(et.id)
				delEndpoint := c.BuildRizhiyiURL(nil, "..", "v3", "dashboards", id, "tabs", tabIDStr)
				if r, e := c.Delete(delEndpoint); e == nil {
					r.Body.Close()
				} else {
					return fmt.Errorf("failed to delete tab id=%d name=%s: %s", et.id, et.name, e)
				}
			}
		}
	}

	return resourceDashboardsRead(d, m)
}

func resourceDashboardsDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	id := d.Id()

	endpoint := c.BuildRizhiyiURL(nil, "..", "v3", "dashboards", id)
	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}
	resp.Body.Close()
	d.SetId("")
	return nil
}
