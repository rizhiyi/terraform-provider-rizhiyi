provider "rizhiyi" {
  host  = "YOTTAWEB-IP"
  token = "" # Base64 encoded username:password
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
  name = "terraform_test_update"
}

//rizhiyi alert create
resource "rizhiyi_alert" "test_alert" {
  name            = "terraform_test_alert_update"
  category        = 0
  query           = "index=\"monitor\" *"
  check_condition = "{\"threshold\":\"info:500\",\"function\":\"count\",\"operator\":\">\",\"timerange\":\"-30min\"}"
  executor_id     = 1

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
  email      = "terrafor_testm@rizhiyi.com"
  passwd   = "Changeme"
  role_ids   = rizhiyi_role.new_role.app_id
  depends_on = [
    rizhiyi_role.new_role
  ]
}
