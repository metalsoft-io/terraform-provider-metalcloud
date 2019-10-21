package metalcloud

import "log"
import (
	"github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceInstanceArray() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_array_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_array_instance_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"drive_array": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceDriveArray(),
			},
		},
	}
}

func resourceInstanceArrayCreate(infrastructureID int64, d map[string]interface{}, meta interface{}) error {

	client := meta.(*metalcloud.MetalCloudClient)

	instanceArray := metalcloud.InstanceArray{
		InstanceArrayLabel:         d["instance_array_label"].(string),
		InstanceArrayInstanceCount: d["instance_array_instance_count"].(int64),
	}

	createdInstanceArray, err := client.InstanceArrayCreate(infrastructureID, instanceArray)
	if err != nil || createdInstanceArray == nil {
		return err
	}

	driveArrays := d["drive_array"].(*schema.Set)

	log.Printf("Created InstanceArray %d", createdInstanceArray.InstanceArrayID)

	for _, driveArray := range driveArrays.List() {
		err := resourceDriveArrayCreate(infrastructureID,
			createdInstanceArray.InstanceArrayID,
			driveArray.(map[string]interface{}),
			meta)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceInstanceArrayRead(instanceArray metalcloud.InstanceArray, meta interface{}) (*map[string]interface{}, error) {

	client := meta.(*metalcloud.MetalCloudClient)

	var instanceArrayMap = make(map[string]interface{})

	instanceArrayMap["instance_array_label"] = instanceArray.InstanceArrayLabel
	instanceArrayMap["instance_array_instance_count"] = int64(instanceArray.InstanceArrayInstanceCount)

	var driveArraysOfThisInstanceArray []interface{}
	driveArrays, err := client.DriveArrays(instanceArray.InfrastructureID)
	if err != nil {
		return nil, err
	}
	for _, driveArray := range *driveArrays {
		if driveArray.InstanceArrayID == instanceArray.InstanceArrayID {
			driveArrayMap, err := resourceDriveArrayRead(driveArray, meta)
			if err != nil {
				return nil, err
			}
			driveArraysOfThisInstanceArray = append(driveArraysOfThisInstanceArray, *driveArrayMap)
		}
	}

	instanceArrayMap["drive_array"] = schema.NewSet(
		schema.HashResource(resourceDriveArray()),
		driveArraysOfThisInstanceArray)

	return &instanceArrayMap, nil
}
