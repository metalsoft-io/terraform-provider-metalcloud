package metalcloud

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
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
			"infrastructure_id": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == 0 {
						errs = append(errs, fmt.Errorf("%q is required. Provided value: %d", key, v))
					}
					return
				},
			},
			"instance_array_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_array_label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				Computed: true,
				//this is required because on the serverside the labels are converted to lowercase automatically
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.EqualFold(old, new) {
						return true
					}

					if new == "" {
						return true
					}
					return false
				},
				ValidateDiagFunc: validateLabel,
			},
			"instance_array_instance_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
				// ValidateDiagFunc: validateMaxOne,
			},
			"instance_array_boot_method": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:  "pxe_iscsi",
				Default:  nil,
				Computed: true, //default is computed serverside
			},
			"instance_array_ram_gbytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_processor_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_processor_core_mhz": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_processor_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_disk_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_disk_size_mbytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_additional_wan_ipv4_json": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"instance_array_custom_variables": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
				Computed: true, //default is computed serverside
			},
			"instance_custom_variables": {
				Type:     schema.TypeList,
				Elem:     resourceInstanceCustomVariable(),
				Optional: true,
			},
			"volume_template_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				//Computed: true, //default is computed serverside
			},
			"instance_array_firewall_managed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"firewall_rule": {
				Type:     schema.TypeSet,
				Optional: true, //default is computed serverside
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
				Elem:     resourceFirewallRule(),
			},
			"interface": {
				Type:     schema.TypeSet,
				Optional: true,
				Default:  nil,
				Computed: true,
				Elem:     resourceInstanceArrayInterface(),
			},
			"network_profile": {
				Type:     schema.TypeSet,
				Optional: true,
				Default:  nil,
				Computed: true,
				Elem:     resourceInstanceArrayNetworkProfile(),
			},
			"instance_server_type": {
				Type:     schema.TypeList,
				Elem:     resourceInstanceServerType(),
				Optional: true,
			},

			"drive_array_id_boot": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				//Computed: true,
			},
			"cluster_group_role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				//Computed: true,
			},
		},
	}
}

func resourceInstanceArrayNetworkProfile() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network_profile_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceInstanceCustomVariable() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_index": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"custom_variables": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceInstanceServerType() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_index": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server_type_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"firewall_rule_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"firewall_rule_port_range_start": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:  1,
			},
			"firewall_rule_port_range_end": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:  65535,
			},
			"firewall_rule_source_ip_address_range_start": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"firewall_rule_source_ip_address_range_end": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"firewall_rule_destination_ip_address_range_start": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"firewall_rule_destination_ip_address_range_end": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"firewall_rule_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tcp",
			},
			"firewall_rule_ip_address_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ipv4",
			},
			"firewall_rule_enabled": {
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
			"network_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceInstanceArrayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)

	var diags diag.Diagnostics

	infrastructure_id := d.Get("infrastructure_id").(int)
	_, err := client.InfrastructureGet(infrastructure_id)

	if err != nil {
		return diag.Errorf("Infrastructure with id %+v not found.", infrastructure_id)
	}
	ia := expandInstanceArray(d)

	if ia.VolumeTemplateID != 0 && ia.DriveArrayIDBoot == 0 && ia.InstanceArrayBootMethod != LOCAL_DRIVES {
		return diag.Errorf("Current instance array configuration is only valid for %s boot method", LOCAL_DRIVES)
	}

	if ia.DriveArrayIDBoot != 0 && ia.InstanceArrayBootMethod != PXE_ISCSI {
		return diag.Errorf("Current instance array configuration is only valid for %s boot method", PXE_ISCSI)
	}

	iaC, err := client.InstanceArrayCreate(infrastructure_id, ia)
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%d", iaC.InstanceArrayID)

	d.SetId(id)

	/* custom variables for instances */
	cvList := d.Get("instance_custom_variables").([]interface{})
	dg := updateInstancesCustomVariables(cvList, iaC.InstanceArrayID, client)

	if dg.HasError() {
		resourceInstanceArrayRead(ctx, d, meta)
		return dg
	}

	diags = append(diags, dg...)

	/* update server types */
	iList := d.Get("instance_server_type").([]interface{})
	dg = updateInstancesServerTypes(iList, iaC.InstanceArrayID, client)

	if dg.HasError() {
		resourceInstanceArrayRead(ctx, d, meta)
		return dg
	}

	diags = append(diags, dg...)

	for _, intf := range ia.InstanceArrayInterfaces {
		_, err := client.InstanceArrayInterfaceAttachNetwork(iaC.InstanceArrayID, intf.InstanceArrayInterfaceIndex, intf.NetworkID)
		if err != nil {
			resourceInstanceArrayRead(ctx, d, meta)
			return diag.FromErr(err)
		}
	}

	if d.Get("network_profile") != nil {
		profileSet := d.Get("network_profile").(*schema.Set)

		for _, profileIntf := range profileSet.List() {
			profileMap := profileIntf.(map[string]interface{})

			_, err := client.InstanceArrayNetworkProfileSet(iaC.InstanceArrayID, profileMap["network_id"].(int), profileMap["network_profile_id"].(int))

			if err != nil {
				resourceInstanceArrayRead(ctx, d, meta)
				return diag.FromErr(err)
			}
		}
	}

	dg = resourceInstanceArrayRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags
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

	flattenInstanceArray(d, *ia)

	networkToNetworkProfileMap, err := client.NetworkProfileListByInstanceArray(ia.InstanceArrayID)

	if err != nil {
		return diag.FromErr(err)
	}

	var networkProfileList []interface{}
	for networkID, networkProfileID := range *networkToNetworkProfileMap {
		if networkProfileID > 0 {
			networkProfileEntry := flattenInstanceArrayNetworkProfile(networkID, networkProfileID)
			networkProfileList = append(networkProfileList, networkProfileEntry)
		}
	}

	if len(networkProfileList) > 0 {
		d.Set("network_profile", schema.NewSet(schema.HashResource(resourceInstanceArrayNetworkProfile()), networkProfileList))
	}

	/* INSTANCES */
	retInstances, err := client.InstanceArrayInstances(ia.InstanceArrayID)
	if err != nil {
		return diag.FromErr(err)
	}

	//instances
	keys := []int{}
	instances := []mc.Instance{}

	for _, v := range *retInstances {

		keys = append(keys, v.InstanceID)
	}

	sort.Ints(keys)

	for _, id := range keys {
		i, err := client.InstanceGet(id)
		if err != nil {
			return diag.FromErr(err)
		}
		instances = append(instances, *i)
	}

	/* INSTANCES CUSTOM VARS */
	instancesCustomVariables := flattenInstancesCustomVariables(retInstances)

	if len(instancesCustomVariables) > 0 {
		d.Set("instance_custom_variables", instancesCustomVariables)

	}

	return diags
}

func resourceInstanceArrayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	retIA, err := client.InstanceArrayGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ia := expandInstanceArray(d)

	if ia.VolumeTemplateID != 0 && ia.DriveArrayIDBoot == 0 && ia.InstanceArrayBootMethod != LOCAL_DRIVES {
		return diag.Errorf("Current instance array configuration is only valid for '%s' boot method", LOCAL_DRIVES)
	}

	if ia.DriveArrayIDBoot != 0 && ia.InstanceArrayBootMethod != PXE_ISCSI {
		return diag.Errorf("Current instance array configuration is only valid for '%s' boot method", PXE_ISCSI)
	}

	//update interface operations
	for _, intf := range ia.InstanceArrayInterfaces {
		for _, opIntf := range retIA.InstanceArrayOperation.InstanceArrayInterfaces {
			if opIntf.InstanceArrayInterfaceID == intf.InstanceArrayInterfaceID {
				copyInstanceArrayInterfaceToOperation(intf, &opIntf)
			}
		}
	}

	//update the main operation object
	copyInstanceArrayToOperation(ia, retIA.InstanceArrayOperation)

	bSwapExistingInstancesHardware := false
	bkeepDetachingDrives := false

	editedIA, err := client.InstanceArrayEdit(id, *retIA.InstanceArrayOperation, &bSwapExistingInstancesHardware, &bkeepDetachingDrives, nil, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var dg diag.Diagnostics

	/* custom variables for instances */
	cvList := d.Get("instance_custom_variables").([]interface{})
	if d.HasChange("instance_custom_variables") {
		dg = updateInstancesCustomVariables(cvList, id, client)

		if dg.HasError() {
			resourceInstanceArrayRead(ctx, d, meta)
			return dg
		}

		diags = append(diags, dg...)
	}

	/* update server types */
	iList := d.Get("instance_server_type").([]interface{})
	dg = updateInstancesServerTypes(iList, id, client)

	if dg.HasError() {
		resourceInstanceArrayRead(ctx, d, meta)
		return dg
	}

	diags = append(diags, dg...)

	//update interfaces
	if d.HasChange("interface") {
		dg := updateInstanceArrayInterfaces(ia.InstanceArrayInterfaces, editedIA.InstanceArrayInterfaces, client)

		if dg.HasError() {
			resourceInstanceArrayRead(ctx, d, meta)
			return dg
		}
	}

	if d.Get("network_profile") != nil {
		profileSet := d.Get("network_profile").(*schema.Set)

		for _, profileIntf := range profileSet.List() {
			profileMap := profileIntf.(map[string]interface{})
			if profileMap["network_profile_id"].(int) != 0 {
				_, err := client.InstanceArrayNetworkProfileSet(id, profileMap["network_id"].(int), profileMap["network_profile_id"].(int))

				if err != nil {
					resourceInstanceArrayRead(ctx, d, meta)
					return diag.FromErr(err)
				}
			} else {
				err := client.InstanceArrayNetworkProfileClear(id, profileMap["network_id"].(int))

				if err != nil {
					resourceInstanceArrayRead(ctx, d, meta)
					return diag.FromErr(err)
				}
			}
		}
	}

	dg = resourceInstanceArrayRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags
}

func updateInstanceArrayInterfaces(configuredInterfaces []mc.InstanceArrayInterface, intfList []mc.InstanceArrayInterface, client *mc.Client) diag.Diagnostics {
	for _, intf := range configuredInterfaces {
		_, err := client.InstanceArrayInterfaceAttachNetwork(intf.InstanceArrayID, intf.InstanceArrayInterfaceIndex, intf.NetworkID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	for _, existingIntf := range intfList {
		found := false

		for _, intf := range configuredInterfaces {
			if existingIntf.InstanceArrayInterfaceIndex == intf.InstanceArrayInterfaceIndex {
				found = true
			}
		}
		if found == false && existingIntf.NetworkID != 0 {
			_, err := client.InstanceArrayInterfaceDetach(existingIntf.InstanceArrayID, existingIntf.InstanceArrayInterfaceIndex)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func resourceInstanceArrayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ia, err := client.InstanceArrayGet(id)

	if err == nil && ia.InstanceArrayServiceStatus != SERVICE_STATUS_DELETED {
		err = client.InstanceArrayDelete(id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return diags
}

func flattenInstanceArray(d *schema.ResourceData, instanceArray mc.InstanceArray) error {

	d.Set("instance_array_id", instanceArray.InstanceArrayID)
	d.Set("instance_array_label", instanceArray.InstanceArrayLabel)
	d.Set("instance_array_boot_method", instanceArray.InstanceArrayBootMethod)
	d.Set("instance_array_ram_gbytes", instanceArray.InstanceArrayRAMGbytes)
	d.Set("instance_array_processor_count", instanceArray.InstanceArrayProcessorCount)
	d.Set("instance_array_processor_core_mhz", instanceArray.InstanceArrayProcessorCoreMHZ)
	d.Set("instance_array_processor_core_count", instanceArray.InstanceArrayProcessorCoreCount)
	d.Set("instance_array_disk_count", instanceArray.InstanceArrayDiskCount)
	d.Set("instance_array_disk_size_mbytes", instanceArray.InstanceArrayDiskSizeMBytes)
	d.Set("volume_template_id", instanceArray.VolumeTemplateID)
	d.Set("instance_array_firewall_managed", instanceArray.InstanceArrayFirewallManaged)
	d.Set("instance_array_additional_wan_ipv4_json", instanceArray.InstanceArrayAdditionalWanIPv4JSON)
	d.Set("infrastructure_id", instanceArray.InfrastructureID)
	d.Set("drive_array_id_boot", instanceArray.DriveArrayIDBoot)
	d.Set("instance_array_instance_count", instanceArray.InstanceArrayOperation.InstanceArrayInstanceCount)

	/* INSTANCE ARRAY CUSTOM VARIABLES */
	switch instanceArray.InstanceArrayCustomVariables.(type) {
	case []interface{}:
		d.Set("instance_array_custom_variables", make(map[string]string))
	default:
		iacv := make(map[string]string)

		for k, v := range instanceArray.InstanceArrayCustomVariables.(map[string]interface{}) {
			iacv[k] = v.(string)
		}
		d.Set("instance_array_custom_variables", iacv)
	}

	/* FIREWALL RULES */
	fwRules := []interface{}{}

	for _, fw := range instanceArray.InstanceArrayFirewallRules {
		fwRules = append(fwRules, flattenFirewallRule(fw))
	}
	if len(fwRules) > 0 {
		d.Set("firewall_rule", schema.NewSet(schema.HashResource(resourceFirewallRule()), fwRules))
	}

	/* INTERFACES */
	interfaces := []interface{}{}

	//iterate over interfaces
	if intfList, ok := d.GetOkExists("interface"); ok {
		for _, iIntf := range intfList.(*schema.Set).List() {
			iaInterface := iIntf.(map[string]interface{})
			interfaceIndex := iaInterface["interface_index"].(int)

			//locate interface with index in returned data
			for _, intf := range instanceArray.InstanceArrayInterfaces {
				//if we found it, locate the network it's connected to add it to the list
				if intf.InstanceArrayInterfaceIndex == interfaceIndex && intf.NetworkID != 0 {
					interfaces = append(interfaces, flattenInstanceArrayInterface(intf))
				}
			}
		}
	}
	if len(interfaces) > 0 {
		d.Set("interface", schema.NewSet(schema.HashResource(resourceInstanceArrayInterface()), interfaces))
	}

	return nil
}

func expandInstanceArray(d *schema.ResourceData) mc.InstanceArray {

	var ia mc.InstanceArray
	if d.Get("instance_array_id") != nil {
		ia.InstanceArrayID = d.Get("instance_array_id").(int)
	}
	ia.InstanceArrayLabel = d.Get("instance_array_label").(string)
	ia.InstanceArrayInstanceCount = d.Get("instance_array_instance_count").(int)

	//ia.InstanceArraySubdomain = d.Get("instance_array_subdomain").(string)

	if d.Get("instance_array_boot_method") != nil {
		ia.InstanceArrayBootMethod = d.Get("instance_array_boot_method").(string)
	} else {
		ia.InstanceArrayBootMethod = LOCAL_DRIVES
	}
	ia.InstanceArrayRAMGbytes = d.Get("instance_array_ram_gbytes").(int)
	ia.InstanceArrayProcessorCount = d.Get("instance_array_processor_count").(int)
	ia.InstanceArrayProcessorCoreMHZ = d.Get("instance_array_processor_core_mhz").(int)
	ia.InstanceArrayProcessorCoreCount = d.Get("instance_array_processor_core_count").(int)
	ia.InstanceArrayDiskCount = d.Get("instance_array_disk_count").(int)
	ia.InstanceArrayDiskSizeMBytes = d.Get("instance_array_disk_size_mbytes").(int)
	ia.VolumeTemplateID = d.Get("volume_template_id").(int)
	ia.InstanceArrayAdditionalWanIPv4JSON = d.Get("instance_array_additional_wan_ipv4_json").(string)
	ia.DriveArrayIDBoot = d.Get("drive_array_id_boot").(int)
	ia.InstanceArrayFirewallManaged = d.Get("instance_array_firewall_managed").(bool)

	if d.Get("firewall_rule") != nil {
		fwRulesSet := d.Get("firewall_rule").(*schema.Set)
		fwRules := []mc.FirewallRule{}

		for _, fwMap := range fwRulesSet.List() {
			fwRules = append(fwRules, expandFirewallRule(fwMap.(map[string]interface{})))
		}

		ia.InstanceArrayFirewallRules = fwRules
	}

	if d.Get("interface") != nil {
		interfaceSet := d.Get("interface").(*schema.Set)
		interfaces := []mc.InstanceArrayInterface{}

		for _, intfList := range interfaceSet.List() {
			intfMap := intfList.(map[string]interface{})
			intfMap["instance_array_id"] = ia.InstanceArrayID
			interfaces = append(interfaces, expandInstanceArrayInterface(intfMap))
		}

		ia.InstanceArrayInterfaces = interfaces
	}

	if d.Get("instance_array_custom_variables") != nil {
		iacv := make(map[string]string)

		for k, v := range d.Get("instance_array_custom_variables").(map[string]interface{}) {
			iacv[k] = v.(string)
		}

		ia.InstanceArrayCustomVariables = iacv
	}

	return ia
}

func flattenFirewallRule(fw mc.FirewallRule) map[string]interface{} {
	var d = make(map[string]interface{})

	d["firewall_rule_description"] = fw.FirewallRuleDescription
	d["firewall_rule_port_range_start"] = fw.FirewallRulePortRangeStart
	d["firewall_rule_port_range_end"] = fw.FirewallRulePortRangeEnd
	d["firewall_rule_source_ip_address_range_start"] = fw.FirewallRuleSourceIPAddressRangeStart
	d["firewall_rule_source_ip_address_range_end"] = fw.FirewallRuleSourceIPAddressRangeEnd
	d["firewall_rule_destination_ip_address_range_start"] = fw.FirewallRuleDestinationIPAddressRangeStart
	d["firewall_rule_destination_ip_address_range_end"] = fw.FirewallRuleDestinationIPAddressRangeEnd
	d["firewall_rule_protocol"] = fw.FirewallRuleProtocol
	d["firewall_rule_ip_address_type"] = fw.FirewallRuleIPAddressType
	d["firewall_rule_enabled"] = fw.FirewallRuleEnabled

	return d
}

func expandFirewallRule(d map[string]interface{}) mc.FirewallRule {
	var fw mc.FirewallRule

	fw.FirewallRuleDescription = d["firewall_rule_description"].(string)
	fw.FirewallRulePortRangeStart = d["firewall_rule_port_range_start"].(int)
	fw.FirewallRulePortRangeEnd = d["firewall_rule_port_range_end"].(int)
	fw.FirewallRuleSourceIPAddressRangeStart = d["firewall_rule_source_ip_address_range_start"].(string)
	fw.FirewallRuleSourceIPAddressRangeEnd = d["firewall_rule_source_ip_address_range_end"].(string)
	fw.FirewallRuleDestinationIPAddressRangeStart = d["firewall_rule_destination_ip_address_range_start"].(string)
	fw.FirewallRuleDestinationIPAddressRangeEnd = d["firewall_rule_destination_ip_address_range_end"].(string)
	fw.FirewallRuleProtocol = d["firewall_rule_protocol"].(string)
	fw.FirewallRuleIPAddressType = d["firewall_rule_ip_address_type"].(string)
	fw.FirewallRuleEnabled = d["firewall_rule_enabled"].(bool)

	return fw
}

func flattenInstanceArrayInterface(i mc.InstanceArrayInterface) map[string]interface{} {
	var d = make(map[string]interface{})

	d["interface_index"] = i.InstanceArrayInterfaceIndex
	d["network_id"] = i.NetworkID

	return d
}

func expandInstanceArrayInterface(d map[string]interface{}) mc.InstanceArrayInterface {

	var i mc.InstanceArrayInterface

	i.InstanceArrayInterfaceIndex = d["interface_index"].(int)
	i.NetworkID = d["network_id"].(int)
	i.InstanceArrayID = d["instance_array_id"].(int)

	return i
}

/*
//I am unsure why this was implemented this way
func flattenInstanceArrayNetworkProfile(networkProfiles map[int]int, d *schema.ResourceData) []interface{} {
	profiles := []interface{}{}

	if profileList, ok := d.GetOkExists("network_profile"); ok {
		for _, pIntf := range profileList.(*schema.Set).List() {
			profile := pIntf.(map[string]interface{})
			network_id := profile["network_id"].(int)

			if _, ok := networkProfiles[network_id]; ok {
				profiles = append(profiles, profile)
			}
		}
	}

	return profiles
}
*/

func flattenInstanceArrayNetworkProfile(networkID int, networkProfileID int) map[string]interface{} {
	var d = make(map[string]interface{})

	d["network_profile_id"] = networkProfileID
	d["network_id"] = networkID

	return d
}

func flattenInstancesCustomVariables(retInstances *map[string]mc.Instance) []interface{} {

	instanceMap := make(map[int]mc.Instance, len(*retInstances))
	keys := []int{}
	instances := []mc.Instance{}

	for _, v := range *retInstances {
		instanceMap[v.InstanceID] = v
		keys = append(keys, v.InstanceID)
	}

	sort.Ints(keys)

	for _, id := range keys {
		instances = append(instances, instanceMap[id])
	}

	customVars := []interface{}{}

	for index, instance := range instances {
		i := make(map[string]interface{})
		cv := make(map[string]interface{})
		i["instance_index"] = index
		switch instance.InstanceCustomVariables.(type) {
		//todo: add nil
		case []interface{}:
			cv = make(map[string]interface{})
		default:
			for k, v := range instance.InstanceCustomVariables.(map[string]interface{}) {
				cv[k] = v.(string)
			}
		}
		i["custom_variables"] = cv
		if len(cv) > 0 {
			customVars = append(customVars, i)
		}
	}

	return customVars
}

func flattenInstanceServerTypes(retInstances *map[string]mc.Instance) []interface{} {
	instanceServerTypes := []interface{}{}
	instances := getSortedListOfInstances(retInstances)

	for instanceIndex, inst := range instances {
		entry := make(map[string]interface{})
		entry["instance_index"] = instanceIndex
		entry["server_type_id"] = inst.ServerTypeID
		instanceServerTypes = append(instanceServerTypes, entry)
	}
	return instanceServerTypes
}

// this is used to get a sorted list of instances ordered by their id
// so that instance_index is consistent
func getSortedListOfInstances(instanceList *map[string]mc.Instance) []mc.Instance {
	instanceMap := make(map[int]mc.Instance, len(*instanceList))

	keys := []int{}
	instances := []mc.Instance{}

	for _, v := range *instanceList {
		instanceMap[v.InstanceID] = v
		keys = append(keys, v.InstanceID)
	}

	sort.Ints(keys)

	for _, id := range keys {
		instances = append(instances, instanceMap[id])
	}
	return instances
}

// * sets the custom variables on the instances object. Used by the Upgrade function
// TODO: convert tot an actual expand function that doesn't use the client to set them to make it easier to test
func updateInstancesCustomVariables(cvList []interface{}, instanceArrayID int, client *mc.Client) diag.Diagnostics {

	var diags diag.Diagnostics
	instanceList, err := client.InstanceArrayInstances(instanceArrayID)
	if err != nil {
		return diag.FromErr(err)
	}

	instances := getSortedListOfInstances(instanceList)

	nInstances := len(*instanceList)

	currentCVLabelList := make(map[string]int, len(*instanceList))

	for _, icvIntf := range cvList {
		icv := icvIntf.(map[string]interface{})
		cvIntf := icv["custom_variables"].(map[string]interface{})
		instance_custom_variables := make(map[string]string)
		for k, v := range cvIntf {
			instance_custom_variables[k] = v.(string)
		}

		instance_index := icv["instance_index"].(int)

		if instance_index < 0 || instance_index >= nInstances {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "instance_index in custom_variables block out of bounds",
				Detail:   fmt.Sprintf("Use a number between 0 and %v (instance_array_instance_count-1). ", nInstances-1),
			})
		}

		instance := instances[instance_index]
		currentCVLabelList[instance.InstanceLabel] = instance.InstanceID
		instance.InstanceOperation.InstanceCustomVariables = instance_custom_variables
		_, err := client.InstanceEdit(instance.InstanceID, instance.InstanceOperation)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	for _, instance := range *instanceList {
		if _, ok := currentCVLabelList[instance.InstanceLabel]; !ok {
			instance.InstanceOperation.InstanceCustomVariables = make(map[string]string)
			_, err := client.InstanceEdit(instance.InstanceID, instance.InstanceOperation)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return diags
}

// * sets the server types on each of the instances
func updateInstancesServerTypes(iList []interface{}, instanceArrayID int, client *mc.Client) diag.Diagnostics {

	var diags diag.Diagnostics
	instanceList, err := client.InstanceArrayInstances(instanceArrayID)
	if err != nil {
		return diag.FromErr(err)
	}
	//TODO: flatten instances
	instanceMap := make(map[int]mc.Instance, len(*instanceList))
	nInstances := len(*instanceList)
	keys := []int{}
	instances := []mc.Instance{}

	for _, v := range *instanceList {
		instanceMap[v.InstanceID] = v
		keys = append(keys, v.InstanceID)
	}

	sort.Ints(keys)

	for _, id := range keys {
		instances = append(instances, instanceMap[id])
	}

	for _, iIntf := range iList {
		imap := iIntf.(map[string]interface{})

		instance_index := imap["instance_index"].(int)
		server_type_id := imap["server_type_id"].(int)

		if instance_index < 0 || instance_index >= nInstances {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("instance_index in instance_server_type block out of bounds"),
				Detail:   fmt.Sprintf("Use a number between 0 and %v (instance_array_instance_count-1). ", nInstances-1),
			})
		}

		instance := instances[instance_index]

		instance.InstanceOperation.ServerTypeID = server_type_id

		_, err := client.InstanceEdit(instance.InstanceID, instance.InstanceOperation)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	return diags
}

func copyInstanceArrayToOperation(ia mc.InstanceArray, iao *mc.InstanceArrayOperation) {

	iao.InstanceArrayID = ia.InstanceArrayID
	iao.InstanceArrayLabel = ia.InstanceArrayLabel
	iao.InstanceArrayBootMethod = ia.InstanceArrayBootMethod
	iao.InstanceArrayInstanceCount = ia.InstanceArrayInstanceCount
	iao.InstanceArrayRAMGbytes = ia.InstanceArrayRAMGbytes
	iao.InstanceArrayProcessorCount = ia.InstanceArrayProcessorCount
	iao.InstanceArrayProcessorCoreMHZ = ia.InstanceArrayProcessorCoreMHZ
	iao.InstanceArrayDiskCount = ia.InstanceArrayDiskCount
	iao.InstanceArrayDiskSizeMBytes = ia.InstanceArrayDiskSizeMBytes
	iao.InstanceArrayDiskTypes = ia.InstanceArrayDiskTypes
	//iao.ClusterID = ia.ClusterID
	iao.InstanceArrayFirewallManaged = ia.InstanceArrayFirewallManaged
	iao.InstanceArrayFirewallRules = ia.InstanceArrayFirewallRules
	iao.VolumeTemplateID = ia.VolumeTemplateID
	iao.InstanceArrayAdditionalWanIPv4JSON = ia.InstanceArrayAdditionalWanIPv4JSON
	iao.InstanceArrayCustomVariables = ia.InstanceArrayCustomVariables
}

func copyInstanceArrayInterfaceToOperation(i mc.InstanceArrayInterface, io *mc.InstanceArrayInterfaceOperation) {
	io.InstanceArrayInterfaceLAGGIndexes = i.InstanceArrayInterfaceLAGGIndexes
	io.InstanceArrayInterfaceIndex = i.InstanceArrayInterfaceIndex
	io.NetworkID = i.NetworkID
	io.InstanceArrayInterfaceChangeID = i.InstanceArrayInterfaceChangeID
}

const PXE_ISCSI = "pxe_iscsi"
const LOCAL_DRIVES = "local_drives"
