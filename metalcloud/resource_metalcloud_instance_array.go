package metalcloud

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func resourceInstanceArray() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceArrayCreate,
		ReadContext:   resourceInstanceArrayRead,
		UpdateContext: resourceInstanceArrayUpdate,
		DeleteContext: resourceInstanceArrayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"infrastructure_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance_array_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_array_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//this is required because on the serverside the labels are converted to lowercase automatically
				ValidateDiagFunc: validateLabel,
			},
			"instance_array_instance_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance_array_boot_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pxe_iscsi",
			},
			"instance_array_ram_gbytes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_processor_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_processor_core_mhz": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_processor_core_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_disk_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_disk_size_mbytes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_additional_wan_ipv4_json": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_array_custom_variables": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
			},
			"instance_custom_variables": {
				Type:     schema.TypeList,
				Elem:     instanceCustomVariableResource(),
				Optional: true,
			},
			"volume_template_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"instance_array_firewall_managed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"firewall_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceFirewallRule(),
				Set:      firewallRuleResourceHash,
				//TODO: set defaults so that we don't get the big list of serverside generated rules
			},

			"interface": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceInstanceArrayInterface(),
				Set:      interfaceResourceHash,
				//TODO: set defaults so that we don't get the big list of serverside generated rules
			},

			"instances": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func instanceCustomVariableResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_index": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"custom_variables": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"firewall_rule_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"firewall_rule_port_range_start": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				//Default:  1,
			},
			"firewall_rule_port_range_end": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				//Default:  65535,
			},
			"firewall_rule_source_ip_address_range_start": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"firewall_rule_source_ip_address_range_end": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"firewall_rule_destination_ip_address_range_start": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"firewall_rule_destination_ip_address_range_end": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"firewall_rule_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tcp",
			},
			"firewall_rule_ip_address_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ipv4",
			},
			"firewall_rule_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceInstanceArrayInterface() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"interface_index": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network_label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceInstanceArrayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	ia := mc.InstanceArray{
		InstanceArrayLabel: d.Get("instance_array_label").(string),
	}

	infrastructure_id := d.Get("infrastructure_id").(int)

	iaC, err := client.InstanceArrayCreate(infrastructure_id, ia)
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%d", iaC.InstanceArrayID)

	d.SetId(id)

	return resourceInstanceArrayRead(ctx, d, meta)
}

func resourceInstanceArrayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ia, err := client.InstanceArrayGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("instance_array_label", ia.InstanceArrayLabel)
	d.Set("instance_array_id", ia.InstanceArrayID)

	return diags

}

func resourceInstanceArrayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange("instance_array_label") {
		client := meta.(*mc.Client)

		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		ia, err := client.InstanceArrayGet(id)
		if err != nil {
			return diag.FromErr(err)
		}

		bSwapExistingInstancesHardware := false
		bkeepDetachingDrives := false
		ia.InstanceArrayOperation.InstanceArrayLabel = d.Get("instance_array_label").(string)

		ia, err = client.InstanceArrayEdit(id, *ia.InstanceArrayOperation, &bSwapExistingInstancesHardware, &bkeepDetachingDrives, nil, nil)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	return resourceInstanceArrayRead(ctx, d, meta)
}

func resourceInstanceArrayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.InstanceArrayDelete(id)
	d.SetId("")
	return diags
}

func interfaceToString(v interface{}) string {
	var buf bytes.Buffer

	i := v.(map[string]interface{})

	// instance_array_interface_label := i["instance_array_interface_label"].(string)
	// instance_array_interface_service_status := i["instance_array_interface_service_status"].(string)
	instance_array_interface_index := strconv.Itoa(i["interface_index"].(int))
	network_id := strconv.Itoa(i["network_id"].(int))
	network_label := i["network_label"].(string)

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(network_label)))
	// buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_interface_service_status)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_interface_index)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(network_id)))

	return buf.String()
}

func interfaceResourceHash(v interface{}) int {
	return hash(interfaceToString(v))
}

func firewallRuleResourceHash(v interface{}) int {
	var buf bytes.Buffer
	fr := v.(map[string]interface{})

	firewall_rule_description := fr["firewall_rule_description"].(string)
	firewall_rule_source_ip_address_range_start := fr["firewall_rule_source_ip_address_range_start"].(string)
	firewall_rule_source_ip_address_range_end := fr["firewall_rule_source_ip_address_range_end"].(string)
	firewall_rule_destination_ip_address_range_start := fr["firewall_rule_destination_ip_address_range_start"].(string)
	firewall_rule_destination_ip_address_range_end := fr["firewall_rule_destination_ip_address_range_end"].(string)
	firewall_rule_protocol := fr["firewall_rule_protocol"].(string)
	firewall_rule_ip_address_type := fr["firewall_rule_ip_address_type"].(string)
	firewall_rule_port_range_start := strconv.Itoa(fr["firewall_rule_port_range_start"].(int))
	firewall_rule_port_range_end := strconv.Itoa(fr["firewall_rule_port_range_end"].(int))
	firewall_rule_enabled := strconv.FormatBool(fr["firewall_rule_enabled"].(bool))

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_description)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_source_ip_address_range_start)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_source_ip_address_range_end)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_destination_ip_address_range_start)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_destination_ip_address_range_end)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_protocol)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_ip_address_type)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_port_range_start)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_port_range_end)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_enabled)))

	return hash(buf.String())
}

func hash(v string) int {
	hash := crc32.ChecksumIEEE([]byte(v))

	return int(hash)

}
