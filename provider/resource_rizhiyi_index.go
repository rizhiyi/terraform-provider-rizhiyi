package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/url"
	"terraform-provider-rizhiyi/yottaweb"
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
	advanced_strategy := d.Get("advanced_strategy").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	disabled := d.Get("disabled").(int)
	expired_time := d.Get("expired_time").(string)
	freeze := d.Get("freeze").(string)
	rotation_period := d.Get("rotation_period").(string)
	discard_backup := d.Get("discard_backup").(string)
	sink_to_nas := d.Get("sink_to_nas").(string)
	number_of_replicas := d.Get("number_of_replicas").(int)
	domain_id := d.Get("domain_id").(int)
	sink_to_hdd := d.Get("sink_to_hdd").(string)
	index_name_pattern := d.Get("index_name_pattern").(string)
	discard_stored_field := d.Get("discard_stored_field").(string)
	use_zstd_compress := d.Get("use_zstd_compress").(bool)
	change_disabled_state := d.Get("change_disabled_state").(bool)
	reduce_inner_fields := d.Get("reduce_inner_fields").(bool)
	inject_reduce := d.Get("inject_reduce").(map[string]interface{})
	tokenizer := d.Get("tokenizer").(map[string]interface{})

	requestBody := map[string]interface{}{
		"pattern":               pattern,
		"advanced_strategy":     advanced_strategy,
		"name":                  name,
		"description":           description,
		"disabled":              disabled,
		"expired_time":          expired_time,
		"discard_backup":        discard_backup,
		"freeze":                freeze,
		"rotation_period":       rotation_period,
		"index_name_pattern":    index_name_pattern,
		"sink_to_nas":           sink_to_nas,
		"domain_id":             domain_id,
		"sink_to_hdd":           sink_to_hdd,
		"discard_stored_field":  discard_stored_field,
		"use_zstd_compress":     use_zstd_compress,
		"number_of_replicas":    number_of_replicas,
		"change_disabled_state": change_disabled_state,
		"reduce_inner_fields":   reduce_inner_fields,
		"inject_reduce":         inject_reduce,
		"tokenizer":             tokenizer,
	}

	endpoint := c.BuildRizhiyiURL(nil, "indexinfo")
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId(name)
	return nil
}

func resourceIndexesRead(d *schema.ResourceData, m interface{}) error {
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

	appID, err := c.GetResourceIdByName(name, "indexinfo")
	if err != nil {
		return err
	}
	if appID == "" {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceIndexesUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	pattern := d.Get("pattern").(string)
	advanced_strategy := d.Get("advanced_strategy").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	disabled := d.Get("disabled").(int)
	expired_time := d.Get("expired_time").(string)
	freeze := d.Get("freeze").(string)
	rotation_period := d.Get("rotation_period").(string)
	discard_backup := d.Get("discard_backup").(string)
	sink_to_nas := d.Get("sink_to_nas").(string)
	number_of_replicas := d.Get("number_of_replicas").(int)
	domain_id := d.Get("domain_id").(int)
	sink_to_hdd := d.Get("sink_to_hdd").(string)
	index_name_pattern := d.Get("index_name_pattern").(string)
	discard_stored_field := d.Get("discard_stored_field").(string)
	use_zstd_compress := d.Get("use_zstd_compress").(bool)
	change_disabled_state := d.Get("change_disabled_state").(bool)
	reduce_inner_fields := d.Get("reduce_inner_fields").(bool)
	inject_reduce := d.Get("inject_reduce").(map[string]interface{})
	tokenizer := d.Get("tokenizer").(map[string]interface{})

	requestBody := map[string]interface{}{
		"pattern":               pattern,
		"advanced_strategy":     advanced_strategy,
		"name":                  name,
		"description":           description,
		"disabled":              disabled,
		"expired_time":          expired_time,
		"discard_backup":        discard_backup,
		"freeze":                freeze,
		"rotation_period":       rotation_period,
		"index_name_pattern":    index_name_pattern,
		"sink_to_nas":           sink_to_nas,
		"domain_id":             domain_id,
		"sink_to_hdd":           sink_to_hdd,
		"discard_stored_field":  discard_stored_field,
		"use_zstd_compress":     use_zstd_compress,
		"number_of_replicas":    number_of_replicas,
		"change_disabled_state": change_disabled_state,
		"reduce_inner_fields":   reduce_inner_fields,
		"inject_reduce":         inject_reduce,
		"tokenizer":             tokenizer,
	}

	update_id, _ := c.GetResourceIdByName(d.Id(), "indexinfo")
	endpoint := c.BuildRizhiyiURL(nil, "indexinfo", update_id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId(name)
	return nil
}

func resourceIndexesDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	del_id, err := c.GetResourceIdByName(name, "indexinfo")
	if err != nil {
		return err
	}

	// build delete index parameters
	parametersValues := url.Values{}
	parametersValues.Add("engine", "beaver")
	parametersValues.Add("id", del_id)
	parametersValues.Add("name", name)
	parametersValues.Add("index_name", name)
	endpoint := c.BuildRizhiyiURL(parametersValues, "indexinfo", del_id)

	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId("")

	return nil
}
