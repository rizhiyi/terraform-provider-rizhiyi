package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-rizhiyi/yottaweb"
)

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertCreate,
		Read:   resourceAlertRead,
		Update: resourceAlertUpdate,
		Delete: resourceAlertDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Resource name for the new Alert resource.",
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: " Description of the new Alert resource.",
			},
			"category": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				Description: "The monitoring types for the new Alert resource are as follows: 0.Event Count Monitoring, 1.Field Statistics Monitoring, 2.Continuous Statistics Monitoring, 3.Baseline Comparison Monitoring, 4.SPL Statistics Monitoring. (default value 0)",
			},
			"check_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The scheduled execution plan for the new Alert resource, fill in the interval in seconds for the scheduled execution plan, where 0 indicates no scheduled execution plan. (default value 0)",
			},
			"continuous_trigger_value": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "",
			},
			"executor_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "User ID for executing the new Alert resource.（default value 0）",
			},
			"market_day": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Whether the new Alert resource is executed only on the transaction day. (default value false)",
			},
			"max_restrain_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The doubling time (in seconds) for suppressing the cancellation of alert monitoring. (default value 0)",
			},
			"restrain_interval": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Monitoring suppression time (in seconds) for the alert resource. (default value 0)",
			},
			"schedule_priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "",
			},
			"query": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Search content for the alert resource.",
			},
			"check_condition_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"extend_query": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Search content for the extended search of the alert resource.",
			},
			"crontab": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "The crontab execution schedule for the alert resource, please provide the corresponding cron statement, for example, 0 * * * * ？, where 0 indicates not using the crontab execution schedule.",
			},

			"app_ids": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "List of application IDs associated with the Alert resource, for example: 1, 2, 3.",
			},
			"group_suppress_field": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},
			"rt_names": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Resource group name to which the Alert resource belongs, for example: default_Alert, test.",
			},
			"schedule_window": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},
			"segmentation_field": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Device split field for the Alert resource, an empty string indicates no device split.",
			},
			"statistics_field": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},

			"topic": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},
			"window": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "",
			},

			"alert_when_recover": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Whether the Alert resource uses monitoring reply prompts.(default value false)",
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "The field to enable monitoring for the Alert resource. (default value false)",
			},
			"extend_use_spark": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Whether the extended search of the Alert resource uses advanced mode. (default value false)",
			},
			"graph_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Whether the extended search of the Alert resource has the effect illustration enabled. (default value false)",
			},
			"group_trigger_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "",
			},
			"hosted_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "",
			},
			"use_spark": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Whether the Alert resource uses advanced mode. (default value false)",
			},
			"alert_metas": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "Plugin data for the Alert resource (in JSON array string format), where each item in the array should provide the plugin's name, trigger level, configuration information, and change data.",
			},
			"dataset_ids": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "JSON string for the dataset node ID of the Alert resource.",
			},
			"extend_dataset_ids": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "JSON string for the dataset node ID of the extended search in the Alert resource.",
			},

			"check_condition": &schema.Schema{
				Type: schema.TypeString,
				Required:    true,
				Description: "Monitoring trigger conditions for the Alert resource.",
			},
			"extend_conf": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
				Description: "Fixed key-value for the extended search of the Alert resource.",
			},
		},
	}
}

func resourceAlertCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	category := d.Get("category").(int)
	check_interval := d.Get("check_interval").(int)
	continuous_trigger_value := d.Get("continuous_trigger_value").(int)
	market_day := d.Get("market_day").(bool)
	max_restrain_interval := d.Get("max_restrain_interval").(int)
	restrain_interval := d.Get("restrain_interval").(int)
	schedule_priority := d.Get("schedule_priority").(int)
	executor_id := d.Get("executor_id").(int)
	query := d.Get("query").(string)
	check_condition_group := d.Get("check_condition_group").(string)
	extend_query := d.Get("extend_query").(string)
	crontab := d.Get("crontab").(string)
	app_ids := d.Get("app_ids").(string)
	group_suppress_field := d.Get("group_suppress_field").(string)
	rt_names := d.Get("rt_names").(string)
	schedule_window := d.Get("schedule_window").(string)
	segmentation_field := d.Get("segmentation_field").(string)
	statistics_field := d.Get("statistics_field").(string)
	topic := d.Get("topic").(string)
	window := d.Get("window").(string)
	alert_when_recover := d.Get("alert_when_recover").(bool)
	enabled := d.Get("enabled").(bool)
	extend_use_spark := d.Get("extend_use_spark").(bool)
	graph_enabled := d.Get("graph_enabled").(bool)
	group_trigger_flag := d.Get("group_trigger_flag").(bool)
	hosted_flag := d.Get("hosted_flag").(bool)
	use_spark := d.Get("use_spark").(bool)
	alert_metas := d.Get("alert_metas").([]interface{})
	dataset_ids := d.Get("dataset_ids").([]interface{})
	extend_dataset_ids := d.Get("extend_dataset_ids").([]interface{})
	check_condition := d.Get("check_condition").(string)
	extend_conf := d.Get("extend_conf").(string)

	requestBody := map[string]interface{}{
		"name":                     name,
		"description":              description,
		"category":                 category,
		"check_interval":           check_interval,
		"continuous_trigger_value": continuous_trigger_value,
		"market_day":               market_day,
		"max_restrain_interval":    max_restrain_interval,
		"restrain_interval":        restrain_interval,
		"schedule_priority":        schedule_priority,
		"query":                    query,
		"check_condition_group":    check_condition_group,
		"extend_query":             extend_query,
		"crontab":                  crontab,
		"app_ids":                  app_ids,
		"group_suppress_field":     group_suppress_field,
		"rt_names":                 rt_names,
		"schedule_window":          schedule_window,
		"executor_id":              executor_id,
		"segmentation_field":       segmentation_field,
		"statistics_field":         statistics_field,
		"topic":                    topic,
		"window":                   window,
		"alert_when_recover":       alert_when_recover,
		"enabled":                  enabled,
		"extend_use_spark":         extend_use_spark,
		"graph_enabled":            graph_enabled,
		"group_trigger_flag":       group_trigger_flag,
		"hosted_flag":              hosted_flag,
		"use_spark":                use_spark,
		"alert_metas":              fmt.Sprintf("%q", alert_metas),
		"dataset_ids":              fmt.Sprintf("%q", dataset_ids),
		"extend_dataset_ids":       fmt.Sprintf("%q", extend_dataset_ids),
		"check_condition":          check_condition,
		"extend_conf":              extend_conf,
	}

	endpoint := c.BuildRizhiyiURL(nil, "alerts")
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId(name)

	return nil
}

func resourceAlertRead(d *schema.ResourceData, m interface{}) error {
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

	appID, err := c.GetResourceIdByName(name, "alerts")
	if err != nil {
		return err
	}
	if appID == "" {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	category := d.Get("category").(int)
	check_interval := d.Get("check_interval").(int)
	continuous_trigger_value := d.Get("continuous_trigger_value").(int)
	market_day := d.Get("market_day").(bool)
	max_restrain_interval := d.Get("max_restrain_interval").(int)
	restrain_interval := d.Get("restrain_interval").(int)
	schedule_priority := d.Get("schedule_priority").(int)
	executor_id := d.Get("executor_id").(int)
	query := d.Get("query").(string)
	check_condition_group := d.Get("check_condition_group").(string)
	extend_query := d.Get("extend_query").(string)
	crontab := d.Get("crontab").(string)
	app_ids := d.Get("app_ids").(string)
	group_suppress_field := d.Get("group_suppress_field").(string)
	rt_names := d.Get("rt_names").(string)
	schedule_window := d.Get("schedule_window").(string)
	segmentation_field := d.Get("segmentation_field").(string)
	statistics_field := d.Get("statistics_field").(string)
	topic := d.Get("topic").(string)
	window := d.Get("window").(string)
	alert_when_recover := d.Get("alert_when_recover").(bool)
	enabled := d.Get("enabled").(bool)
	extend_use_spark := d.Get("extend_use_spark").(bool)
	graph_enabled := d.Get("graph_enabled").(bool)
	group_trigger_flag := d.Get("group_trigger_flag").(bool)
	hosted_flag := d.Get("hosted_flag").(bool)
	use_spark := d.Get("use_spark").(bool)
	alert_metas := d.Get("alert_metas").([]interface{})
	dataset_ids := d.Get("dataset_ids").([]interface{})
	extend_dataset_ids := d.Get("extend_dataset_ids").([]interface{})
	check_condition := d.Get("check_condition").(string)
	extend_conf := d.Get("extend_conf").(string)

	requestBody := map[string]interface{}{
		"name":                     name,
		"description":              description,
		"category":                 category,
		"check_interval":           check_interval,
		"continuous_trigger_value": continuous_trigger_value,
		"market_day":               market_day,
		"max_restrain_interval":    max_restrain_interval,
		"restrain_interval":        restrain_interval,
		"schedule_priority":        schedule_priority,
		"query":                    query,
		"check_condition_group":    check_condition_group,
		"extend_query":             extend_query,
		"crontab":                  crontab,
		"app_ids":                  app_ids,
		"group_suppress_field":     group_suppress_field,
		"rt_names":                 rt_names,
		"schedule_window":          schedule_window,
		"executor_id":              executor_id,
		"segmentation_field":       segmentation_field,
		"statistics_field":         statistics_field,
		"topic":                    topic,
		"window":                   window,
		"alert_when_recover":       alert_when_recover,
		"enabled":                  enabled,
		"extend_use_spark":         extend_use_spark,
		"graph_enabled":            graph_enabled,
		"group_trigger_flag":       group_trigger_flag,
		"hosted_flag":              hosted_flag,
		"use_spark":                use_spark,
		"alert_metas":              fmt.Sprintf("%q", alert_metas),
		"dataset_ids":              fmt.Sprintf("%q", dataset_ids),
		"extend_dataset_ids":       fmt.Sprintf("%q", extend_dataset_ids),
		"check_condition":          check_condition,
		"extend_conf":              extend_conf,
	}

	update_id, _ := c.GetResourceIdByName(d.Id(), "alerts")
	endpoint := c.BuildRizhiyiURL(nil, "alerts", update_id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d.SetId(name)
	return nil
}

func resourceAlertDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	del_id, err := c.GetResourceIdByName(name, "alerts")
	if err != nil {
		return err
	}

	endpoint := c.BuildRizhiyiURL(nil, "alerts", del_id)

	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d.SetId("")
	return nil
}
