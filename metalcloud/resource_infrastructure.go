package metalcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//ResourceInfrastructure returns the top infrastructure resource
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
				Type:     schema.TypeList,
				Required: true,
				Elem:     resourceInstanceArray(),
			},
			"prevent_deploy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceInstanceArray() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_array_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"instance_array_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_array_instance_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			/*
				"instance_array_subdomain": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
			*/
			"instance_array_boot_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pxe_iscsi",
			},
			"instance_array_ram_gbytes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"instance_array_processor_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"instance_array_processor_core_mhz": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},
			"instance_array_processor_core_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"instance_array_disk_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"instance_array_disk_size_mbytes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
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
			},
			"drive_array": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceDriveArray(),
			},
		},
	}
}

func resourceDriveArray() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"drive_array_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"drive_array_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"volume_template_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"drive_array_storage_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "auto" {
						return true
					}
					return false
				},
			},
			"drive_size_mbytes_default": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  40960,
			},
			"instance_array_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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

func resourceInfrastructureCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*metalcloud.MetalCloudClient)

	infrastructure := metalcloud.Infrastructure{
		InfrastructureLabel: d.Get("infrastructure_label").(string),
		DatacenterName:      d.Get("datacenter_name").(string),
	}

	createdInfra, err := client.InfrastructureCreate(infrastructure)
	if err != nil || createdInfra == nil {
		return err
	}

	if instanceArrays, ok := d.GetOk("instance_array"); ok {

		for i, resIA := range instanceArrays.([]interface{}) {

			ia, daList := expandInstanceArrayWithDriveArrays(resIA.(map[string]interface{}))

			iaCreated, err := client.InstanceArrayCreate(createdInfra.InfrastructureID, ia)
			if err != nil {
				return err
			}

			d.Set(fmt.Sprintf("instance_array.%d.instance_array_id", i), iaCreated.InstanceArrayID)

			for di, da := range daList {
				da.InstanceArrayID = iaCreated.InstanceArrayID
				daCreated, err := client.DriveArrayCreate(createdInfra.InfrastructureID, da)
				if err != nil {
					return err
				}
				d.Set(fmt.Sprintf("instance_array.%d.drive_array.%d.drive_array_id", i, di), daCreated.DriveArrayID)
			}
		}
	}

	d.SetId(fmt.Sprintf("%d", createdInfra.InfrastructureID))

	if preventDeploy, ok := d.GetOk("prevent_deploy"); !ok || preventDeploy == false {
		client.InfrastructureDeploy(createdInfra.InfrastructureID, metalcloud.ShutdownOptions{}, true, false)
	}

	return resourceInfrastructureRead(d, meta)
}

//resourceInfrastructureRead reads the serverside status of elements
//it merges the serverside status with what is stored in the current state
func resourceInfrastructureRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*metalcloud.MetalCloudClient)

	infrastructureID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	log.Printf("resourceInfrastructureRead for infrastructureID=%d", infrastructureID)

	infrastructure, err := client.InfrastructureGet(infrastructureID)
	if err != nil {
		return err
	}

	if infrastructure == nil {
		log.Printf("Could not get infrastructure with id=%d", infrastructureID)
		d.SetId("") //404
		return nil
	}

	if err := d.Set("infrastructure_label", infrastructure.InfrastructureLabel); err != nil {
		return fmt.Errorf("error setting infrastructure_label: %s", err)
	}

	if err := d.Set("datacenter_name", infrastructure.DatacenterName); err != nil {
		return fmt.Errorf("error setting datacenter_name: %s", err)
	}

	retInstanceArrays, err := client.InstanceArrays(infrastructureID)
	if err != nil {
		return err
	}

	retDriveArrays, err := client.DriveArrays(infrastructureID)
	if err != nil {
		return err
	}

	iaList := []interface{}{}

	if iaCount, ok := d.Get("instance_array.#").(int); ok {
		for i := 0; i < iaCount; i++ {
			if instanceArrayID, ok := d.GetOk(fmt.Sprintf("instance_array.%d.instance_array_id", i)); ok {
				//we get the instance array again because the label might have changed so we cannot use the label index of the retInstanceArrays
				ia, err := client.InstanceArrayGet(instanceArrayID.(int))
				if err != nil {
					return err
				}
				//get the drive arrays of the current instance array
				daList := []metalcloud.DriveArray{}
				if daCount, ok := d.Get("instance_array.%d.drive_array.#").(int); ok {
					for di := 0; di < daCount; di++ {
						if driveArrayID, ok := d.GetOk(fmt.Sprintf("instance_array.%d.drive_array.%d.drive_array_id", i, di)); ok {
							da, err := client.DriveArrayGet(driveArrayID.(int))
							if err != nil {
								return err
							}
							daList = append(daList, *da)
							//we delete this from the serverside list so at the end we have only new elements
							//that we will be appending to the state
							delete(*retDriveArrays, fmt.Sprintf("%s.vanilla", da.DriveArrayLabel))
						}
					}
				}

				//iterate over the remaining drive arrays to see if we got any new drive
				//arrays for this instance array
				for _, da := range *retDriveArrays {
					if da.InstanceArrayID == ia.InstanceArrayID {
						daList = append(daList, da)
					}
				}

				iaList = append(iaList, flattenInstanceArrayWithDriveArrays(*ia, daList))

				//delete record from serverside list so at the end we have only new elements
				//that we will be appending to the state
				delete(*retInstanceArrays, fmt.Sprintf("%s.vanilla", ia.InstanceArrayLabel))

			}
		}
	}

	//append remaining elements (new on the serverside)
	for _, ia := range *retInstanceArrays {
		var daList []metalcloud.DriveArray
		for _, da := range *retDriveArrays {
			if da.InstanceArrayID == ia.InstanceArrayID {
				daList = append(daList, da)
			}
		}
		iaList = append(iaList, flattenInstanceArrayWithDriveArrays(ia, daList))
	}

	j, _ := json.MarshalIndent(iaList, "", "\t")
	log.Printf("flattened list of instance arrays is now %s", j)

	if err := d.Set("instance_array", iaList); err != nil {
		return fmt.Errorf("error setting instance_array: %s", err)
	}

	return nil
}

func resourceInfrastructureUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	client := meta.(*metalcloud.MetalCloudClient)

	infrastructureID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	needsDeploy := false

	if d.HasChange("infrastructure_label") || d.HasChange("datacenter_name") {

		infrastructure, err := client.InfrastructureGet(int(infrastructureID))
		if err != nil {
			return err
		}

		operation := infrastructure.InfrastructureOperation
		operation.InfrastructureLabel = d.Get("infrastructure_label").(string)
		operation.DatacenterName = d.Get("datacenter_name").(string)

		if _, err = client.InfrastructureEdit(int(infrastructureID), operation); err != nil {
			return err
		}

		needsDeploy = true
	}

	if d.HasChange("instance_array") {
		//take each instance array and apply changes
		currentInstanceArraysMap := d.Get("instance_array").([]interface{})

		for i, iaMap := range currentInstanceArraysMap {

			if !d.HasChange(fmt.Sprintf("instance_array.%d", i)) {
				continue
			}

			ia, daList := expandInstanceArrayWithDriveArrays(iaMap.(map[string]interface{}))

			retIA, err := createOrUpdateInstanceArray(infrastructureID, ia, client)
			if err != nil {
				return err
			}

			for di, da := range daList {
				if !d.HasChange(fmt.Sprintf("instance_array.%d.drive_array.%d", i, di)) {
					continue
				}
				da.InstanceArrayID = retIA.InstanceArrayID
				if _, err := createOrUpdateDriveArray(infrastructureID, da, client); err != nil {
					return err
				}
				needsDeploy = true
			}

			needsDeploy = true
		}
	}

	if needsDeploy {
		if preventDeploy, ok := d.GetOk("prevent_deploy"); !ok || preventDeploy == false {
			client.InfrastructureDeploy(infrastructureID, metalcloud.ShutdownOptions{}, true, false)
		}
	}

	d.Partial(false)

	return resourceInfrastructureRead(d, meta)
}

func createOrUpdateInstanceArray(infrastructureID int, ia metalcloud.InstanceArray, client *metalcloud.MetalCloudClient) (*metalcloud.InstanceArray, error) {
	var instanceArrayID = ia.InstanceArrayID

	var iaToReturn *metalcloud.InstanceArray

	if instanceArrayID == 0 {

		retIA, err := client.InstanceArrayCreate(infrastructureID, ia)
		if err != nil {
			return nil, err
		}
		iaToReturn = retIA

	} else {

		retIA, err := client.InstanceArrayGet(ia.InstanceArrayID)
		if err != nil {
			return nil, err
		}

		copyInstanceArrayToOperation(ia, retIA.InstanceArrayOperation)

		retIA2, err2 := client.InstanceArrayEdit(retIA.InstanceArrayID, *retIA.InstanceArrayOperation)
		if err2 != nil {
			return nil, err2
		}
		iaToReturn = retIA2
	}
	return iaToReturn, nil
}

func createOrUpdateDriveArray(infrastructureID int, da metalcloud.DriveArray, client *metalcloud.MetalCloudClient) (*metalcloud.DriveArray, error) {
	var driveArrayToReturn *metalcloud.DriveArray
	if da.DriveArrayID == 0 {
		retDA, err := client.DriveArrayCreate(infrastructureID, da)
		if err != nil {
			return nil, err
		}
		driveArrayToReturn = retDA
	} else {
		retDA, err := client.DriveArrayGet(da.DriveArrayID)
		if err != nil {
			return nil, err
		}

		copyDriveArrayToOperation(da, retDA.DriveArrayOperation)
		retDA, err2 := client.DriveArrayEdit(da.DriveArrayID, *retDA.DriveArrayOperation)
		if err2 != nil {
			return nil, err2
		}
		driveArrayToReturn = retDA
	}
	return driveArrayToReturn, nil
}

func resourceInfrastructureDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*metalcloud.MetalCloudClient)

	infrastructureID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	client.InfrastructureDelete(infrastructureID)

	d.SetId("")
	return nil
}
