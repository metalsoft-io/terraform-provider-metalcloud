package metalcloud

import "time"
import "strconv"
import "fmt"
import "log"
import (
	"github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ResourceInfrastructure() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfrastructureCreate,
		Read:   resourceInfrastructureRead,
		Update: resourceInfrastructureUpdate,
		Delete: resourceInfrastructureDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"infrastructure_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_array": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceInstanceArray(),
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceInfrastructureCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*metalcloud.MetalCloudClient)

	d.Partial(true)

	infrastructure := metalcloud.Infrastructure{
		InfrastructureLabel: d.Get("infrastructure_label").(string),
		DatacenterName:      d.Get("datacenter_name").(string),
	}

	created_infrastructure, err := client.InfrastructureCreate(infrastructure)
	if err != nil || created_infrastructure == nil {
		return err
	}

	infrastructureID := created_infrastructure.InfrastructureID

	d.SetId(fmt.Sprintf("%d", int(infrastructureID)))

	if d.HasChange("instance_array") {
		instanceArraysSet := d.Get("instance_array").(*schema.Set)

		log.Printf("instanceArraysSet=%s", instanceArraysSet.GoString())

		for _, instanceArray := range instanceArraysSet.List() {

			err := resourceInstanceArrayCreate(infrastructureID,
				instanceArray.(map[string]interface{}),
				meta)
			if err != nil {
				return err
			}
		}
	}
	d.Partial(false)

	return resourceInfrastructureRead(d, meta)
}

func resourceInfrastructureRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*metalcloud.MetalCloudClient)

	infrastructureID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	infrastructure, err := client.InfrastructureGet(int64(infrastructureID))
	if err != nil {
		return err
	}

	if infrastructure == nil {
		log.Printf("Could not get infrastructure with id=%d", infrastructureID)
		d.SetId("") //404
		return nil
	}

	d.Set("infrastructure_label", infrastructure.InfrastructureLabel)
	d.Set("datacenter_name", infrastructure.DatacenterName)

	var instanceArraysList []interface{}
	instanceArrays, err := client.InstanceArrays(int64(infrastructureID))
	if err != nil {
		return err
	}

	for _, instanceArray := range *instanceArrays {
		instanceArrayMap, err := resourceInstanceArrayRead(instanceArray, meta)
		if err != nil {
			return err
		}

		instanceArraysList = append(instanceArraysList, *instanceArrayMap)
	}

	d.Set("instance_array", schema.NewSet(
		schema.HashResource(resourceInstanceArray()),
		instanceArraysList))

	return nil
}

func resourceInfrastructureUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	client := meta.(*metalcloud.MetalCloudClient)

	infrastructureID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("infrastructure_label") || d.HasChange("datacenter_name") {

		infrastructure, err := client.InfrastructureGet(int64(infrastructureID))
		if err != nil {
			return err
		}

		operation := infrastructure.InfrastructureOperation
		operation.InfrastructureLabel = d.Get("infrastructure_label").(string)
		operation.DatacenterName = d.Get("datacenter_name").(string)

		_, err = client.InfrastructureEdit(int64(infrastructureID), operation)
		if err != nil {
			return err
		}
	}

	//if d.HasChange("instance_array")

	d.Partial(false)

	return resourceInfrastructureRead(d, meta)
}

func resourceInfrastructureDelete(d *schema.ResourceData, meta interface{}) error {

	//client := meta.(*metalcloud.MetalCloudClient)

	d.SetId("")
	return nil
}
