package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-rizhiyi/yottaweb"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"io"
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

	endpoint := c.BuildRizhiyiURL(nil, "v3", "accounts")
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
				if idv, ok := v["object"]; ok {
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
				} else if idv2, ok := v["id"]; ok {
					switch iv := idv2.(type) {
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
		if rid, _ := c.GetResourceIdByName(name, "v3", "accounts"); rid != "" {
			d.SetId(rid)
		}
	}
	return resourceAccountRead(d, m)
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	id := d.Id()
	// 如果当前 state 的 ID 是 name（非纯数字），尝试按 name 解析出数值 id
	if id != "" {
		if _, err := strconv.Atoi(id); err != nil {
			if rid, _ := c.GetResourceIdByName(id, "v3", "accounts"); rid != "" {
				id = rid
				d.SetId(id)
			} else {
				if v, ok := d.GetOk("name"); ok {
					if rid2, _ := c.GetResourceIdByName(v.(string), "v3", "accounts"); rid2 != "" {
						id = rid2
						d.SetId(id)
					}
				}
			}
		}
	}
	if id == "" {
		if v, ok := d.GetOk("name"); ok {
			if rid, _ := c.GetResourceIdByName(v.(string), "v3", "accounts"); rid != "" {
				id = rid
				d.SetId(id)
			}
		}
	}
	if id == "" {
		d.SetId("")
		return nil
	}

	data, err := c.GetResourceById(id, "v3", "accounts")
	if err != nil {
		return err
	}

	d.Set("name", data["name"])
	d.Set("email", data["email"])
	d.Set("full_name", data["full_name"])
	d.Set("group_ids", data["group_ids"])
	d.Set("phone", data["phone"])
	d.Set("role_assign_ids", data["role_assign_ids"])
	roleIDsVal, ok := data["role_ids"]
	if ok && roleIDsVal != nil {
		switch v := roleIDsVal.(type) {
		case string:
			d.Set("role_ids", v)
		case []interface{}:
			var parts []string
			for _, it := range v {
				switch iv := it.(type) {
				case float64:
					parts = append(parts, strconv.Itoa(int(iv)))
				case int:
					parts = append(parts, strconv.Itoa(iv))
				case string:
					parts = append(parts, iv)
				default:
					parts = append(parts, fmt.Sprintf("%v", iv))
				}
			}
			d.Set("role_ids", strings.Join(parts, ","))
		case []string:
			d.Set("role_ids", strings.Join(v, ","))
		case float64:
			d.Set("role_ids", strconv.Itoa(int(v)))
		case int:
			d.Set("role_ids", strconv.Itoa(v))
		default:
			d.Set("role_ids", fmt.Sprintf("%v", v))
		}
	} else {
		if v, ok := d.GetOk("role_ids"); ok {
			d.Set("role_ids", v.(string))
		}
	}
	d.Set("additional_info", data["additional_info"])

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

	update_id := d.Id()
	endpoint := c.BuildRizhiyiURL(nil, "v3", "accounts", update_id)
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	del_id := d.Id()

	endpoint := c.BuildRizhiyiURL(nil, "v3", "accounts", del_id)

	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId("")
	return nil
}
