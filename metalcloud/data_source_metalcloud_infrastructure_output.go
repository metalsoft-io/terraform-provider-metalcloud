package metalcloud

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

// DataSourceInfrastructureOutput provides a way to export infrastructure information through terraform output blocks
func DataSourceInfrastructureOutput() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInfrastructureOutputRead,
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
			"drives": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: false,
				Default:  nil,
			},
			"instances": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: false,
				Default:  nil,
			},
			"shared_drives": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: false,
				Default:  nil,
			},
			"clusters": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: false,
				Default:  nil,
			},
		},
	}
}

func dataSourceInfrastructureOutputRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	infrastructure_id := d.Get("infrastructure_id").(int)

	d.Set("infrastructure_id", infrastructure_id)

	driveArrays, err := client.DriveArrays(infrastructure_id)

	if err != nil {
		return diag.FromErr(err)
	}

	var drivesMap = make(map[string]map[string]mc.Drive)

	for _, driveArray := range *driveArrays {
		drives, err := client.DriveArrayDrives(driveArray.DriveArrayID)

		if err != nil {
			return diag.FromErr(err)
		}

		drivesMap[driveArray.DriveArrayLabel] = *drives
	}

	drivesOutput, err := flattenDrives(&drivesMap)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("drives", drivesOutput)

	instances := []mc.Instance{}

	instanceArrays, err := client.InstanceArrays(infrastructure_id)

	if err != nil {
		return diag.FromErr(err)
	}

	for _, instanceArray := range *instanceArrays {
		retInstances, err := client.InstanceArrayInstances(instanceArray.InstanceArrayID)

		if err != nil {
			return diag.FromErr(err)
		}

		for _, instance := range *retInstances {
			i, err := client.InstanceGet(instance.InstanceID)

			if err != nil {
				return diag.FromErr(err)
			}
			instances = append(instances, *i)
		}
	}

	instancesOutput, err := flattenInstancesInfo(instances)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("instances", instancesOutput)

	sharedDrives, err := client.SharedDrives(infrastructure_id)

	if err != nil {
		return diag.FromErr(err)
	}

	sharedDrivesOutput, err := flattenSharedDrives(sharedDrives)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("shared_drives", sharedDrivesOutput)

	clusters, err := client.Clusters(infrastructure_id)
	clusterAppObjects := map[string]interface{}{}

	for label, cluster := range *clusters {

		switch cluster.ClusterType {
		case mc.CLUSTER_TYPE_VMWARE_VSPHERE:

			c, err := client.ClusterAppVMWareVSphere(cluster.ClusterID, true)

			if err != nil {
				return diag.FromErr(err)
			}
			clusterAppObjects[label] = c

		case mc.CLUSTER_TYPE_KUBERNETES:
			c, err := client.ClusterAppKubernetes(cluster.ClusterID, true)

			if err != nil {
				return diag.FromErr(err)
			}
			clusterAppObjects[label] = c
		}

	}

	clustersOutput, err := flattenClusters(&clusterAppObjects)

	d.Set("clusters", clustersOutput)

	d.SetId(fmt.Sprintf("%d", infrastructure_id))

	return diags
}

func flattenDrives(drivesMap *map[string]map[string]mc.Drive) (string, error) {
	drivesOutput := make(map[string]interface{})

	for label, drives := range *drivesMap {

		driveDetails := make(map[string]map[string]string)

		for k, v := range drives {
			driveDetails[k] = make(map[string]string)
			driveDetails[k]["drive_wwn"] = v.DriveWWN
		}

		drivesOutput[label] = driveDetails
	}

	bytes, err := json.Marshal(drivesOutput)

	if err != nil {
		return "", fmt.Errorf("error serializing drives: %s", err)
	}

	return string(bytes), nil
}

func flattenInstancesInfo(instances []mc.Instance) (string, error) {
	instancesOutput := make(map[string]interface{})

	for _, instance := range instances {
		label := instance.InstanceLabel

		instanceDetails := make(map[string]interface{})
		instanceDetails["instance_credentials"] = instance.InstanceCredentials
		instanceDetails["instance_array_id"] = instance.InstanceArrayID
		instancesOutput[label] = instanceDetails
	}

	bytes, err := json.Marshal(instancesOutput)

	if err != nil {
		return "", fmt.Errorf("error serializing instances array: %s", err)
	}

	return string(bytes), nil
}

func flattenSharedDrives(sharedDrives *map[string]mc.SharedDrive) (string, error) {
	sharedDrivesOutput := make(map[string]interface{})

	for label, sharedDrive := range *sharedDrives {
		sharedDriveDetails := make(map[string]interface{})
		sharedDriveDetails["shared_drive_targets_json"] = sharedDrive.SharedDriveTargetsJSON
		sharedDriveDetails["shared_drive_wwn"] = sharedDrive.SharedDriveWWN
		sharedDrivesOutput[label] = sharedDriveDetails
	}

	bytes, err := json.Marshal(sharedDrivesOutput)

	if err != nil {
		return "", fmt.Errorf("error serializing shared drives: %s", err)
	}

	return string(bytes), nil
}

func flattenClusters(clusters *map[string]interface{}) (string, error) {

	bytes, err := json.Marshal(clusters)

	if err != nil {
		return "", fmt.Errorf("error serializing clusters: %s", err)
	}

	return string(bytes), nil
}
