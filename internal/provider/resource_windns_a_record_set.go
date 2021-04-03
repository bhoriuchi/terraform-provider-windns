package provider

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/bhoriuchi/terraform-provider-windns/windns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDnsARecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsARecordSetCreate,
		Read:   resourceDnsARecordSetRead,
		Update: resourceDnsARecordSetUpdate,
		Delete: resourceDnsARecordSetDelete,
		/*
			Importer: &schema.ResourceImporter{
				State: resourceDnsImport,
			},
		*/

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateZone,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"addresses": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      hashIPString,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  3600,
			},
		},
	}
}

func resourceDnsARecordSetCreate(d *schema.ResourceData, meta interface{}) error {
	if meta == nil {
		return fmt.Errorf("client not created")
	}

	client := meta.(*windns.Client)
	d.SetId(resourceFQDN(d))
	name := d.Get("name").(string)
	zone := d.Get("zone").(string)
	ttl := d.Get("ttl").(int)
	addresses := d.Get("addresses").(*schema.Set).List()

	for _, address := range addresses {
		rsp, err := client.AddARecord(&windns.AddARecordOptions{
			Name:     name,
			Address:  address.(string),
			ZoneName: zone,
			TTL:      ttl,
		})
		if err != nil {
			d.SetId("")
			return err
		} else if rsp.Code != http.StatusOK {
			d.SetId("")
			return fmt.Errorf(rsp.Detail)
		}
	}

	return resourceDnsARecordSetRead(d, meta)
}

func resourceDnsARecordSetRead(d *schema.ResourceData, meta interface{}) error {
	if meta == nil {
		return fmt.Errorf("client not created")
	}

	client := meta.(*windns.Client)
	rsp, err := client.ReadARecord(&windns.ReadARecordOptions{
		Name:     d.Get("name").(string),
		ZoneName: d.Get("zone").(string),
	})
	if err != nil {
		d.SetId("")
		return err
	} else if rsp.Code != http.StatusOK {
		d.SetId("")
		return fmt.Errorf(rsp.Detail)
	}

	if len(rsp.Records) > 0 {
		var ttl sort.IntSlice
		addresses := schema.NewSet(hashIPString, nil)
		for _, record := range rsp.Records {
			addresses.Add(record.Data)
			ttl = append(ttl, record.TTL)
		}
		sort.Sort(ttl)

		d.Set("addresses", addresses)
		d.Set("ttl", ttl[0])
	} else {
		d.SetId("")
	}

	return nil
}

func resourceDnsARecordSetUpdate(d *schema.ResourceData, meta interface{}) error {
	if meta == nil {
		return fmt.Errorf("client not created")
	}

	client := meta.(*windns.Client)
	name := d.Get("name").(string)
	zone := d.Get("zone").(string)
	ttl := d.Get("ttl").(int)

	if d.HasChange("addresses") {
		o, n := d.GetChange("addresses")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		// Loop through all the old addresses and remove them
		for _, addr := range remove {
			rsp, err := client.DeleteARecord(&windns.DeleteARecordOptions{
				Name:     name,
				ZoneName: zone,
				Address:  addr.(string),
			})
			if err != nil {
				d.SetId("")
				return fmt.Errorf("Error updating DNS record: %s", err)
			} else if rsp.Code != http.StatusOK {
				d.SetId("")
				return fmt.Errorf(rsp.Detail)
			}
		}
		// Loop through all the new addresses and insert them
		for _, addr := range add {
			rsp, err := client.AddARecord(&windns.AddARecordOptions{
				Name:     name,
				ZoneName: zone,
				Address:  addr.(string),
				TTL:      ttl,
			})
			if err != nil {
				d.SetId("")
				return fmt.Errorf("Error updating DNS record: %s", err)
			} else if rsp.Code != http.StatusOK {
				d.SetId("")
				return fmt.Errorf(rsp.Detail)
			}
		}
	}

	return resourceDnsARecordSetRead(d, meta)
}

func resourceDnsARecordSetDelete(d *schema.ResourceData, meta interface{}) error {
	if meta == nil {
		return fmt.Errorf("client not created")
	}

	client := meta.(*windns.Client)
	name := d.Get("name").(string)
	zone := d.Get("zone").(string)
	addresses := d.Get("addresses").(*schema.Set).List()

	for _, address := range addresses {
		rsp, err := client.DeleteARecord(&windns.DeleteARecordOptions{
			Name:     name,
			Address:  address.(string),
			ZoneName: zone,
		})
		if err != nil {
			d.SetId("")
			return err
		} else if rsp.Code != http.StatusOK {
			d.SetId("")
			return fmt.Errorf(rsp.Detail)
		}
	}

	return nil
}
