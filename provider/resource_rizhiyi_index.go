package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/url"
	"terraform-provider-rizhiyi/yottaweb"
	"encoding/json"
	"io"
	"strconv"
)

func resourceIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceIndexesCreate,
		Read:   resourceIndexesRead,
		Update: resourceIndexesUpdate,
		Delete: resourceIndexesDelete,

		Schema: map[string]*schema.Schema{
			"advanced_strategy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},
			"pattern": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Index info mode, value: 'kCompression' for compression mode, default enabling forward optimization/forward compression; 'kNumeric' for numeric mode, default enabling forward optimization/dropping some built-in fields; 'kNormal' for custom mode.",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Index info name",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Index info description",
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Enable status of Index info resource, 0 - enabled, 1 - disabled. (default value 0)",
			},
			"number_of_replicas": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				Description: "",
			},

			"expired_time": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Retention time for Index info resource, for example: 10d.",
			},

			"rotation_period": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Partitioning time for Index info resource, for example: 5d.",
			},
			"sink_to_nas": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},
			"domain_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				Description: "Domain ID for Index info resource, for example: 1. (default value 1)",
			},
			"sink_to_hdd": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},
			"discard_stored_field": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Forward optimization of Index info.",
			},
			"index_name_pattern": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},

			"discard_backup": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"freeze": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},

			"change_disabled_state": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "",
			},
			"use_zstd_compress": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Forward compression of Index info.",
			},
			"reduce_inner_fields": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Dropping some built-in fields of Index info. (default value false)",
			},
			"inject_reduce": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "",
			},
			"tokenizer": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "Indexinfo tokenized fields, currently only standard tokenization (i.e., 'standard' attribute). Fields can be separated by a comma.",
			},
		},
	}
}

func resourceIndexesCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	pattern := d.Get("pattern").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	disabled := d.Get("disabled").(int)
	expired_time := d.Get("expired_time").(string)
	rotation_period := d.Get("rotation_period").(string)
	number_of_replicas := d.Get("number_of_replicas").(int)
	domain_id := d.Get("domain_id").(int)
	index_name_pattern := d.Get("index_name_pattern").(string)
	discard_stored_field := d.Get("discard_stored_field").(string)
	use_zstd_compress := d.Get("use_zstd_compress").(bool)
	reduce_inner_fields := d.Get("reduce_inner_fields").(bool)

	requestBody := map[string]interface{}{
		"pattern":               pattern,
		"name":                  name,
		"description":           description,
		"disabled":              disabled != 0,
		"expired":               expired_time,
		"rotation_period":       rotation_period,
		"number_of_replicas":    number_of_replicas,
		"domain_id":             domain_id,
		"index_name_pattern":    index_name_pattern,
		"discard_stored_field":  discard_stored_field,
		"use_zstd_compress":     use_zstd_compress,
		"reduce_inner_fields":   reduce_inner_fields,
	}

	endpoint := c.BuildRizhiyiURL(nil, "..", "v3", "indexes")
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
		if rid, _ := c.GetResourceIdByName(name, "..", "v3", "indexes"); rid != "" {
			d.SetId(rid)
		}
	}
	return nil
}

func resourceIndexesRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	id := d.Id()
	// 如果当前 state 的 ID 是 name（非纯数字），尝试按 name 解析出数值 id
	if id != "" {
		if _, err := strconv.Atoi(id); err != nil {
			if rid, _ := c.GetResourceIdByName(id, "..", "v3", "indexes"); rid != "" {
				id = rid
				d.SetId(id)
			} else {
				if v, ok := d.GetOk("name"); ok {
					if rid2, _ := c.GetResourceIdByName(v.(string), "..", "v3", "indexes"); rid2 != "" {
						id = rid2
						d.SetId(id)
					}
				}
			}
		}
	}
	if id == "" {
		if v, ok := d.GetOk("name"); ok {
			if rid, _ := c.GetResourceIdByName(v.(string), "..", "v3", "indexes"); rid != "" {
				id = rid
				d.SetId(id)
			}
		}
	}
	if id == "" {
		d.SetId("")
		return nil
	}

	data, err := c.GetResourceById(id, "..", "v3", "indexes")
	if err != nil {
		return err
	}

	d.Set("name", data["name"])
	d.Set("description", data["description"])
	if disabled, ok := data["disabled"].(bool); ok {
		if disabled {
			d.Set("disabled", 1)
		} else {
			d.Set("disabled", 0)
		}
	}
	d.Set("rotation_period", data["rotation_period"])
	d.Set("expired_time", data["expired"])
	d.Set("pattern", data["pattern"])
	d.Set("domain_id", data["domain_id"])
	d.Set("number_of_replicas", data["number_of_replicas"])
	d.Set("index_name_pattern", data["index_name_pattern"])
	d.Set("discard_stored_field", data["discard_stored_field"])
	d.Set("use_zstd_compress", data["use_zstd_compress"])
	d.Set("reduce_inner_fields", data["reduce_inner_fields"])
	d.Set("freeze", data["freeze"])
	d.Set("sink_to_nas", data["sink_to_nas"])
	d.Set("sink_to_hdd", data["sink_to_hdd"])
	d.Set("discard_backup", data["discard_backup"])

	return nil
}

func resourceIndexesUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	pattern := d.Get("pattern").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	disabled := d.Get("disabled").(int)
	expired_time := d.Get("expired_time").(string)
	rotation_period := d.Get("rotation_period").(string)
	number_of_replicas := d.Get("number_of_replicas").(int)
	domain_id := d.Get("domain_id").(int)
	index_name_pattern := d.Get("index_name_pattern").(string)
	discard_stored_field := d.Get("discard_stored_field").(string)
	use_zstd_compress := d.Get("use_zstd_compress").(bool)
	reduce_inner_fields := d.Get("reduce_inner_fields").(bool)

	requestBody := map[string]interface{}{
		"pattern":               pattern,
		"name":                  name,
		"description":           description,
		"disabled":              disabled != 0,
		"expired":               expired_time,
		"rotation_period":       rotation_period,
		"number_of_replicas":    number_of_replicas,
		"domain_id":             domain_id,
		"index_name_pattern":    index_name_pattern,
		"discard_stored_field":  discard_stored_field,
		"use_zstd_compress":     use_zstd_compress,
		"reduce_inner_fields":   reduce_inner_fields,
	}

	update_id := d.Id()
	endpoint := c.BuildRizhiyiURL(nil, "..", "v3", "indexes", update_id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func resourceIndexesDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	del_id := d.Id()

	// build delete index parameters
	parametersValues := url.Values{}
	parametersValues.Add("engine", "beaver")
	parametersValues.Add("index_name", name)
	endpoint := c.BuildRizhiyiURL(parametersValues, "..", "v3", "indexes", del_id)

	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId("")

	return nil
}
