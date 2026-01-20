package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-rizhiyi/yottaweb"
	"encoding/json"
	"io"
	"strconv"
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

			"app_id": &schema.Schema{
				Type:     schema.TypeInt,
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
	app_id := d.Get("app_id").(int)
	rt_names := d.Get("rt_names").(string)
	assign_data := d.Get("assign_data").([]interface{})
	conf := d.Get("conf").(string)
	event_list := d.Get("event_list").([]interface{})

	requestBody := map[string]interface{}{
		"name":        name,
		"logtype":     logtype,
		"enable":      enable,
		"category_id": category_id,
		"app_id":      app_id,
		"rt_names":    rt_names,
		"assign_data": assign_data,
		"conf":        conf,
		"event_list":  event_list,
	}

	endpoint := c.BuildRizhiyiURL(nil, "v3", "parserrules")
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
		if rid, _ := c.GetResourceIdByName(name, "v3", "parserrules"); rid != "" {
			d.SetId(rid)
		}
	}
	return resourceParserRuleRead(d, m)
}

func resourceParserRuleRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	id := d.Id()
	// 如果当前 state 的 ID 是 name（非纯数字），尝试按 name 解析出数值 id
	if id != "" {
		if _, err := strconv.Atoi(id); err != nil {
			if rid, _ := c.GetResourceIdByName(id, "v3", "parserrules"); rid != "" {
				id = rid
				d.SetId(id)
			} else {
				if v, ok := d.GetOk("name"); ok {
					if rid2, _ := c.GetResourceIdByName(v.(string), "v3", "parserrules"); rid2 != "" {
						id = rid2
						d.SetId(id)
					}
				}
			}
		}
	}
	if id == "" {
		if v, ok := d.GetOk("name"); ok {
			if rid, _ := c.GetResourceIdByName(v.(string), "v3", "parserrules"); rid != "" {
				id = rid
				d.SetId(id)
			}
		}
	}
	if id == "" {
		d.SetId("")
		return nil
	}

	data, err := c.GetResourceById(id, "v3", "parserrules")
	if err != nil {
		return err
	}

	d.Set("name", data["name"])
	d.Set("logtype", data["logtype"])
	d.Set("enable", data["enable"])
	d.Set("category_id", data["category_id"])
	d.Set("app_id", data["app_id"])
	d.Set("rt_names", data["rt_names"])
	d.Set("assign_data", data["assign_data"])
	d.Set("conf", data["conf"])
	d.Set("event_list", data["event_list"])

	return nil
}

func resourceParserRuleUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	logtype := d.Get("logtype").(string)
	enable := d.Get("enable").(int)
	category_id := d.Get("category_id").(int)
	app_id := d.Get("app_id").(int)
	rt_names := d.Get("rt_names").(string)
	assign_data := d.Get("assign_data").([]interface{})
	conf := d.Get("conf").(string)
	event_list := d.Get("event_list").([]interface{})

	requestBody := map[string]interface{}{
		"name":        name,
		"logtype":     logtype,
		"enable":      enable,
		"category_id": category_id,
		"app_id":      app_id,
		"rt_names":    rt_names,
		"assign_data": assign_data,
		"conf":        conf,
		"event_list":  event_list,
	}

	update_id := d.Id()
	endpoint := c.BuildRizhiyiURL(nil, "v3", "parserrules", update_id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// keep numeric id stable
	return nil
}

func resourceParserRuleDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	del_id := d.Id()

	endpoint := c.BuildRizhiyiURL(nil, "v3", "parserrules", del_id)

	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId("")
	return nil
}
