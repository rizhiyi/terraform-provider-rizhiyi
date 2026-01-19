package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"terraform-provider-rizhiyi/yottaweb"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RIZHIYI_HOST", nil),
				Description: "Rizhiyi host (e.g. 192.168.1.224)",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RIZHIYI_TOKEN", nil),
				Description: "Rizhiyi authorization token (Base64 encoded username:password)",
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rizhiyi_role":        resourceRoles(),
			"rizhiyi_index":       resourceIndex(),
			"rizhiyi_dashboard":   resourceDashboards(),
			"rizhiyi_alert":       resourceAlert(),
			"rizhiyi_parser_rule": resourceParserRule(),
			"rizhiyi_account":     resourceAccount(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	host := d.Get("host").(string)
	token := d.Get("token").(string)
	return yottaweb.NewClient(host, token), nil
}
