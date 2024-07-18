package metalcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
)

func resourceClusterAppSchema(groupRolesSuffixes map[string]string) map[string]*schema.Schema {

	schemaForOneInstanceArray := map[string]*schema.Schema{
		"instance_array_instance_count": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
			// ValidateDiagFunc: validateMaxOne,
		},

		"instance_server_type": {
			Type:     schema.TypeList,
			Elem:     resourceInstanceServerType(),
			Optional: true,
		},

		"instance_array_network_profile": {
			Type:     schema.TypeSet,
			Optional: true,
			Default:  nil,
			Computed: true,
			Elem:     resourceInstanceArrayNetworkProfile(),
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
		"interface": {
			Type:     schema.TypeSet,
			Optional: true,
			Default:  nil,
			Computed: true,
			Elem:     resourceInstanceArrayInterface(),
		},
	}

	schema := map[string]*schema.Schema{
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
		"instance_array": {
			Type:     schema.TypeList,
			Elem:     resourceInstanceArray(),
			Optional: true,
		},
	}

	for _, suffix := range groupRolesSuffixes {
		for key, value := range schemaForOneInstanceArray {
			schema[key+suffix] = value
		}
	}

	return schema
}

func resourceClusterAppCreate(clusterAppType string, groupRolesSuffixes map[string]string, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	dg := updateClusterInstanceArrays(groupRolesSuffixes, ctx, d, meta, retCl.ClusterID)
	diags = append(diags, dg...)

	dg = resourceClusterAppRead(groupRolesSuffixes, ctx, d, meta)
	diags = append(diags, dg...)

	return diags
}

func updateClusterInstanceArrays(groupRolesSuffixes map[string]string, ctx context.Context, d *schema.ResourceData, meta interface{}, clusterID int) diag.Diagnostics {
	client := meta.(*mc.Client)

	var diags diag.Diagnostics

	retIa, err := client.ClusterInstanceArrays(clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, ia := range *retIa {

		suffix := groupRolesSuffixes[ia.ClusterRoleGroup]

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

			//create extra interfaces if needed
			interfacesToCreate := len(interfaces) - len(ia.InstanceArrayInterfaces)
			if interfacesToCreate > 0 {
				log.Printf("Creating %d interfaces for instance array %d", interfacesToCreate, ia.InstanceArrayID)
			}
			for i := 0; i < interfacesToCreate; i++ {
				client.InstanceArrayInterfaceCreate(ia.InstanceArrayID)
			}

			for _, intf := range interfaces {
				_, err := client.InstanceArrayInterfaceAttachNetwork(ia.InstanceArrayID, intf.InstanceArrayInterfaceIndex, intf.NetworkID)
				if err != nil {
					return diag.FromErr(err)
				}
			}

			ia.InstanceArrayInterfaces = interfaces
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

func resourceClusterAppRead(groupRolesSuffixes map[string]string, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if err != nil {
		return diag.FromErr(err)
	}

	flattenAppCluster(groupRolesSuffixes, d, *cluster, *ia, client)

	return diags
}

func flattenAppCluster(groupRolesSuffixes map[string]string, d *schema.ResourceData, cluster mc.Cluster, ia map[string]mc.InstanceArray, client *mc.Client) error {

	d.Set("cluster_id", cluster.ClusterID)
	d.Set("infrastructure_id", cluster.InfrastructureID)
	d.Set("cluster_label", cluster.ClusterLabel)

	for _, ia := range ia {

		suffix := groupRolesSuffixes[ia.ClusterRoleGroup]

		log.Printf("importing role %s with suffix %s", ia.ClusterRoleGroup, suffix)

		d.Set("instance_array_instance_count"+suffix, ia.InstanceArrayInstanceCount)

		var intfList []interface{}
		for _, intf := range ia.InstanceArrayInterfaces {

			if intf.NetworkID != 0 { //we ignore unconnected interfaces
				intfList = append(intfList, flattenInstanceArrayInterface(intf))
			}
		}

		if len(intfList) > 0 {
			d.Set("interface"+suffix, schema.NewSet(schema.HashResource(resourceInstanceArrayInterface()), intfList))
		}

		networkToNetworkProfileMap, err := client.NetworkProfileListByInstanceArray(ia.InstanceArrayID)
		if err != nil {
			return err
		}

		var networkProfileList []interface{}
		for networkID, networkProfileID := range *networkToNetworkProfileMap {
			if networkProfileID > 0 {
				networkProfileEntry := flattenInstanceArrayNetworkProfile(networkID, networkProfileID)
				networkProfileList = append(networkProfileList, networkProfileEntry)
			}
		}

		if len(networkProfileList) > 0 {
			d.Set("instance_array_network_profile"+suffix, schema.NewSet(schema.HashResource(resourceInstanceArrayNetworkProfile()), networkProfileList))
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

		instancesCustomVariables := flattenInstancesCustomVariables(retInstances)

		if len(instancesCustomVariables) > 0 {
			d.Set("instance_custom_variables"+suffix, instancesCustomVariables)
		}

		instanceServerTypes := flattenInstanceServerTypes(retInstances)

		if len(instanceServerTypes) > 0 {
			d.Set("instance_server_type"+suffix, instanceServerTypes)
		}
	}
	return nil
}

func resourceClusterAppUpdate(groupRolesSuffixes map[string]string, ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	dg := updateClusterInstanceArrays(groupRolesSuffixes, ctx, d, meta, retCl.ClusterID)
	diags = append(diags, dg...)
	//}

	dg = resourceClusterAppRead(groupRolesSuffixes, ctx, d, meta)
	diags = append(diags, dg...)

	return diags

}

func expandClusterApp(d *schema.ResourceData) mc.Cluster {
	var c mc.Cluster

	c.ClusterID = d.Get("cluster_id").(int)
	c.ClusterLabel = d.Get("cluster_label").(string)
	c.ClusterSoftwareVersion = d.Get("cluster_software_version").(string)

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
