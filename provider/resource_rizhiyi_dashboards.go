package provider

import (
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
				Type:     schema.TypeString,
				Required: true,
				Description: "Name for Dashboard resource",
			},
			"rt_names": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Resource group name to which the Dashboard resource belongs. You can pass one or multiple resource tag IDs, separated by commas, for example: test1, test2, test3.",
			},
			"app_ids": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "List of associated apps for the Dashboard resource. You can pass one or multiple application IDs, separated by commas, for example: 5, 12, 32.",
			},
			"data_user": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "viewer",
				Description: "The user's role in accessing the Dashboard，the optional parameters are 'viewer' and 'creator'. (default value viewer)",
			},
		},
	}
}

func resourceDashboardsCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	rt_names := d.Get("rt_names").(string)
	app_ids := d.Get("app_ids").(string)
	data_user := d.Get("data_user").(string)

	requestBody := map[string]interface{}{
		"name":      name,
		"rt_names":  rt_names,
		"app_ids":   app_ids,
		"data_user": data_user,
	}

	endpoint := c.BuildRizhiyiURL(nil, "dashboards")
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId(name)
	return nil
}

func resourceDashboardsRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Id()
	if name == "" {
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
	}
	if name == "" {
		d.SetId("")
		return nil
	}

	appID, err := c.GetResourceIdByName(name, "dashboards")
	if err != nil {
		return err
	}
	if appID == "" {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceDashboardsUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	rt_names := d.Get("rt_names").(string)
	app_ids := d.Get("app_ids").(string)
	data_user := d.Get("data_user").(string)

	requestBody := map[string]interface{}{
		"name":      name,
		"rt_names":  rt_names,
		"app_ids":   app_ids,
		"data_user": data_user,
	}
	update_id, _ := c.GetResourceIdByName(d.Id(), "dashboards")
	endpoint := c.BuildRizhiyiURL(nil, "dashboards", update_id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId(name)
	return nil
}

func resourceDashboardsDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	del_id, err := c.GetResourceIdByName(name, "dashboards")
	if err != nil {
		return err
	}

	endpoint := c.BuildRizhiyiURL(nil, "dashboards", del_id)

	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId("")
	return nil
}
