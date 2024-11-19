package metalcloud

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
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
			fieldInfrastructureId: {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validateRequired,
			},
			fieldVmInstanceGroupId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			fieldVmInstanceGroupLabel: {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          nil,
				Computed:         true,
				DiffSuppressFunc: caseInsensitiveDiff,
				ValidateDiagFunc: validateLabel,
			},
			fieldVmInstanceGroupInstanceCount: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			fieldVmTypeId: {
				Type:     schema.TypeInt,
				Required: true,
			},
			fieldVmInstanceGroupDiskSizeGbytes: {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			fieldVmInstanceGroupTemplateId: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			fieldVmInstanceCustomVariables: {
				Type:     schema.TypeList,
				Elem:     resourceVmInstanceCustomVariable(),
				Optional: true,
			},
			fieldVmInstanceGroupInterfaces: {
				Type:     schema.TypeList,
				Optional: true,
				Default:  nil,
				Computed: true,
				Elem:     resourceVmInstanceGroupInterface(),
			},
			fieldVmInstanceGroupNetworkProfiles: {
				Type:     schema.TypeList,
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
			fieldVmInstanceIndex: {
				Type:     schema.TypeInt,
				Required: true,
			},
			fieldCustomVariables: {
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
			fieldVmInterfaceIndex: {
				Type:     schema.TypeInt,
				Required: true,
			},
			fieldNetworkId: {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceVmInstanceGroupNetworkProfile() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			fieldNetworkId: {
				Type:     schema.TypeInt,
				Required: true,
			},
			fieldNetworkProfileId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceVmInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	// Verify infrastructure exists
	infrastructureId := getInfrastructureId(d)

	_, _, err = client.InfrastructureApi.GetInfrastructure(ctx, float64(infrastructureId))
	if err != nil {
		return diag.Errorf("Infrastructure with Id %+v not found.", infrastructureId)
	}

	// Create VM instance group
	vmGroupCreate, _ := expandCreateVmInstanceGroup(d)

	vmInstanceGroupCreated, _, err := client.VMInstanceGroupApi.CreateVMInstanceGroup(ctx, vmGroupCreate, float64(infrastructureId))
	if err != nil {
		return extractApiError(err)
	}

	vmInstanceGroupId := int(vmInstanceGroupCreated.Id)
	d.SetId(fmt.Sprintf("%d", vmInstanceGroupId))

	// Create VM instances custom variables
	diags := updateVmInstanceCustomVariables(ctx, client, d, infrastructureId, vmInstanceGroupId)
	if diags.HasError() {
		resourceVmInstanceGroupRead(ctx, d, meta)
		return diags
	}

	// Create VM instance group interfaces
	err = createMissingInterfaces(ctx, client, d, nil, infrastructureId, vmInstanceGroupId)
	if err != nil {
		resourceVmInstanceGroupRead(ctx, d, meta)
		return extractApiError(err)
	}

	// Set VM instance group network profiles
	err = updateNetworkProfiles(ctx, client, d, infrastructureId, vmInstanceGroupId)
	if err != nil {
		resourceVmInstanceGroupRead(ctx, d, meta)
		return extractApiError(err)
	}

	return resourceVmInstanceGroupRead(ctx, d, meta)
}

func resourceVmInstanceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstanceGroupId, infrastructureId, err := getVmInstanceGroupId(d)
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstanceGroup, _, err := client.VMInstanceGroupApi.GetVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	err = flattenVmInstanceGroup(ctx, d, vmInstanceGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceVmInstanceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstanceGroupId, infrastructureId, err := getVmInstanceGroupId(d)
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstanceGroup, _, err := client.VMInstanceGroupApi.GetVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	// Update VM instance group, including custom variables and interfaces
	vmInstanceGroupUpdates := expandUpdateVmInstanceGroup(d, &vmInstanceGroup)

	_, _, err = client.VMInstanceGroupApi.UpdateVMInstanceGroup(ctx, vmInstanceGroupUpdates, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	// Update VM instance custom variables
	diags := updateVmInstanceCustomVariables(ctx, client, d, infrastructureId, vmInstanceGroupId)
	if diags.HasError() {
		resourceVmInstanceGroupRead(ctx, d, meta)
		return diags
	}

	// Create missing interfaces
	err = createMissingInterfaces(ctx, client, d, &vmInstanceGroup, infrastructureId, vmInstanceGroupId)
	if err != nil {
		resourceVmInstanceGroupRead(ctx, d, meta)
		return extractApiError(err)
	}

	// Update VM instance group network profiles
	err = updateNetworkProfiles(ctx, client, d, infrastructureId, vmInstanceGroupId)
	if err != nil {
		resourceVmInstanceGroupRead(ctx, d, meta)
		return extractApiError(err)
	}

	return resourceVmInstanceGroupRead(ctx, d, meta)
}

func resourceVmInstanceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstanceGroupId, infrastructureId, err := getVmInstanceGroupId(d)
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

	return diag.Diagnostics{}
}

func flattenVmInstanceGroup(ctx context.Context, d *schema.ResourceData, vmInstanceGroup mc2.VmInstanceGroup) error {
	d.Set(fieldInfrastructureId, vmInstanceGroup.InfrastructureId)
	d.Set(fieldVmInstanceGroupId, vmInstanceGroup.Id)
	d.Set(fieldVmInstanceGroupLabel, vmInstanceGroup.Label)
	d.Set(fieldVmInstanceGroupInstanceCount, vmInstanceGroup.InstanceCount)
	d.Set(fieldVmInstanceGroupDiskSizeGbytes, vmInstanceGroup.DiskSizeGB)
	d.Set(fieldVmInstanceGroupTemplateId, vmInstanceGroup.VolumeTemplateId)

	/* NETWORK PROFILES */
	if vmInstanceGroup.NetworkIdToNetworkProfileId != nil {
		networkProfiles := flattenVmInstanceGroupNetworkProfiles(vmInstanceGroup)

		if len(networkProfiles) > 0 {
			d.Set(fieldVmInstanceGroupNetworkProfiles, schema.NewSet(schema.HashResource(resourceVmInstanceGroupNetworkProfile()), networkProfiles))
		}
	}

	/* INTERFACES */
	vmInterfacesData, ok := d.GetOk(fieldVmInstanceGroupInterfaces)
	if ok && vmInterfacesData != nil {
		vmInterfaces := []interface{}{}

		for _, vmInterface := range vmInterfacesData.([]interface{}) {
			vmInterfaceMap := vmInterface.(map[string]interface{})
			interfaceIndex := vmInterfaceMap[fieldVmInterfaceIndex].(int)

			// locate interface with index in returned data
			for _, vmInterface := range vmInstanceGroup.VmInstanceGroupInterfaces {
				// if we found it, locate the network it's connected to add it to the list
				if int(vmInterface.InterfaceIndex) == interfaceIndex && vmInterface.NetworkId != 0 {
					vmInterfaces = append(vmInterfaces, flattenVmInstanceGroupInterface(vmInterface))
				}
			}
		}

		if len(vmInterfaces) > 0 {
			d.Set(fieldVmInstanceGroupInterfaces, schema.NewSet(schema.HashResource(resourceVmInstanceGroupInterface()), vmInterfaces))
		}
	}

	/* CUSTOM VARIABLES */
	customVariablesMap := []interface{}{}
	for index, vmInstance := range vmInstanceGroup.VmInstances {
		vmInstanceEntry := make(map[string]interface{})
		vmInstanceEntryCustomVariables := make(map[string]interface{})

		vmInstanceEntry[fieldVmInstanceIndex] = index

		vmInstanceCustomVariables := *(vmInstance.CustomVariables)
		switch vmInstanceCustomVariables.(type) {
		case []interface{}:
			vmInstanceEntryCustomVariables = make(map[string]interface{})

		default:
			for k, v := range vmInstanceCustomVariables.(map[string]interface{}) {
				vmInstanceEntryCustomVariables[k] = v.(string)
			}
		}

		vmInstanceEntry[fieldCustomVariables] = vmInstanceEntryCustomVariables

		if len(vmInstanceEntryCustomVariables) > 0 {
			customVariablesMap = append(customVariablesMap, vmInstanceEntry)
		}
	}

	d.Set(fieldVmInstanceCustomVariables, customVariablesMap)

	return nil
}

func flattenVmInstanceGroupNetworkProfiles(vmInstanceGroup mc2.VmInstanceGroup) []interface{} {
	networkIdToNetworkProfileId := *vmInstanceGroup.NetworkIdToNetworkProfileId
	networkIdToNetworkProfileIdMap := networkIdToNetworkProfileId.(map[string]interface{})

	networkProfiles := []interface{}{}
	for networkId, networkProfileId := range networkIdToNetworkProfileIdMap {
		networkIdInt, err := strconv.Atoi(networkId)
		if err != nil {
			continue
		}

		networkProfiles = append(networkProfiles, map[string]interface{}{
			fieldNetworkId:        networkIdInt,
			fieldNetworkProfileId: int(networkProfileId.(float64)),
		})
	}

	return networkProfiles
}

func flattenVmInstanceGroupInterface(vmInstanceGroupInterface mc2.VmInstanceGroupInterface) map[string]interface{} {
	var interfaceMap = make(map[string]interface{})

	interfaceMap[fieldVmInterfaceIndex] = int(vmInstanceGroupInterface.InterfaceIndex)
	interfaceMap[fieldNetworkId] = int(vmInstanceGroupInterface.NetworkId)

	return interfaceMap
}

func expandCreateVmInstanceGroup(d *schema.ResourceData) (ig mc2.CreateVmInstanceGroup, interfaces []mc2.VmInstanceGroupInterface) {
	ig.InstanceCount = float64(d.Get(fieldVmInstanceGroupInstanceCount).(int))
	ig.TypeId = float64(d.Get(fieldVmTypeId).(int))
	ig.DiskSizeGB = d.Get(fieldVmInstanceGroupDiskSizeGbytes).(float64)
	ig.VolumeTemplateId = float64(d.Get(fieldVmInstanceGroupTemplateId).(int))

	return
}

func expandUpdateVmInstanceGroup(d *schema.ResourceData, vmInstanceGroup *mc2.VmInstanceGroup) (vmInstanceGroupUpdates mc2.UpdateVmInstanceGroup) {
	vmInstanceGroupUpdates.Label = d.Get(fieldVmInstanceGroupLabel).(string)

	vmInstanceGroupUpdates.VmInstanceGroupInterfaces = []mc2.UpdateVmInstanceGroupInterface{}

	interfacesData := d.Get(fieldVmInstanceGroupInterfaces)
	if interfacesData != nil {
		vmInterfaces := getInterfacesMap(vmInstanceGroup)

		for _, interfaceEntry := range interfacesData.([]interface{}) {
			interfaceEntryMap := interfaceEntry.(map[string]interface{})
			if vmInterface, ok := vmInterfaces[interfaceEntryMap[fieldVmInterfaceIndex].(int)]; ok {
				vmInterfaceUpdate := mc2.UpdateVmInstanceGroupInterface{
					Id:        float64(vmInterface.Id),
					NetworkId: float64(interfaceEntryMap[fieldNetworkId].(int)),
				}

				vmInstanceGroupUpdates.VmInstanceGroupInterfaces = append(vmInstanceGroupUpdates.VmInstanceGroupInterfaces, vmInterfaceUpdate)
			}
		}
	}

	return
}

func getInfrastructureId(d *schema.ResourceData) int {
	return d.Get(fieldInfrastructureId).(int)
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
func getSortedListOfVmInstances(ctx context.Context, client *mc2.APIClient, infrastructureId int, vmInstanceGroupId int) ([]mc2.VmInstance, error) {
	vmInstanceList, _, err := client.VMInstanceGroupApi.GetVMInstanceGroupVMInstances(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return nil, err
	}

	instanceMap := make(map[int]mc2.VmInstance, len(vmInstanceList))

	keys := []int{}
	instances := []mc2.VmInstance{}

	for _, v := range vmInstanceList {
		instanceMap[int(v.Id)] = v
		keys = append(keys, int(v.Id))
	}

	sort.Ints(keys)

	for _, id := range keys {
		instances = append(instances, instanceMap[id])
	}

	return instances, nil
}

func getInterfacesMap(vmInstanceGroup *mc2.VmInstanceGroup) map[int]mc2.VmInstanceGroupInterface {
	vmInterfaces := make(map[int]mc2.VmInstanceGroupInterface)

	if vmInstanceGroup != nil && vmInstanceGroup.VmInstanceGroupInterfaces != nil {
		for _, vmInterface := range vmInstanceGroup.VmInstanceGroupInterfaces {
			vmInterfaces[int(vmInterface.InterfaceIndex)] = vmInterface
		}
	}

	return vmInterfaces
}

func updateVmInstanceCustomVariables(ctx context.Context, client *mc2.APIClient, d *schema.ResourceData, infrastructureId int, vmInstanceGroupId int) diag.Diagnostics {
	vmInstancesList, err := getSortedListOfVmInstances(ctx, client, infrastructureId, vmInstanceGroupId)
	if err != nil {
		return extractApiError(err)
	}

	customVariablesData := d.Get(fieldVmInstanceCustomVariables)
	if customVariablesData != nil {
		for _, vmInstanceCustomVariablesData := range customVariablesData.([]interface{}) {
			vmInstanceCustomVariables := vmInstanceCustomVariablesData.(map[string]interface{})
			vmInstanceIndex := vmInstanceCustomVariables[fieldVmInstanceIndex].(int)

			if vmInstanceIndex < 0 || vmInstanceIndex >= len(vmInstancesList) {
				return diag.Diagnostics{diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("%s in %s block out of bounds", fieldVmInstanceIndex, fieldVmInstanceCustomVariables),
					Detail:   fmt.Sprintf("Use a number between 0 and %v (%s-1). ", len(vmInstancesList)-1, fieldVmInstanceGroupInstanceCount),
				}}
			}

			var customVariablesValue interface{} = vmInstanceCustomVariables[fieldCustomVariables].(map[string]interface{})
			vmInstanceUpdates := mc2.UpdateVmInstance{
				CustomVariables: &customVariablesValue,
			}

			_, _, err := client.VMInstanceApi.UpdateVMInstance(ctx, vmInstanceUpdates, float64(infrastructureId), vmInstancesList[vmInstanceIndex].Id)
			if err != nil {
				return extractApiError(err)
			}

			vmInstancesList[vmInstanceIndex].Id = 0
		}
	}

	// Remove custom variables from instances that were not updated
	for _, vmInstance := range vmInstancesList {
		if vmInstance.Id != 0 {
			vmInstanceUpdates := mc2.UpdateVmInstance{
				CustomVariables: nil,
			}

			_, _, err := client.VMInstanceApi.UpdateVMInstance(ctx, vmInstanceUpdates, float64(infrastructureId), vmInstance.Id)
			if err != nil {
				return extractApiError(err)
			}
		}
	}

	return diag.Diagnostics{}
}

func createMissingInterfaces(ctx context.Context, client *mc2.APIClient, d *schema.ResourceData, vmInstanceGroup *mc2.VmInstanceGroup, infrastructureId int, vmInstanceGroupId int) error {
	interfacesData := d.Get(fieldVmInstanceGroupInterfaces)
	if interfacesData != nil {
		vmInterfaces := getInterfacesMap(vmInstanceGroup)

		for _, interfaceEntry := range interfacesData.([]interface{}) {
			interfaceEntryMap := interfaceEntry.(map[string]interface{})
			if _, ok := vmInterfaces[interfaceEntryMap[fieldVmInterfaceIndex].(int)]; !ok {
				vmInterfaceCreate := mc2.CreateVmInstanceGroupInterface{
					NetworkId: float64(interfaceEntryMap[fieldNetworkId].(int)),
				}

				_, _, err := client.VMInstanceGroupApi.CreateVMInterfaceOnVMInstanceGroup(
					ctx,
					vmInterfaceCreate,
					float64(infrastructureId),
					float64(vmInstanceGroupId),
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func updateNetworkProfiles(ctx context.Context, client *mc2.APIClient, d *schema.ResourceData, infrastructureId int, vmInstanceGroupId int) error {
	networkProfilesData := d.Get(fieldVmInstanceGroupNetworkProfiles)
	if networkProfilesData != nil {
		for _, networkProfileEntry := range networkProfilesData.([]interface{}) {
			networkProfileEntryMap := networkProfileEntry.(map[string]interface{})
			networkId := networkProfileEntryMap[fieldNetworkId].(int)
			networkProfileUpdates := mc2.UpdateVmInstanceGroupNetwork{
				NetworkProfileId: float64(networkProfileEntryMap[fieldNetworkProfileId].(int)),
			}

			_, _, err := client.VMInstanceGroupApi.UpdateNetworkProfileOnVMInstanceGroupNetwork(
				ctx,
				networkProfileUpdates,
				float64(infrastructureId),
				float64(vmInstanceGroupId),
				float64(networkId),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
