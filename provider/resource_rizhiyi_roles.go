package provider

import (
	"encoding/json"
	"io"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-rizhiyi/yottaweb"
	"strconv"
)

func resourceRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceRolesCreate,
		Read:   resourceRolesRead,
		Update: resourceRolesUpdate,
		Delete: resourceRolesDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Nickname for the new Role resource.",
			},
			"memo": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Resource description for the new Role resource.",
			},
			"app_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceRolesCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	memo := d.Get("memo").(string)

	requestBody := map[string]interface{}{
		"name": name,
		"memo": memo,
	}
	// add request path parameters
	endpoint := c.BuildRizhiyiURL(nil, "v3", "roles")
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var respData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &respData); err == nil {
		if obj, ok := respData["object"]; ok {
			switch v := obj.(type) {
			case float64:
				d.SetId(strconv.Itoa(int(v)))
			case map[string]interface{}:
				if idv, ok := v["id"]; ok {
					switch iv := idv.(type) {
					case float64:
						d.SetId(strconv.Itoa(int(iv)))
					case string:
						d.SetId(iv)
					case int:
						d.SetId(strconv.Itoa(iv))
					default:
						d.SetId("")
					}
				}
			case string:
				d.SetId(v)
			}
		}
	}
	if d.Id() == "" {
		if rid, _ := c.GetResourceIdByName(name, "v3", "roles"); rid != "" {
			d.SetId(rid)
		}
	}
	return resourceRolesRead(d, m)

}

func resourceRolesRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	id := d.Id()
	// 如果当前 state 的 ID 是 name（非纯数字），尝试按 name 解析出数值 id
	if id != "" {
		if _, err := strconv.Atoi(id); err != nil {
			if rid, _ := c.GetResourceIdByName(id, "v3", "roles"); rid != "" {
				id = rid
				d.SetId(id)
			} else {
				// 再尝试从属性 name 获取
				if v, ok := d.GetOk("name"); ok {
					if rid2, _ := c.GetResourceIdByName(v.(string), "v3", "roles"); rid2 != "" {
						id = rid2
						d.SetId(id)
					}
				}
			}
		}
	}
	if id == "" {
		if v, ok := d.GetOk("name"); ok {
			if rid, _ := c.GetResourceIdByName(v.(string), "v3", "roles"); rid != "" {
				id = rid
				d.SetId(id)
			}
		}
	}
	if id == "" {
		d.SetId("")
		return nil
	}

	data, err := c.GetResourceById(id, "v3", "roles")
	if err != nil {
		return err
	}

	d.Set("name", data["name"])
	d.Set("memo", data["memo"])
	return nil
}

func resourceRolesUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	memo := d.Get("memo").(string)

	requestBody := map[string]interface{}{
		"name": name,
		"memo": memo,
	}

	update_id := d.Id()
	endpoint := c.BuildRizhiyiURL(nil, "v3", "roles", update_id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func resourceRolesDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	del_id := d.Id()

	endpoint := c.BuildRizhiyiURL(nil, "v3", "roles", del_id)
	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId("")

	return nil
}
