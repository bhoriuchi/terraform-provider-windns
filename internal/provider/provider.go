package provider

import (
	"fmt"

	"github.com/bhoriuchi/terraform-provider-windns/windns"
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
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password for username.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain for username.",
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
			"kdc_servers": {
				Type:     schema.TypeList,
				MinItems: 1,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of KDCs. Will assume the dns_server is a KDC if omitted.",
			},
			"krb5_conf": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom krb5.conf configuration specified as a string.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Time in seconds before the operation should timeout. Defaults to 60 seconds",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dns_a_record_set": resourceDnsARecordSet(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	opts := &windns.Options{
		DnsServer:       d.Get("dns_server").(string),
		Username:        d.Get("username").(string),
		Password:        d.Get("password").(string),
		Domain:          d.Get("domain").(string),
		SecureTransport: d.Get("secure_transport").(bool),
		SkipSSLVerify:   d.Get("ignore_ssl_checks").(bool),
	}

	if port, ok := d.GetOk("port"); ok {
		opts.Port = port.(int)
	}

	if timeout, ok := d.GetOk("timeout"); ok {
		opts.TimeoutSeconds = timeout.(int)
	}

	if krb5conf, ok := d.GetOk("krb5_conf"); ok {
		opts.KRB5Conf = krb5conf.(string)
	}

	if proxyHost, ok := d.GetOk("proxy_host"); ok {
		opts.ProxyHost = proxyHost.(string)
	}

	if val, ok := d.GetOk("kdc_servers"); ok {
		opts.KDCServers = []string{}
		for _, kdc := range val.([]interface{}) {
			opts.KDCServers = append(opts.KDCServers, kdc.(string))
		}
	}

	return windns.NewClient(opts)
}

// https://github.com/hashicorp/terraform-provider-dns/blob/main/internal/provider/provider.go
func resourceDnsImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	/*
			record := d.Id()
			if !IsFqdn(record) {
				return nil, fmt.Errorf("Not a fully-qualified DNS name: %s", record)
			}

			labels := dns.SplitDomainName(record)

			msg := new(dns.Msg)

			var zone *string

		Loop:
			for l := range labels {

				msg.SetQuestion(dns.Fqdn(strings.Join(labels[l:], ".")), dns.TypeSOA)

				r, err := exchange(msg, true, meta)
				if err != nil {
					return nil, fmt.Errorf("Error querying DNS record: %s", err)
				}

				switch r.Rcode {
				case dns.RcodeSuccess:

					if len(r.Answer) == 0 {
						continue
					}

					for _, ans := range r.Answer {
						switch t := ans.(type) {
						case *dns.SOA:
							zone = &t.Hdr.Name
						case *dns.CNAME:
							continue Loop
						}
					}

					break Loop
				case dns.RcodeNameError:
					continue
				default:
					return nil, fmt.Errorf("Error querying DNS record: %v (%s)", r.Rcode, dns.RcodeToString[r.Rcode])
				}
			}

			if zone == nil {
				return nil, fmt.Errorf("No SOA record in authority section in response for %s", record)
			}

			common := dns.CompareDomainName(record, *zone)
			if common == 0 {
				return nil, fmt.Errorf("DNS record %s shares no common labels with zone %s", record, *zone)
			}

			d.Set("zone", *zone)
			if name := strings.Join(labels[:len(labels)-common], "."); name != "" {
				d.Set("name", name)
			}

			return []*schema.ResourceData{d}, nil
	*/
	return []*schema.ResourceData{d}, nil
}

func resourceFQDN(d *schema.ResourceData) string {

	fqdn := d.Get("zone").(string)

	if name, ok := d.GetOk("name"); ok {
		fqdn = fmt.Sprintf("%s.%s", name.(string), fqdn)
	}

	return fqdn
}
