package metalcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetCreate,
		ReadContext:   resourceSubnetRead,
		DeleteContext: resourceSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"infrastructure_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == 0 {
						errs = append(errs, fmt.Errorf("%q is required. Provided value: %d", key, v))
					}
					return
				},
			},
			"network_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == 0 {
						errs = append(errs, fmt.Errorf("%q is required. Provided value: %d", key, v))
					}
					return
				},
			},
			"subnet_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"subnet_label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				Computed: true,
				ForceNew: true,
				//this is required because on the serverside the labels are converted to lowercase automatically
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}

					if new == "" {
						return true
					}
					return false
				},
			},
			"subnet_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
			"subnet_automatic_allocation": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
				Optional: true,
				ForceNew: true,
			},
			"subnet_prefix_size": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
			"subnet_is_ip_range": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
			"subnet_ip_range_ip_count": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
			"subnet_pool_id": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
			"subnet_range_start_human_readable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_range_end_human_readable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_netmask_human_readable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_gateway_human_readable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSubnetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	var infrastructure_id int
	var err error

	switch d.Get("infrastructure_id").(type) {
	case int:
		infrastructure_id = d.Get("infrastructure_id").(int)
	case string:
		infrastructure_id, err = strconv.Atoi(d.Get("infrastructure_id").(string))
		if err != nil {
			return diag.Errorf("Could not convert input '%s' to int", d.Get("infrastructure_id").(string))
		}
	}

	_, err = client.InfrastructureGet(infrastructure_id)

	if err != nil {
		return diag.Errorf("Infrastructure with id %+v not found.", infrastructure_id)
	}

	subnet := expandSubnet(d)

	Subnet, err := client.SubnetCreate(subnet)
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%d", Subnet.SubnetID)
	d.SetId(id)

	return resourceSubnetRead(ctx, d, meta)
}

func resourceSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	n, err := client.SubnetGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenSubnet(d, *n)

	return diags

}

func resourceSubnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.SubnetGet(id)

	//if the Subnet has already been deleted by infrastructure delete we ignore it, else we actually delete because
	//that might have been the intent. Note that only LAN Subnets can be deleted.
	//SAN and WAN are automatically created and deleted
	if err == nil {

		err = client.SubnetDelete(id)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId("")
	return diags
}

func flattenSubnet(d *schema.ResourceData, Subnet mc.Subnet) map[string]interface{} {

	d.Set("subnet_id", Subnet.SubnetID)
	d.Set("network_id", Subnet.NetworkID)
	d.Set("subnet_label", Subnet.SubnetLabel)
	d.Set("subnet_type", Subnet.SubnetType)
	d.Set("subnet_automatic_allocation", Subnet.SubnetAutomaticAllocation)
	d.Set("infrastructure_id", Subnet.InfrastructureID)
	d.Set("subnet_range_start_human_readable", Subnet.SubnetRangeStartHumanReadable)
	d.Set("subnet_range_end_human_readable", Subnet.SubnetRangeEndHumanReadable)
	d.Set("subnet_netmask_human_readable", Subnet.SubnetNetmaskHumanReadable)
	d.Set("subnet_gateway_human_readable", Subnet.SubnetGatewayHumanReadable)
	d.Set("subnet_pool_id", Subnet.SubnetPoolID)
	d.Set("cluster_id", Subnet.ClusterID)
	d.Set("subnet_is_ip_range", Subnet.SubnetIsIPRange)
	d.Set("subnet_ip_range_ip_count", Subnet.SubnetIPRangeCount)
	d.Set("subnet_prefix_size", Subnet.SubnetPrefixSize)

	return nil
}

func expandSubnet(d *schema.ResourceData) mc.Subnet {
	var n mc.Subnet

	if d.Get("subnet_id") != nil {
		n.SubnetID = d.Get("subnet_id").(int)
	}
	n.NetworkID = d.Get("network_id").(int)
	n.SubnetLabel = d.Get("subnet_label").(string)
	n.SubnetType = d.Get("subnet_type").(string)
	n.SubnetAutomaticAllocation = d.Get("subnet_automatic_allocation").(bool)
	n.SubnetPrefixSize = d.Get("subnet_prefix_size").(int)
	n.ClusterID = d.Get("cluster_id").(int)
	n.SubnetPoolID = d.Get("subnet_pool_id").(int)
	n.SubnetIsIPRange = d.Get("subnet_is_ip_range").(bool)
	n.SubnetIPRangeCount = d.Get("subnet_ip_range_ip_count").(int)

	return n
}
