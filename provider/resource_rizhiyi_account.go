package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-rizhiyi/yottaweb"
)


func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Nickname for the new Account resource.",

			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Email address for the new Account resource.",
			},
			"passwd": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "For the new Account resource,The encryption password for the current encryption algorithm (default is MD5).",
			},
			"full_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Full name of the new Account resource.",
			},

			"group_ids": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "The ID list of user groups to which the new Account resource belongs.",
			},

			"phone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Phone number for the new Account resource.",
			},
			"role_assign_ids": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_ids": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "The ID list of roles assigned to the new Account resource (only admin users can assign).",
			},

			"additional_info": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "Additional information for the new Account resource.",
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	email := d.Get("email").(string)
	passwd := d.Get("passwd").(string)
	full_name := d.Get("full_name").(string)
	group_ids := d.Get("group_ids").(string)
	phone := d.Get("phone").(string)
	role_assign_ids := d.Get("role_assign_ids").(string)
	role_ids := d.Get("role_ids").(string)
	additional_info := d.Get("additional_info").([]interface{})

	requestBody := map[string]interface{}{
		"name":            name,
		"email":           email,
		"passwd":          passwd,
		"full_name":       full_name,
		"group_ids":       group_ids,
		"phone":           phone,
		"role_assign_ids": role_assign_ids,
		"role_ids":        role_ids,
		"additional_info": additional_info,
	}

	endpoint := c.BuildRizhiyiURL(nil, "accounts")
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId(name)
	return nil
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
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

	appID, err := c.GetResourceIdByName(name, "accounts")
	if err != nil {
		return err
	}
	if appID == "" {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	email := d.Get("email").(string)
	passwd := d.Get("passwd").(string)
	full_name := d.Get("full_name").(string)
	group_ids := d.Get("group_ids").(string)
	phone := d.Get("phone").(string)
	role_assign_ids := d.Get("role_assign_ids").(string)
	role_ids := d.Get("role_ids").(string)
	additional_info := d.Get("additional_info").([]interface{})

	requestBody := map[string]interface{}{
		"name":            name,
		"email":           email,
		"passwd":          passwd,
		"full_name":       full_name,
		"group_ids":       group_ids,
		"phone":           phone,
		"role_assign_ids": role_assign_ids,
		"role_ids":        role_ids,
		"additional_info": additional_info,
	}

	update_id, _ := c.GetResourceIdByName(d.Id(), "accounts")
	endpoint := c.BuildRizhiyiURL(nil, "accounts", update_id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId(name)
	return nil
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	del_id, err := c.GetResourceIdByName(name, "accounts")
	if err != nil {
		return err
	}

	endpoint := c.BuildRizhiyiURL(nil, "accounts", del_id)

	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId("")
	return nil
}
