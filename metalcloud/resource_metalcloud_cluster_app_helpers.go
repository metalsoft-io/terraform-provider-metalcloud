package metalcloud

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func resourceClusterAppTwoIASchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"cluster_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"cluster_label": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  nil,
			Computed: true,
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
			ValidateDiagFunc: validateLabel,
		},

		"instance_array_instance_count_master": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
			// ValidateDiagFunc: validateMaxOne,
		},

		"instance_server_type_master": {
			Type:     schema.TypeList,
			Elem:     resourceInstanceServerType(),
			Optional: true,
		},

		"instance_array_network_profile_master": {
			Type:     schema.TypeSet,
			Optional: true,
			Default:  nil,
			Computed: true,
			Elem:     resourceInstanceArrayNetworkProfile(),
		},
		"instance_array_custom_variables_master": {
			Type:     schema.TypeMap,
			Elem:     schema.TypeString,
			Optional: true,
			Computed: true, //default is computed serverside
		},
		"instance_custom_variables_master": {
			Type:     schema.TypeList,
			Elem:     resourceInstanceCustomVariable(),
			Optional: true,
		},

		"instance_array_instance_count_worker": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
			// ValidateDiagFunc: validateMaxOne,
		},

		"instance_server_type_worker": {
			Type:     schema.TypeList,
			Elem:     resourceInstanceServerType(),
			Optional: true,
		},

		"instance_array_network_profile_worker": {
			Type:     schema.TypeSet,
			Optional: true,
			Default:  nil,
			Computed: true,
			Elem:     resourceInstanceArrayNetworkProfile(),
		},
		"instance_array_custom_variables_worker": {
			Type:     schema.TypeMap,
			Elem:     schema.TypeString,
			Optional: true,
			Computed: true, //default is computed serverside
		},
		"instance_custom_variables_worker": {
			Type:     schema.TypeList,
			Elem:     resourceInstanceCustomVariable(),
			Optional: true,
		},
		"interface_master": {
			Type:     schema.TypeSet,
			Optional: true,
			Default:  nil,
			Computed: true,
			Elem:     resourceInstanceArrayInterface(),
		},
		"interface_worker": {
			Type:     schema.TypeSet,
			Optional: true,
			Default:  nil,
			Computed: true,
			Elem:     resourceInstanceArrayInterface(),
		},
	}
}

func resourceClusterAppCreate(clusterAppType string, masterRoleGroupName string, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)
	var diags diag.Diagnostics

	infrastructure_id := d.Get("infrastructure_id").(int)
	_, err := client.InfrastructureGet(infrastructure_id)

	if err != nil {
		return diag.Errorf("Infrastructure with id %+v not found.", infrastructure_id)
	}

	var cluster = expandClusterApp(d)

	cluster.ClusterType = clusterAppType

	retCl, err := client.ClusterCreate(infrastructure_id, cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", retCl.ClusterID))

	dg := updateClusterInstanceArrays(masterRoleGroupName, ctx, d, meta, retCl.ClusterID)
	diags = append(diags, dg...)

	dg = resourceClusterAppRead(masterRoleGroupName, ctx, d, meta)
	diags = append(diags, dg...)

	return diags
}

func updateClusterInstanceArrays(masterRoleGroupName string, ctx context.Context, d *schema.ResourceData, meta interface{}, clusterID int) diag.Diagnostics {
	client := meta.(*mc.Client)

	var diags diag.Diagnostics

	retIa, err := client.ClusterInstanceArrays(clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, ia := range *retIa {

		var suffix string
		if ia.ClusterRoleGroup == masterRoleGroupName {
			suffix = "_master"
		} else {
			suffix = "_worker"
		}

		instanceArrayInstanceCount := d.Get("instance_array_instance_count" + suffix).(int)
		instanceArrayServerTypeList := d.Get("instance_server_type" + suffix).([]interface{})
		instanceArrayNetworkProfileList := d.Get("instance_array_network_profile" + suffix).(*schema.Set)
		instanceArrayInstanceVariables := d.Get("instance_custom_variables" + suffix).([]interface{})

		if d.Get("instance_array_custom_variables"+suffix) != nil {
			iacv := make(map[string]string)

			for k, v := range d.Get("instance_array_custom_variables" + suffix).(map[string]interface{}) {
				iacv[k] = v.(string)
			}

			ia.InstanceArrayCustomVariables = iacv
		}

		ia.InstanceArrayInstanceCount = instanceArrayInstanceCount
		copyInstanceArrayToOperation(ia, ia.InstanceArrayOperation)

		detachDrives := true
		swapHardware := false

		_, err := client.InstanceArrayEdit(ia.InstanceArrayID, *ia.InstanceArrayOperation, &swapHardware, &detachDrives, nil, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		dg := updateInstancesServerTypes(instanceArrayServerTypeList, ia.InstanceArrayID, client)
		if dg.HasError() {
			return dg
		}

		diags = append(diags, dg...)

		//attach interfaces to network

		if d.Get("interface"+suffix) != nil {
			interfaceSet := d.Get("interface" + suffix).(*schema.Set)
			interfaces := []mc.InstanceArrayInterface{}

			for _, intfList := range interfaceSet.List() {
				intfMap := intfList.(map[string]interface{})
				intfMap["instance_array_id"] = ia.InstanceArrayID
				interfaces = append(interfaces, expandInstanceArrayInterface(intfMap))
			}

			ia.InstanceArrayInterfaces = interfaces

			for _, intf := range ia.InstanceArrayInterfaces {
				_, err := client.InstanceArrayInterfaceAttachNetwork(ia.InstanceArrayID, intf.InstanceArrayInterfaceIndex, intf.NetworkID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		//network profiles
		if instanceArrayNetworkProfileList != nil {
			for _, profileIntf := range instanceArrayNetworkProfileList.List() {
				profileMap := profileIntf.(map[string]interface{})

				_, err := client.InstanceArrayNetworkProfileSet(ia.InstanceArrayID, profileMap["network_id"].(int), profileMap["network_profile_id"].(int))
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		/* custom variables for instances */
		dg = updateInstancesCustomVariables(instanceArrayInstanceVariables, ia.InstanceArrayID, client)
		diags = append(diags, dg...)
	}
	return diags
}

func resourceClusterAppRead(masterRoleGroupName string, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cluster, err := client.ClusterGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ia, err := client.ClusterInstanceArrays(id)

	flattenAppCluster(masterRoleGroupName, d, *cluster, *ia, client)

	return diags
}

func flattenAppCluster(masterRoleGroupName string, d *schema.ResourceData, cluster mc.Cluster, ia map[string]mc.InstanceArray, client *mc.Client) error {

	d.Set("cluster_id", cluster.ClusterID)
	d.Set("cluster_label", cluster.ClusterLabel)

	for _, ia := range ia {

		var suffix string
		if ia.ClusterRoleGroup == masterRoleGroupName {
			suffix = "_master"
		} else {
			suffix = "_worker"
		}

		d.Set("instance_array_instance_count"+suffix, ia.InstanceArrayInstanceCount)

		//iterate over interfaces
		interfaces := []interface{}{}
		if intfList, ok := d.GetOkExists("interface" + suffix); ok {
			for _, iIntf := range intfList.(*schema.Set).List() {
				iaInterface := iIntf.(map[string]interface{})
				interfaceIndex := iaInterface["interface_index"].(int)

				//locate interface with index in returned data
				for _, intf := range ia.InstanceArrayInterfaces {
					//if we found it, locate the network it's connected to add it to the list
					if intf.InstanceArrayInterfaceIndex == interfaceIndex && intf.NetworkID != 0 {
						interfaces = append(interfaces, flattenInstanceArrayInterface(intf))
					}
				}
			}
		}
		if len(interfaces) > 0 {
			d.Set("interface"+suffix, schema.NewSet(schema.HashResource(resourceInstanceArrayInterface()), interfaces))
		}

		networkProfiles, err := client.NetworkProfileListByInstanceArray(ia.InstanceArrayID)
		if err != nil {
			return err
		}

		profiles := flattenInstanceArrayNetworkProfile(*networkProfiles, d)

		if len(profiles) > 0 {
			d.Set("instance_array_network_profile"+suffix, schema.NewSet(schema.HashResource(resourceInstanceArrayNetworkProfile()), profiles))
		}

		/* INSTANCE ARRAY CUSTOM VARIABLES */
		switch ia.InstanceArrayCustomVariables.(type) {
		case []interface{}:
			d.Set("instance_array_custom_variables"+suffix, make(map[string]string))
		default:
			iacv := make(map[string]string)

			for k, v := range ia.InstanceArrayCustomVariables.(map[string]interface{}) {
				iacv[k] = v.(string)
			}
			d.Set("instance_array_custom_variables"+suffix, iacv)
		}

		/* INSTANCES CUSTOM VARS */
		retInstances, err := client.InstanceArrayInstances(ia.InstanceArrayID)
		if err != nil {
			return err
		}

		keys := []int{}
		instances := []mc.Instance{}

		for _, v := range *retInstances {

			keys = append(keys, v.InstanceID)
		}

		sort.Ints(keys)

		for _, id := range keys {
			i, err := client.InstanceGet(id)
			if err != nil {
				return err
			}
			instances = append(instances, *i)
		}

		instancesCustomVariables := flattenInstancesCustomVariables(retInstances)

		if len(instancesCustomVariables) > 0 {
			d.Set("instance_custom_variables"+suffix, instancesCustomVariables)

		}
	}

	return nil
}

func resourceClusterAppUpdate(masterRoleGroupName string, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)
	var diags diag.Diagnostics
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	retCl, err := client.ClusterGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	cluster := expandClusterApp(d)

	copyClusterToOperation(cluster, &cluster.ClusterOperation)

	/*if d.HasChange("instance_array_instance_count_master") ||
	d.HasChange("instance_array_instance_count_woker") ||
	d.HasChange("instance_server_type_master") ||
	d.HasChange("instance_server_type_worker") {
	*/
	dg := updateClusterInstanceArrays(masterRoleGroupName, ctx, d, meta, retCl.ClusterID)
	diags = append(diags, dg...)
	//}

	dg = resourceClusterAppRead(masterRoleGroupName, ctx, d, meta)
	diags = append(diags, dg...)

	return diags

}

func expandClusterApp(d *schema.ResourceData) mc.Cluster {
	var c mc.Cluster

	c.ClusterID = d.Get("cluster_id").(int)
	c.ClusterLabel = d.Get("cluster_label").(string)

	return c
}

func resourceClusterAppDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)
	var diags diag.Diagnostics
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.ClusterDelete(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
