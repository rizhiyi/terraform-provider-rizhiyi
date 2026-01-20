variable "rizhiyi_host" {
  type        = string
  description = "Rizhiyi host URL"
  default     = null
}

variable "rizhiyi_token" {
  type        = string
  description = "Rizhiyi API token"
  default     = null
}

provider "rizhiyi" {
  host  = var.rizhiyi_host
  token = var.rizhiyi_token
}

//rizhiyi roles create
resource "rizhiyi_role" "new_role" {
  name = "new_roles_01"
  memo = "This is test roles1"
}

// rizhiyi index create
resource "rizhiyi_index" "test_index" {
  pattern         = "kCompression"
  name            = "terraform_index"
  description     = "This is terraform test index"
  rotation_period = "10d"
  expired_time    = "25d"

}

//rizhiyi dashboard create
resource "rizhiyi_dashboard" "test_dashboard" {
  name            = "terraform_test_update"
  app_id          = 1
  export          = "local"
  data_user       = "viewer"
  default_display = 0
  active_tab      = 0
  manage_tabs     = true
  
  tabs {
    name    = "循环图"
    content = jsonencode({
      refresh = {
        time = 3
        unit = "m"
        on = false
        showRefreshProcess = true
      }
      showFilters = false
      showTitle = true
      editable = true
      scheme = "schemecat1"
      theme = "day"
      activeDrilldown = false
      autoUpdate = true
      filters = []
      widgets = [
        {
          y = 0
          x = 0
          w = 12
          h = 5
          type = "trend"
          trendId = 107
          importType = "clone"
          id = "05d97f209d08421c94aceb96ca0c73fe"
          searchData = {
            query = "starttime=\"now/d\" endtime=\"now/d+24h\" tag:sample04061424_chart | stats count() by hostname,apache.clientip |sort by +apache.clientip |limit 5"
            time_range = "-10m,now"
            now = ""
            market_day = 0
            highlight = false
            onlySortByTimestamp = false
            chartType = "sequence"
            use_spark = false
            scheme = "schemecat1"
            dataset_ids = "[]"
            trendName = "仪表盘循环图"
            config = [
              { xField = "" },
              { fromField = "" },
              { toField = "" },
              { byFields = [] },
              { labelField = "" }
            ]
            xField = "hostname"
            fromField = "apache.clientip"
            toField = "hostname"
            byFields = ["apache.clientip"]
            labelField = "apache.clientip"
            chartStartingColor = "#3661EB"
          }
        }
      ]
    })
  }
}

resource "rizhiyi_dashboard" "web_app" {
  name            = "Web应用监控"
  app_id          = 1
  export          = "local"
  data_user       = "viewer"
  default_display = 0
  active_tab      = 0
  manage_tabs     = true

  tabs {
    name    = "总览"
    content = jsonencode({
      refresh = {
        time = 3
        unit = "m"
        on = false
        showRefreshProcess = true
      }
      showFilters = true
      showTitle = true
      editable = true
      scheme = "schemecat1"
      theme = "day"
      activeDrilldown = false
      autoUpdate = true
      filters = [
        {
          visible                       = true
          filterType                    = "token"
          title                         = "filterTitle"
          token                         = "filterId"
          description                   = "an input token here"
          width                         = "200px"
          searchOnChange                = false
          type                          = "text"
          textValue                     = "default value"
          prefix                        = ""
          suffix                        = ""
        },
        {
          visible                       = true
          filterType                    = "token"
          title                         = "selectorTitle"
          token                         = "selectorId"
          description                   = "a selection token here"
          width                         = "200px"
          searchOnChange                = false
          field                         = ""
          type                          = "dynamicDropdown"
          selectionMode                 = "multiple" // single
          dynamicFieldName              = ""
          dynamicFieldValue             = "ip"
          dynamicQuery                  = "* | top 10 ip"
          textValue                     = "vendors"
          dropdownSelectedValue         = "*"
          dynamicSearchTimeRange        = "-10m,now"
          dynamicSearchUseGlobalTimeRange = true
          setAsGlobal                   = false
          dropdownValues                = [{
            "label": "all",
            "value": "*"
          }]
          prefix                        = "src_ip:("
          suffix                        = ")"
          valueSuffix                   = ""
          valuePrefix                   = ""
          delimiter                     = " OR "
        },
        {
          visible                       = true
          filterType                    = "token"
          title                         = "globalTimerange"
          token                         = "globalTimeRange"
          description                   = "a time picker token here"
          width                         = "200px"
          searchOnChange                = true
          timeValue                     = "-10m,now"
          type                          = "timerange"
          setAsGlobal                   = false
        }
      ]
      widgets = [
        {
          y = 0
          x = 0
          w = 6
          h = 5
          type = "trend"
          importType = "clone"
          id = "overview-reqs"
          searchData = {
            query = "starttime=\"now/d\" endtime=\"now/d+24h\" index=\"yotta\" * | timechart span=1m count()"
            time_range = "-1h,now"
            chartType = "line"
            xField = ""
          }
        }
      ]
    })
  }
}


//rizhiyi alert create
resource "rizhiyi_alert" "test1_alert" {
  name            = "terraform_test1_alert_update"
  category        = 0
  query           = "index=monitor *"
  check_interval  = 300
  crontab         = "0"
  extend_conf     = jsonencode({})
  timezone        = "Asia/Shanghai"
  check_condition = "{\"threshold\":\"info:500\",\"function\":\"count\",\"operator\":\">\",\"timerange\":\"-30m\"}"
  executor_id     = 1
  app_id          = 1
}

//rizhiyi parser rules create
resource "rizhiyi_parser_rule" "test_parser_rule" {
  name        = "test_rule_update"
  logtype     = "json"
  conf        = "[{\"json\":{\"rule\":[{\"add_fields\":[],\"source\":\"raw_message\",\"another_name\":\"\",\"paths\":[],\"extract_limit\":\"\"}]}}]"
  category_id = 1000
}

//rizhiyi account create
resource "rizhiyi_account" "test_account" {
  name       = "terraform_test_update"
  email      = "vYzrZ7Cygw+lgXavlVMuAgzzyHPfSRofgd0I4Bd4jwA="
  phone      = "L6lTtIrxjcpl36dn37wEKg=="
  passwd     = "Changeme"
  role_ids   = rizhiyi_role.new_role.app_id
  depends_on = [
    rizhiyi_role.new_role
  ]
}
