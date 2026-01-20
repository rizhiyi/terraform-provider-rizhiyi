package provider

import (
	"encoding/json"
	"io"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"net/url"
	"fmt"
	"time"
	"terraform-provider-rizhiyi/yottaweb"
)

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertCreate,
		Read:   resourceAlertRead,
		Update: resourceAlertUpdate,
		Delete: resourceAlertDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Resource name for the new Alert resource.",
			},
			"category": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The monitoring types for the new Alert resource are as follows: 0.Event Count Monitoring, 1.Field Statistics Monitoring, 2.Continuous Statistics Monitoring, 3.Baseline Comparison Monitoring, 4.SPL Statistics Monitoring. (default value 0)",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Search content for the alert resource.",
			},
			"check_condition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Monitoring trigger conditions for the Alert resource.",
			},
			"executor_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "User ID for executing the new Alert resource.（default value 0）",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the new Alert resource.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The field to enable monitoring for the Alert resource. (default value false)",
			},
			"crontab": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The crontab execution schedule for the alert resource, please provide the corresponding cron statement, for example, 0 * * * * ？, where 0 indicates not using the crontab execution schedule.",
			},
			"check_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The scheduled execution plan for the new Alert resource, fill in the interval in seconds for the scheduled execution plan, where 0 indicates no scheduled execution plan. (default value 0)",
			},
			"restrain_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Monitoring suppression time (in seconds) for the alert resource. (default value 0)",
			},
			"max_restrain_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The doubling time (in seconds) for suppressing the cancellation of alert monitoring. (default value 0)",
			},
			"continuous_trigger_value": {
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"use_spark": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the Alert resource uses advanced mode. (default value false)",
			},
			"extend_use_spark": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the extended search of the Alert resource uses advanced mode. (default value false)",
			},
			"graph_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the extended search of the Alert resource has the effect illustration enabled. (default value false)",
			},
			"extend_query": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search content for the extended search of the alert resource.",
			},
			"extend_conf": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fixed key-value for the extended search of the Alert resource.",
			},
			"dataset_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "JSON string for the dataset node ID of the Alert resource.",
			},
			"extend_dataset_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "JSON string for the dataset node ID of the extended search in the Alert resource.",
			},
			"segmentation_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Device split field for the Alert resource, an empty string indicates no device split.",
			},
			"statistics_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"market_day": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the new Alert resource is executed only on the transaction day. (default value false)",
			},
			"schedule_priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"schedule_window": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"window": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_condition_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_trigger_flag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"hosted_flag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"group_suppress_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_metas": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Plugin data for the Alert resource (in JSON array string format), where each item in the array should provide the plugin's name, trigger level, configuration information, and change data.",
			},
			"alert_when_recover": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the Alert resource uses monitoring reply prompts.(default value false)",
			},
			"app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rt_names": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource group name to which the Alert resource belongs, for example: default_Alert, test.",
			},
		},
	}
}

func resourceAlertCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	category := d.Get("category").(int)
	query := d.Get("query").(string)
	check_condition := d.Get("check_condition").(string)
	executor_id := d.Get("executor_id").(int)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	crontab := d.Get("crontab").(string)
	check_interval := d.Get("check_interval").(int)
	restrain_interval := d.Get("restrain_interval").(int)
	max_restrain_interval := d.Get("max_restrain_interval").(int)
	continuous_trigger_value := d.Get("continuous_trigger_value").(int)
	use_spark := d.Get("use_spark").(bool)
	extend_use_spark := d.Get("extend_use_spark").(bool)
	graph_enabled := d.Get("graph_enabled").(bool)
	extend_query := d.Get("extend_query").(string)
	extend_conf := d.Get("extend_conf").(string)
	dataset_ids := d.Get("dataset_ids").([]interface{})
	extend_dataset_ids := d.Get("extend_dataset_ids").([]interface{})
	datasetIDsStr := ""
	if len(dataset_ids) == 0 {
		datasetIDsStr = "[]"
	} else if b, err := json.Marshal(dataset_ids); err == nil {
		datasetIDsStr = string(b)
	}
	extendDatasetIDsStr := ""
	if len(extend_dataset_ids) == 0 {
		extendDatasetIDsStr = "[]"
	} else if b, err := json.Marshal(extend_dataset_ids); err == nil {
		extendDatasetIDsStr = string(b)
	}
	segmentation_field := d.Get("segmentation_field").(string)
	statistics_field := d.Get("statistics_field").(string)
	market_day := d.Get("market_day").(bool)
	schedule_priority := d.Get("schedule_priority").(int)
	schedule_window := d.Get("schedule_window").(string)
	window := d.Get("window").(string)
	topic := d.Get("topic").(string)
	check_condition_group := d.Get("check_condition_group").(string)
	group_trigger_flag := d.Get("group_trigger_flag").(bool)
	hosted_flag := d.Get("hosted_flag").(bool)
	alert_metas := d.Get("alert_metas").([]interface{})
	alert_when_recover := d.Get("alert_when_recover").(bool)
	app_id := d.Get("app_id").(int)
	group_suppress_field := d.Get("group_suppress_field").(string)
	timezone := d.Get("timezone").(string)
	rt_names := d.Get("rt_names").(string)

	alertMetasStr := ""
	if len(alert_metas) == 0 {
		alertMetasStr = "[]"
	} else if b, err := json.Marshal(alert_metas); err == nil {
		alertMetasStr = string(b)
	}
	if extend_conf == "" {
		extend_conf = "{}"
	}
	if crontab == "" {
		crontab = "0"
	}
	marketDayInt := 0
	if market_day {
		marketDayInt = 1
	}
	if timezone == "" {
		timezone = "Asia/Shanghai"
	}

	requestBody := map[string]interface{}{
		"name":                    name,
		"category":                category,
		"query":                   query,
		"check_condition":         check_condition,
		"executor_id":             executor_id,
		"description":             description,
		"enabled":                 enabled,
		"crontab":                 crontab,
		"check_interval":          check_interval,
		"restrain_interval":       restrain_interval,
		"max_restrain_interval":   max_restrain_interval,
		"continuous_trigger_value": continuous_trigger_value,
		"use_spark":               use_spark,
		"extend_use_spark":        extend_use_spark,
		"graph_enabled":           graph_enabled,
		"extend_query":            extend_query,
		"extend_conf":             extend_conf,
		"dataset_ids":             datasetIDsStr,
		"extend_dataset_ids":      extendDatasetIDsStr,
		"segmentation_field":      segmentation_field,
		"market_day":              marketDayInt,
		"schedule_priority":       schedule_priority,
		"schedule_window":         schedule_window,
		"window":                  window,
		"topic":                   topic,
		"check_condition_group":   check_condition_group,
		"group_trigger_flag":      group_trigger_flag,
		"hosted_flag":             hosted_flag,
		"alert_metas":             alertMetasStr,
		"alert_when_recover":      alert_when_recover,
		"group_suppress_field":    group_suppress_field,
		"timezone":                timezone,
		"rt_names":                rt_names,
		"composite_info":          nil,
		"alert_condition":         nil,
		"recover_condition":       nil,
	}
	if statistics_field != "" {
		requestBody["statistics_field"] = statistics_field
	}
	if app_id > 0 {
		requestBody["app_id"] = app_id
	}

	endpoint := c.BuildRizhiyiURL(nil, "v3", "alerts")
	resp, err := c.Post(endpoint, requestBody)
	if err != nil {
		return err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var result map[string]interface{}
	json.Unmarshal(bodyBytes, &result)
	finalID := ""
	if obj, ok := result["object"].(map[string]interface{}); ok {
		if idVal, ok := obj["id"]; ok {
			switch v := idVal.(type) {
			case float64:
				finalID = strconv.Itoa(int(v))
			case int:
				finalID = strconv.Itoa(v)
			case string:
				finalID = v
			default:
				finalID = ""
			}
		} else {
			finalID = ""
		}
	} else {
		finalID = ""
	}
	// 如果响应未提供 id，则通过名称查询 id，确保 state 使用后端真实 ID
	if finalID == "" {
		if idByName, err := getAlertIdByName(c, name); err == nil && idByName != "" {
			finalID = idByName
		}
	}
	if finalID == "" {
		for i := 0; i < 30 && finalID == ""; i++ {
			if idByName, err := getAlertIdByName(c, name); err == nil && idByName != "" {
				finalID = idByName
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
	if finalID == "" {
		return fmt.Errorf("alert created but id not resolvable within timeout: %s", name)
	}
	d.SetId(finalID)

	return resourceAlertRead(d, m)
}

func resourceAlertRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	id := d.Id()
	if id == "" {
		name := ""
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		}
		if name == "" {
			d.SetId("")
			return nil
		}

		appID, err := c.GetResourceIdByName(name, "v3", "alerts")
		if err != nil {
			return err
		}
		if appID == "" {
			d.SetId("")
			return nil
		}
		id = appID
	}

	data, err := c.GetResourceById(id, "v3", "alerts")
	if err != nil {
		nameVal := ""
		if v, ok := d.GetOk("name"); ok {
			nameVal = v.(string)
		}
		for i := 0; i < 20 && data == nil; i++ {
			if nameVal != "" {
				newID, e := getAlertIdByName(c, nameVal)
				if e == nil && newID != "" {
					d.SetId(newID)
					id = newID
					data, err = c.GetResourceById(id, "v3", "alerts")
					if err == nil && data != nil {
						break
					}
				} else {
				}
			} else {
			}
			time.Sleep(1 * time.Second)
		}
		if data == nil {
			d.SetId("")
			return nil
		}
	}

	d.Set("name", data["name"])
	if v, ok := data["category"]; ok {
		d.Set("category", v)
	}
	d.Set("query", data["query"])
	d.Set("check_condition", data["check_condition"])
	d.Set("executor_id", data["executor_id"])
	d.Set("description", data["description"])
	d.Set("enabled", data["enabled"])
	d.Set("crontab", data["crontab"])
	d.Set("check_interval", data["check_interval"])
	d.Set("restrain_interval", data["restrain_interval"])
	d.Set("max_restrain_interval", data["max_restrain_interval"])
	d.Set("continuous_trigger_value", data["continuous_trigger_value"])
	d.Set("use_spark", data["use_spark"])
	d.Set("extend_use_spark", data["extend_use_spark"])
	d.Set("graph_enabled", data["graph_enabled"])
	d.Set("extend_query", data["extend_query"])
	d.Set("extend_conf", data["extend_conf"])
	if s, ok := data["dataset_ids"].(string); ok {
		var arr []interface{}
		if s != "" {
			json.Unmarshal([]byte(s), &arr)
		}
		d.Set("dataset_ids", arr)
	} else {
		d.Set("dataset_ids", data["dataset_ids"])
	}
	if s, ok := data["extend_dataset_ids"].(string); ok {
		var arr []interface{}
		if s != "" {
			json.Unmarshal([]byte(s), &arr)
		}
		d.Set("extend_dataset_ids", arr)
	} else {
		d.Set("extend_dataset_ids", data["extend_dataset_ids"])
	}
	d.Set("segmentation_field", data["segmentation_field"])
	d.Set("statistics_field", data["statistics_field"])
	d.Set("market_day", data["market_day"])
	d.Set("schedule_priority", data["schedule_priority"])
	d.Set("schedule_window", data["schedule_window"])
	d.Set("window", data["window"])
	d.Set("topic", data["topic"])
	d.Set("check_condition_group", data["check_condition_group"])
	d.Set("group_trigger_flag", data["group_trigger_flag"])
	d.Set("hosted_flag", data["hosted_flag"])
	if s, ok := data["alert_metas"].(string); ok {
		var arr []interface{}
		if s != "" {
			json.Unmarshal([]byte(s), &arr)
		}
		d.Set("alert_metas", arr)
	} else {
		d.Set("alert_metas", data["alert_metas"])
	}
	d.Set("alert_when_recover", data["alert_when_recover"])
	if v, ok := data["app_id"]; ok {
		switch vv := v.(type) {
		case float64:
			d.Set("app_id", int(vv))
		default:
			d.Set("app_id", v)
		}
	}
	d.Set("group_suppress_field", data["group_suppress_field"])
	d.Set("timezone", data["timezone"])
	d.Set("rt_names", data["rt_names"])

	return nil
}

// 精确通过 name 查询 alert 的 ID（使用 name 过滤器提高命中率）
func getAlertIdByName(c *yottaweb.Client, name string) (string, error) {
	params := url.Values{}
	params.Set("name", name)
	params.Set("size", "50")
	endpoint := c.BuildRizhiyiURL(params, "v3", "alerts")
	resp, err := c.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return "", err
	}
	var list []interface{}
	if v, ok := data["list"].([]interface{}); ok {
		list = v
	} else if v, ok := data["objects"].([]interface{}); ok {
		list = v
	}
	for _, item := range list {
		if obj, ok := item.(map[string]interface{}); ok {
			if obj["name"] == name {
				switch idVal := obj["id"].(type) {
				case float64:
					return strconv.Itoa(int(idVal)), nil
				case int:
					return strconv.Itoa(idVal), nil
				case string:
					return idVal, nil
				default:
					return "", nil
				}
			}
		}
	}
	// 回退：不使用过滤参数，拉全量列表再匹配
	endpoint = c.BuildRizhiyiURL(nil, "v3", "alerts")
	resp, err = c.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, _ = io.ReadAll(resp.Body)
	data = map[string]interface{}{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return "", err
	}
	list = nil
	if v, ok := data["list"].([]interface{}); ok {
		list = v
	} else if v, ok := data["objects"].([]interface{}); ok {
		list = v
	}
	for _, item := range list {
		if obj, ok := item.(map[string]interface{}); ok {
			if obj["name"] == name {
				switch idVal := obj["id"].(type) {
				case float64:
					return strconv.Itoa(int(idVal)), nil
				case int:
					return strconv.Itoa(idVal), nil
				case string:
					return idVal, nil
				default:
					return "", nil
				}
			}
		}
	}
	return "", nil
}

func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	name := d.Get("name").(string)
	id := d.Id()
	category := d.Get("category").(int)
	query := d.Get("query").(string)
	check_condition := d.Get("check_condition").(string)
	executor_id := d.Get("executor_id").(int)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	crontab := d.Get("crontab").(string)
	check_interval := d.Get("check_interval").(int)
	restrain_interval := d.Get("restrain_interval").(int)
	max_restrain_interval := d.Get("max_restrain_interval").(int)
	continuous_trigger_value := d.Get("continuous_trigger_value").(int)
	use_spark := d.Get("use_spark").(bool)
	extend_use_spark := d.Get("extend_use_spark").(bool)
	graph_enabled := d.Get("graph_enabled").(bool)
	extend_query := d.Get("extend_query").(string)
	extend_conf := d.Get("extend_conf").(string)
	dataset_ids := d.Get("dataset_ids").([]interface{})
	extend_dataset_ids := d.Get("extend_dataset_ids").([]interface{})
	segmentation_field := d.Get("segmentation_field").(string)
	statistics_field := d.Get("statistics_field").(string)
	market_day := d.Get("market_day").(bool)
	schedule_priority := d.Get("schedule_priority").(int)
	schedule_window := d.Get("schedule_window").(string)
	window := d.Get("window").(string)
	topic := d.Get("topic").(string)
	check_condition_group := d.Get("check_condition_group").(string)
	group_trigger_flag := d.Get("group_trigger_flag").(bool)
	hosted_flag := d.Get("hosted_flag").(bool)
	alert_metas := d.Get("alert_metas").([]interface{})
	alert_when_recover := d.Get("alert_when_recover").(bool)
	app_id := d.Get("app_id").(int)
	group_suppress_field := d.Get("group_suppress_field").(string)
	timezone := d.Get("timezone").(string)
	rt_names := d.Get("rt_names").(string)

	datasetIDsStr := ""
	if len(dataset_ids) == 0 {
		datasetIDsStr = ""
	} else if b, err := json.Marshal(dataset_ids); err == nil {
		datasetIDsStr = string(b)
	}
	extendDatasetIDsStr := ""
	if len(extend_dataset_ids) == 0 {
		extendDatasetIDsStr = ""
	} else if b, err := json.Marshal(extend_dataset_ids); err == nil {
		extendDatasetIDsStr = string(b)
	}

	alertMetasStr := ""
	if len(alert_metas) == 0 {
		alertMetasStr = ""
	} else if b, err := json.Marshal(alert_metas); err == nil {
		alertMetasStr = string(b)
	}

	requestBody := map[string]interface{}{
		"name":                    name,
		"category":                category,
		"query":                   query,
		"check_condition":         check_condition,
		"executor_id":             executor_id,
		"description":             description,
		"enabled":                 enabled,
		"crontab":                 crontab,
		"check_interval":          check_interval,
		"restrain_interval":       restrain_interval,
		"max_restrain_interval":   max_restrain_interval,
		"continuous_trigger_value": continuous_trigger_value,
		"use_spark":               use_spark,
		"extend_use_spark":        extend_use_spark,
		"graph_enabled":           graph_enabled,
		"extend_query":            extend_query,
		"extend_conf":             extend_conf,
		"dataset_ids":             datasetIDsStr,
		"extend_dataset_ids":      extendDatasetIDsStr,
		"segmentation_field":      segmentation_field,
		"statistics_field":        statistics_field,
		"market_day":              market_day,
		"schedule_priority":       schedule_priority,
		"schedule_window":         schedule_window,
		"window":                  window,
		"topic":                   topic,
		"check_condition_group":   check_condition_group,
		"group_trigger_flag":      group_trigger_flag,
		"hosted_flag":             hosted_flag,
		"alert_metas":             alertMetasStr,
		"alert_when_recover":      alert_when_recover,
		"group_suppress_field":    group_suppress_field,
		"timezone":                timezone,
		"rt_names":                rt_names,
	}
	if app_id > 0 {
		requestBody["app_id"] = app_id
	}

	if id == "" {
		update_id, _ := c.GetResourceIdByName(name, "v3", "alerts")
		id = update_id
	}
	endpoint := c.BuildRizhiyiURL(nil, "v3", "alerts", id)
	resp, err := c.Put(endpoint, requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d.SetId(id)
	return resourceAlertRead(d, m)
}

func resourceAlertDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*yottaweb.Client)
	id := d.Id()
	if id == "" {
		name := d.Get("name").(string)
		delID, err := c.GetResourceIdByName(name, "v3", "alerts")
		if err != nil {
			return err
		}
		id = delID
	}
	endpoint := c.BuildRizhiyiURL(nil, "v3", "alerts", id)
	resp, err := c.Delete(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	d.SetId("")
	return nil
}
