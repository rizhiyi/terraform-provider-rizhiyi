package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-rizhiyi/yottaweb"
)


func resourceParserRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceParserRuleCreate,
		Read:   resourceParserRuleRead,
		Update: resourceParserRuleUpdate,
		Delete: resourceParserRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "ParserRule resource name.",
			},
			"logtype": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "ParserRule log type field.for example json,apache",
			},
			"enable": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "ParserRule enable status, 0 - enabled, 1 - disabled. (default value 0)",
			},
			"category_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 1000,
				Description: "ParserRule ownership type, determines whether it is a system default rule. User-created rules are all assigned a value of 1000. (default value 1000)",
			},

			"app_ids": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "App ID to which the ParserRule resource belongs.",
			},

			"rt_names": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Resource group name to which the ParserRule resource belongs.",
			},

			"assign_data": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "ParserRule appname & tag",
			},
			"conf": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Parsing rules included in the ParserRule. for examples: (Parsing rules for JSON) \"[{\"json\":{\"rule\":[{\"add_fields\":[],\"source\":\"raw_message\",\"another_name\":\"\",\"paths\":[],\"extract_limit\":\"\"}]}}]\"",
			},
			"event_list": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "",
			},
		},
	}
}

func resourceParserRuleCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	logtype := d.Get("logtype").(string)
	enable := d.Get("enable").(int)
	category_id := d.Get("category_id").(int)
	app_ids := d.Get("app_ids").(string)
	rt_names := d.Get("rt_names").(string)
	assign_data := d.Get("assign_data").([]interface{})
	conf := d.Get("conf").(string)
	event_list := d.Get("event_list").([]interface{})

	requestBody := map[string]interface{}{
		"name":        name,
		"logtype":     logtype,
		"enable":      enable,
		"category_id": category_id,
		"app_ids":     app_ids,
		"rt_names":    rt_names,
		"assign_data": assign_data,
		"conf":        conf,
		"event_list":  event_list,
	}

	endpoint := c.BuildRizhiyiURL(nil, "parserrules")
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId(name)
	return nil
}

func resourceParserRuleRead(d *schema.ResourceData, m interface{}) error {
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

	appID, err := c.GetResourceIdByName(name, "parserrules")
	if err != nil {
		return err
	}
	if appID == "" {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceParserRuleUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	logtype := d.Get("logtype").(string)
	enable := d.Get("enable").(int)
	category_id := d.Get("category_id").(int)
	app_ids := d.Get("app_ids").(string)
	rt_names := d.Get("rt_names").(string)
	assign_data := d.Get("assign_data").([]interface{})
	conf := d.Get("conf").(string)
	event_list := d.Get("event_list").([]interface{})

	requestBody := map[string]interface{}{
		"name":        name,
		"logtype":     logtype,
		"enable":      enable,
		"category_id": category_id,
		"app_ids":     app_ids,
		"rt_names":    rt_names,
		"assign_data": assign_data,
		"conf":        conf,
		"event_list":  event_list,
	}

	update_id, _ := c.GetResourceIdByName(d.Id(), "parserrules")
	endpoint := c.BuildRizhiyiURL(nil, "parserrules", update_id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId(name)
	return nil
}

func resourceParserRuleDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	del_id, err := c.GetResourceIdByName(name, "parserrules")
	if err != nil {
		return err
	}

	endpoint := c.BuildRizhiyiURL(nil, "parserrules", del_id)

	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId("")
	return nil
}
