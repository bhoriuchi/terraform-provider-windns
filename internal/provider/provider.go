package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// New returns a *schema.Provider for Windows DNS updates.
func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"dns_server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Windows DNS server to issue WinRM commands against.",
			},
			"proxy_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A WinRM proxy host to indirectly run WinRM commands against the dns_server. This is useful if you cannot elevate the permissions on the dns_server for the user.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username to authenticate with.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain for username.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password for username.",
			},
			"secure_transport": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Issue commands on a secure transport.",
			},
			"ignore_ssl_checks": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, connect and ignore any untrusted/invalid SSL.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Port to connect over. Defaults to 5985 when secure_transport is false and 5986 when true.",
			},
			"kdc_server": {
				Type:        schema.TypeList,
				MinItems:    1,
				Optional:    true,
				Elem:        schema.TypeString,
				Description: "List of KDCs. Will assume the dns_server is a KDC if omitted.",
			},
			"krb5_conf": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom krb5.conf configuration specified as a string.",
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time in seconds before the operation should timeout.",
			},
		},
	}
}
