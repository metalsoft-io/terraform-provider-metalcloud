package metalcloud

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
)

func resourceVmInstanceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVmInstanceGroupCreate,
		ReadContext:   resourceVmInstanceGroupRead,
		UpdateContext: resourceVmInstanceGroupUpdate,
		DeleteContext: resourceVmInstanceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"infrastructure_id": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validateRequired,
			},
			"vm_instance_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vm_instance_group_label": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          nil,
				Computed:         true,
				DiffSuppressFunc: caseInsensitiveDiff,
				ValidateDiagFunc: validateLabel,
			},
			"vm_instance_group_instance_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"vm_type_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vm_instance_group_disk_size_gbytes": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"vm_instance_group_template_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"vm_instance_custom_variables": {
				Type:     schema.TypeList,
				Elem:     resourceVmInstanceCustomVariable(),
				Optional: true,
			},
			"vm_instance_group_interfaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Default:  nil,
				Computed: true,
				Elem:     resourceVmInstanceGroupInterface(),
			},
			"vm_instance_group_network_profiles": {
				Type:     schema.TypeSet,
				Optional: true,
				Default:  nil,
				Computed: true,
				Elem:     resourceVmInstanceGroupNetworkProfile(),
			},
		},
	}
}

func resourceVmInstanceCustomVariable() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vm_instance_index": {
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

func resourceVmInstanceGroupInterface() *schema.Resource {
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

func resourceVmInstanceGroupNetworkProfile() *schema.Resource {
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

func resourceVmInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	infrastructureId := getInfrastructureId(d)

	client, err := getAPIClient()
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = client.InfrastructureApi.GetInfrastructure(ctx, float64(infrastructureId))
	if err != nil {
		return diag.Errorf("Infrastructure with Id %+v not found.", infrastructureId)
	}

	vmGroupCreate, _ := expandCreateVmInstanceGroup(d)

	vmInstanceGroupCreated, _, err := client.VMInstanceGroupApi.CreateVMInstanceGroup(ctx, vmGroupCreate, float64(infrastructureId))
	if err != nil {
		return extractApiError(err)
	}

	vmInstanceGroupId := int(vmInstanceGroupCreated.Id)

	d.SetId(fmt.Sprintf("%d", vmInstanceGroupId))

	// VM instances custom variables
	dg := updateVmInstancesCustomVariables(ctx, client, infrastructureId, vmInstanceGroupId, d.Get("vm_instance_custom_variables").([]interface{}))
	if dg.HasError() {
		resourceVmInstanceGroupRead(ctx, d, meta)
		return dg
	}
	diags = append(diags, dg...)

	// VM instance group interfaces
	if d.Get("vm_instance_group_interfaces") != nil {
		interfacesSet := d.Get("vm_instance_group_interfaces").(*schema.Set)

		for _, interfaceEntry := range interfacesSet.List() {
			interfaceEntryMap := interfaceEntry.(map[string]interface{})
			vmInterfaceCreate := sdk2.CreateVmInstanceGroupInterface{
				NetworkId: float64(interfaceEntryMap["network_id"].(int)),
			}

			_, _, err := client.VMInstanceGroupApi.CreateVMInterfaceOnVMInstanceGroup(
				ctx,
				vmInterfaceCreate,
				float64(infrastructureId),
				float64(vmInstanceGroupId),
			)
			if err != nil {
				resourceVmInstanceGroupRead(ctx, d, meta)
				return diag.FromErr(err)
			}
		}
	}

	// VM instance group network profiles
	if d.Get("vm_instance_group_network_profiles") != nil {
		networkProfilesSet := d.Get("vm_instance_group_network_profiles").(*schema.Set)

		for _, networkProfileEntry := range networkProfilesSet.List() {
			networkProfileEntryMap := networkProfileEntry.(map[string]interface{})
			networkId := networkProfileEntryMap["network_id"].(int)
			networkProfileUpdate := sdk2.UpdateVmInstanceGroupNetwork{
				NetworkProfileId: float64(networkProfileEntryMap["network_profile_id"].(int)),
			}

			_, _, err := client.VMInstanceGroupApi.UpdateNetworkProfileOnVMInstanceGroupNetwork(
				ctx,
				networkProfileUpdate,
				float64(infrastructureId),
				float64(vmInstanceGroupId),
				float64(networkId),
			)
			if err != nil {
				resourceVmInstanceGroupRead(ctx, d, meta)
				return diag.FromErr(err)
			}
		}
	}

	dg = resourceVmInstanceGroupRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags
}

func resourceVmInstanceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vmInstanceGroupId, infrastructureId, err := getVmInstanceGroupId(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := getAPIClient()
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstanceGroup, _, err := client.VMInstanceGroupApi.GetVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	flattenVmInstanceGroup(d, vmInstanceGroup)

	return diags
}

func resourceVmInstanceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vmInstanceGroupId, infrastructureId, err := getVmInstanceGroupId(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := getAPIClient()
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = client.VMInstanceGroupApi.GetVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	updateVmInstanceGroup, _ := expandUpdateVmInstanceGroup(d)

	_, _, err = client.VMInstanceGroupApi.UpdateVMInstanceGroup(ctx, updateVmInstanceGroup, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	var dg diag.Diagnostics

	// VM instance group network profiles
	if d.Get("vm_instance_group_network_profiles") != nil {
		networkProfilesSet := d.Get("vm_instance_group_network_profiles").(*schema.Set)

		for _, networkProfileEntry := range networkProfilesSet.List() {
			networkProfileEntryMap := networkProfileEntry.(map[string]interface{})
			networkId := networkProfileEntryMap["network_id"].(int)
			networkProfileUpdate := sdk2.UpdateVmInstanceGroupNetwork{
				NetworkProfileId: float64(networkProfileEntryMap["network_profile_id"].(int)),
			}

			_, _, err := client.VMInstanceGroupApi.UpdateNetworkProfileOnVMInstanceGroupNetwork(
				ctx,
				networkProfileUpdate,
				float64(infrastructureId),
				float64(vmInstanceGroupId),
				float64(networkId),
			)
			if err != nil {
				resourceVmInstanceGroupRead(ctx, d, meta)
				return diag.FromErr(err)
			}
		}
	}

	// /* update VM types */
	// iList := d.Get("vm_instance_type").([]interface{})
	// dg = updateVmInstancesTypes(iList, vmInstanceGroupId, client2)

	// if dg.HasError() {
	// 	resourceVmInstanceGroupRead(ctx, d, meta)
	// 	return dg
	// }

	// diags = append(diags, dg...)

	dg = resourceVmInstanceGroupRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags
}

func resourceVmInstanceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vmInstanceGroupId, infrastructureId, err := getVmInstanceGroupId(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := getAPIClient()
	if err != nil {
		return diag.FromErr(err)
	}

	oldVmInstanceGroup, _, err := client.VMInstanceGroupApi.GetVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	if oldVmInstanceGroup.ServiceStatus != SERVICE_STATUS_DELETED {
		_, err = client.VMInstanceGroupApi.DeleteVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
		if err != nil {
			return extractApiError(err)
		}
	}

	d.SetId("")

	return diags
}

// sets the custom variables on the VM instances
func updateVmInstancesCustomVariables(ctx context.Context, client2 *sdk2.APIClient, infrastructureId int, vmInstanceGroupId int, cvList []interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	vmInstanceList, _, err := client2.VMInstanceGroupApi.GetVMInstanceGroupVMInstances(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstances := getSortedListOfVmInstances(vmInstanceList)

	nVmInstances := len(vmInstanceList)

	currentCVLabelList := make(map[string]int, len(vmInstanceList))

	for _, cvEntry := range cvList {
		cvEntryMap := cvEntry.(map[string]interface{})
		cvEntryVariables := cvEntryMap["custom_variables"].(map[string]interface{})

		customVariables := map[string]string{}
		for k, v := range cvEntryVariables {
			customVariables[k] = v.(string)
		}

		instanceIndex := cvEntryMap["vm_instance_index"].(int)
		if instanceIndex < 0 || instanceIndex >= nVmInstances {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "vm_instance_index in custom_variables block out of bounds",
				Detail:   fmt.Sprintf("Use a number between 0 and %v (vm_instance_group_instance_count-1). ", nVmInstances-1),
			})
		}

		instance := vmInstances[instanceIndex]
		currentCVLabelList[instance.Label] = int(instance.Id)

		var customVariablesValue interface{} = customVariables
		vmInstanceUpdate := sdk2.UpdateVmInstance{
			CustomVariables: &customVariablesValue,
		}

		_, _, err := client2.VMInstanceApi.UpdateVMInstance(ctx, vmInstanceUpdate, float64(infrastructureId), instance.Id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	for _, vmInstance := range vmInstanceList {
		// Remove the custom variables from the instance if they are not in the list
		if _, ok := currentCVLabelList[vmInstance.Label]; !ok {
			vmInstanceUpdate := sdk2.UpdateVmInstance{
				CustomVariables: nil,
			}

			_, _, err := client2.VMInstanceApi.UpdateVMInstance(ctx, vmInstanceUpdate, float64(infrastructureId), vmInstance.Id)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return diags
}

func getInfrastructureId(d *schema.ResourceData) int {
	return d.Get("infrastructure_id").(int)
}

func getVmInstanceGroupId(d *schema.ResourceData) (vmInstanceGroupId int, infrastructureId int, err error) {
	vmInstanceGroupId, err = strconv.Atoi(d.Id())
	if err != nil {
		return
	}

	infrastructureId = getInfrastructureId(d)

	return
}

// Returns sorted list of VM instances ordered by their id, so that vm_instance_index is consistent
func getSortedListOfVmInstances(vmInstanceList []sdk2.VmInstance) []sdk2.VmInstance {
	instanceMap := make(map[int]sdk2.VmInstance, len(vmInstanceList))

	keys := []int{}
	instances := []sdk2.VmInstance{}

	for _, v := range vmInstanceList {
		instanceMap[int(v.Id)] = v
		keys = append(keys, int(v.Id))
	}

	sort.Ints(keys)

	for _, id := range keys {
		instances = append(instances, instanceMap[id])
	}

	return instances
}

func flattenVmInstanceGroup(d *schema.ResourceData, vmInstanceGroup sdk2.VmInstanceGroup) error {
	d.Set("infrastructure_id", vmInstanceGroup.InfrastructureId)
	d.Set("vm_instance_group_id", vmInstanceGroup.Id)
	d.Set("vm_instance_group_label", vmInstanceGroup.Label)
	d.Set("vm_instance_group_instance_count", vmInstanceGroup.InstanceCount)
	d.Set("vm_instance_group_disk_size_gbytes", vmInstanceGroup.DiskSizeGB)
	d.Set("vm_instance_group_template_id", vmInstanceGroup.VolumeTemplateId)

	/* NETWORK PROFILES */
	if vmInstanceGroup.NetworkIdToNetworkProfileId != nil {
		networkIdToNetworkProfileId := *vmInstanceGroup.NetworkIdToNetworkProfileId
		networkIdToNetworkProfileIdMap := networkIdToNetworkProfileId.(map[string]string)
		networkProfiles := []interface{}{}
		for networkId, networkProfileId := range networkIdToNetworkProfileIdMap {
			networkProfiles = append(networkProfiles, map[string]interface{}{
				"network_id":         networkId,
				"network_profile_id": networkProfileId,
			})
		}
		d.Set("vm_instance_group_network_profiles", networkProfiles)
	}

	/* INTERFACES */
	vmInterfaces := []interface{}{}
	vmInterfacesSet, ok := d.GetOk("vm_instance_group_interfaces")
	if ok {
		for _, vmInterface := range vmInterfacesSet.(*schema.Set).List() {
			vmInterfaceMap := vmInterface.(map[string]interface{})
			interfaceIndex := vmInterfaceMap["interface_index"].(int)

			// locate interface with index in returned data
			for _, vmInterface := range vmInstanceGroup.VmInstanceGroupInterfaces {
				// if we found it, locate the network it's connected to add it to the list
				if int(vmInterface.InterfaceIndex) == interfaceIndex && vmInterface.NetworkId != 0 {
					vmInterfaces = append(vmInterfaces, flattenVmInstanceGroupInterface(vmInterface))
				}
			}
		}
	}

	if len(vmInterfaces) > 0 {
		d.Set("vm_instance_group_interfaces", schema.NewSet(schema.HashResource(resourceVmInstanceGroupInterface()), vmInterfaces))
	}

	/* CUSTOM VARIABLES */
	customVariables := *(vmInstanceGroup.CustomVariables)
	switch customVariables.(type) {
	case []interface{}:
		d.Set("vm_instance_custom_variables", make(map[string]string))

	default:
		cv := make(map[string]string)

		for k, v := range customVariables.(map[string]interface{}) {
			cv[k] = v.(string)
		}

		d.Set("vm_instance_custom_variables", cv)
	}

	// /* INSTANCES */
	// for _, vmInstance := range vmInstanceGroup.VmInstance {
	// }

	return nil
}

func flattenVmInstanceGroupInterface(i sdk2.VmInstanceGroupInterface) map[string]interface{} {
	var d = make(map[string]interface{})

	d["interface_index"] = int(i.InterfaceIndex)
	d["network_id"] = int(i.NetworkId)

	return d
}

func expandCreateVmInstanceGroup(d *schema.ResourceData) (ig sdk2.CreateVmInstanceGroup, interfaces []sdk2.VmInstanceGroupInterface) {
	ig.InstanceCount = float64(d.Get("vm_instance_group_instance_count").(int))
	ig.TypeId = float64(d.Get("vm_type_id").(int))
	ig.DiskSizeGB = d.Get("vm_instance_group_disk_size_gbytes").(float64)
	ig.VolumeTemplateId = float64(d.Get("vm_instance_group_template_id").(int))

	return
}

func expandUpdateVmInstanceGroup(d *schema.ResourceData) (ig sdk2.UpdateVmInstanceGroup, interfaces []sdk2.UpdateVmInstanceGroupInterface) {
	ig.Label = d.Get("vm_instance_group_label").(string)

	ig.VmInstanceGroupInterfaces = []sdk2.UpdateVmInstanceGroupInterface{}
	if d.Get("vm_instance_group_interfaces") != nil {
		interfacesSet := d.Get("vm_instance_group_interfaces").(*schema.Set)

		for _, interfaceEntry := range interfacesSet.List() {
			interfaceEntryMap := interfaceEntry.(map[string]interface{})
			vmInterfaceUpdate := sdk2.UpdateVmInstanceGroupInterface{
				Id:        float64(interfaceEntryMap["id"].(int)),
				NetworkId: float64(interfaceEntryMap["network_id"].(int)),
			}

			ig.VmInstanceGroupInterfaces = append(ig.VmInstanceGroupInterfaces, vmInterfaceUpdate)
		}
	}

	customVariables := make(map[string]interface{})
	if d.Get("vm_instance_custom_variables") != nil {
		for k, v := range d.Get("vm_instance_custom_variables").(map[string]interface{}) {
			customVariables[k] = v.(string)
		}
	}
	var customVariablesValue interface{} = customVariables
	ig.CustomVariables = &customVariablesValue

	return
}
