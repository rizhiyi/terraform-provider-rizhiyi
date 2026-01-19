package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-rizhiyi/yottaweb"
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
	endpoint := c.BuildRizhiyiURL(nil, "roles")
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	app_id, _ := c.GetResourceIdByName(name, "roles")
	d.Set("app_id", app_id)
	d.SetId(name)
	return nil

}

func resourceRolesRead(d *schema.ResourceData, m interface{}) error {
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

	appID, err := c.GetResourceIdByName(name, "roles")
	if err != nil {
		return err
	}
	if appID == "" {
		d.SetId("")
		return nil
	}

	d.Set("app_id", appID)
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

	update_id, _ := c.GetResourceIdByName(d.Id(), "roles")
	endpoint := c.BuildRizhiyiURL(nil, "roles", update_id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId(name)
	return nil
}

func resourceRolesDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	del_id, _ := c.GetResourceIdByName(name, "roles")

	endpoint := c.BuildRizhiyiURL(nil, "roles", del_id)
	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId("")

	return nil
}
