package metalcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	mc "github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_DATACENTER", nil),
				Optional:    true,
			},
			"instance_array": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     resourceInstanceArray(),
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceNetwork(),
			},
			"prevent_deploy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"hard_shutdown_after_timeout": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"attempt_soft_shutdown": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"soft_shutdown_timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  180,
			},
			"allow_data_loss": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"skip_ansible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"await_deploy_finished": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"await_delete_finished": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"keep_detaching_drives": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
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

			"interface": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceInstanceArrayInterface(),
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
				Required: true,
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

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"network_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_lan_autoallocate_ips": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceInstanceArrayInterface() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"interface_index": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"network_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceInfrastructureCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*mc.Client)

	infrastructure := mc.Infrastructure{
		InfrastructureLabel: d.Get("infrastructure_label").(string),
		DatacenterName:      d.Get("datacenter_name").(string),
	}

	createdInfra, err := client.InfrastructureCreate(infrastructure)
	if err != nil || createdInfra == nil {
		return err
	}

	//create networks
	//there are some networks already created by infrastructure create so we will be updating their labels only
	retNetworks, err := client.Networks(createdInfra.InfrastructureID)
	if err != nil {
		return err
	}

	nCreatedNetworksMap := map[string]mc.Network{}
	if networks, ok := d.GetOk("network"); ok {
		for _, resN := range networks.([]interface{}) {

			n := expandNetwork(resN.(map[string]interface{}))

			//if existing network and label is different update it
			var network *mc.Network
			for _, existingN := range *retNetworks {
				if existingN.NetworkType == n.NetworkType {
					if existingN.NetworkLabel != n.NetworkLabel {
						existingN.NetworkOperation.NetworkLabel = n.NetworkLabel

						log.Printf("found existing network with id %d with type %s and label %s. Editing...", existingN.NetworkID, existingN.NetworkType, existingN.NetworkLabel)

						network, err = client.NetworkEdit(existingN.NetworkID, *existingN.NetworkOperation)
						if err != nil {
							return err
						}
					} else {
						network = &existingN
					}
					break
				}
			}

			//if no network with that type create it
			if network == nil {
				log.Printf("Creating network %s (%s)", n.NetworkLabel, n.NetworkType)
				network, err = client.NetworkCreate(createdInfra.InfrastructureID, n)
				if err != nil {
					return err
				}
			}
			nCreatedNetworksMap[network.NetworkLabel] = *network
		}
	}

	//create instance arrays (and their drives)
	if instanceArrays, ok := d.GetOkExists("instance_array"); ok {

		for i, resIA := range instanceArrays.([]interface{}) {

			//populate instance array
			ia := expandInstanceArray(resIA.(map[string]interface{}))

			//create the instance array
			iaCreated, err := client.InstanceArrayCreate(createdInfra.InfrastructureID, ia)
			if err != nil {
				return err
			}

			//add interfaces to networks

			intListIntf := d.Get(fmt.Sprintf("instance_array.%d.interface", i))

			for _, intIntf := range intListIntf.([]interface{}) {
				intIntfMap := intIntf.(map[string]interface{})
				networkLabel := intIntfMap["network_label"].(string)

				//look for a network with the given label
				if n, ok := nCreatedNetworksMap[networkLabel]; ok {

					log.Printf("found network %s with id %d, associating with interface with index %d", networkLabel, n.NetworkID, intIntfMap["interface_index"].(int))

					_, err = client.InstanceArrayInterfaceAttachNetwork(iaCreated.InstanceArrayID, intIntfMap["interface_index"].(int), n.NetworkID)
					if err != nil {
						return err
					}

				} else {
					return fmt.Errorf("could not find network with label %s", networkLabel)
				}
			}

			//create drive arrays
			daListIntf := d.Get(fmt.Sprintf("instance_array.%d.drive_array", i))

			for _, daMapIntf := range daListIntf.([]interface{}) {

				da := expandDriveArray(daMapIntf.(map[string]interface{}))

				da.InstanceArrayID = iaCreated.InstanceArrayID

				_, err := client.DriveArrayCreate(createdInfra.InfrastructureID, da)
				if err != nil {
					return err
				}
			}

		}
	}

	d.SetId(fmt.Sprintf("%d", createdInfra.InfrastructureID))

	log.Printf("current state object after create (before read):%v", d.Get("instance_array"))

	if preventDeploy, ok := d.GetOkExists("prevent_deploy"); !ok || preventDeploy == false {
		if err := deployInfrastructure(createdInfra.InfrastructureID, d, meta); err != nil {
			return err
		}

		if d.Get("await_deploy_finished").(bool) {
			return waitForInfrastructureFinished(createdInfra.InfrastructureID, d, meta, d.Timeout(schema.TimeoutCreate))
		}
	}

	return resourceInfrastructureRead(d, meta)

}

//resourceInfrastructureRead reads the serverside status of elements
//it ignores elements added outside of terraform (except of course at deploy time)
func resourceInfrastructureRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*mc.Client)

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

	//discover networks, we use the original network id so that we preserve the order of elements
	retNetworks, err := client.Networks(infrastructureID)
	if err != nil {
		return err
	}

	nList := []interface{}{}
	nListByID := map[int]mc.Network{}
	if nCount, ok := d.Get("network.#").(int); ok {
		for i := 0; i < nCount; i++ {
			if networkLabel, ok := d.GetOkExists(fmt.Sprintf("network.%d.network_label", i)); ok {
				n := (*retNetworks)[networkLabel.(string)]
				nList = append(nList, flattenNetwork(n))
				nListByID[n.NetworkID] = n
			}
		}
	}

	log.Printf("networks:%v", nList)

	d.Set("network", nList)

	//discover drive arrays
	retInstanceArrays, err := client.InstanceArrays(infrastructureID)
	if err != nil {
		return err
	}

	retDriveArrays, err := client.DriveArrays(infrastructureID)
	if err != nil {
		return err
	}

	iaList := []interface{}{}

	iaCount := d.Get("instance_array.#")

	for iai := 0; iai < iaCount.(int); iai++ {

		log.Printf("current state object=%v", d.Get(fmt.Sprintf("instance_array.%d", iai)))

		if instanceArrayLabel, ok := d.GetOkExists(fmt.Sprintf("instance_array.%d.instance_array_label", iai)); ok {

			log.Printf("Retriving existing instance array with label :%s", instanceArrayLabel.(string))

			//locate the instance array
			ia := (*retInstanceArrays)[fmt.Sprintf("%s.vanilla", instanceArrayLabel)]

			iaMap := flattenInstanceArray(ia)

			//get the drive arrays of the current instance array
			daList := []interface{}{}
			daCount := d.Get(fmt.Sprintf("instance_array.%d.drive_array.#", iai)).(int)

			log.Printf("daCount=%d", daCount)

			for di := 0; di < daCount; di++ {

				driveArrayLabel := d.Get(fmt.Sprintf("instance_array.%d.drive_array.%d.drive_array_label", iai, di))

				da := (*retDriveArrays)[fmt.Sprintf("%s.vanilla", driveArrayLabel)]
				daList = append(daList, flattenDriveArray(da))

			}

			log.Printf("daList=%d", daList)
			iaMap["drive_array"] = daList

			interfaces := ia.InstanceArrayInterfaces
			intfList := []interface{}{}
			//iterate over interfaces
			intfCount := d.Get(fmt.Sprintf("instance_array.%d.interface.#", iai)).(int)

			log.Printf("intfCount=%d", intfCount)

			for inti := 0; inti < intfCount; inti++ {
				log.Printf("inti=%d", inti)

				log.Printf("interface index for interface[%d] is %v", inti, d.Get(fmt.Sprintf("instance_array.%d.interface.%d.interface_index", iai, inti)))

				interfaceIndex := d.Get(fmt.Sprintf("instance_array.%d.interface.%d.interface_index", iai, inti))

				log.Printf("Looking for interface with index %d among instance arrays's interfaces %v", interfaceIndex, interfaces)

				//locate interface with index in returned data
				for _, intf := range interfaces {
					//if we found it, locate the network it's connected to add it to the list
					log.Printf("Currently at %d : %v", interfaceIndex, intf.InstanceArrayInterfaceIndex)

					if intf.InstanceArrayInterfaceIndex == interfaceIndex && intf.NetworkID != 0 {

						log.Printf("Interface connected to network %d", intf.NetworkID)

						//if we know the network this is attached to
						if n, ok := nListByID[intf.NetworkID]; ok {

							intfMap := map[string]interface{}{
								"interface_index": interfaceIndex,
								"network_label":   n.NetworkLabel,
								"network_id":      n.NetworkID,
							}

							intfList = append(intfList, intfMap)

							log.Printf("Appended interface %v", intfMap)

						} else {
							return fmt.Errorf("somehow an interface with id %d is connected to an inexistent network with id %d", intf.InstanceArrayInterfaceID, intf.NetworkID)
						}

					}
				}

			}

			iaMap["interface"] = intfList

			log.Printf("Appending instance array %v\n", iaMap)
			//finally append the instance array map to the list of instance arrays
			iaList = append(iaList, iaMap)

		}
	}

	j, _ := json.MarshalIndent(iaList, "", "\t")
	log.Printf("flattened list of instance arrays is now %s", j)

	if err := d.Set("instance_array", iaList); err != nil {
		return fmt.Errorf("error setting instance_array: %s", err)
	}

	return nil
}

//resourceInfrastructureUpdate applies changes on the serverside
//attempts to merge serverside changes into the current state
func resourceInfrastructureUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	client := meta.(*mc.Client)

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

		for i, iaMapIntf := range currentInstanceArraysMap {

			if !d.HasChange(fmt.Sprintf("instance_array.%d", i)) {
				continue
			}

			iaMap := iaMapIntf.(map[string]interface{})

			ia := expandInstanceArray(iaMap)

			//update interfaces
			intMapList := iaMap["interface"].([]interface{})

			var nMapList []interface{}

			if nMapListIntf, ok := d.GetOkExists("networks"); ok {
				nMapList = nMapListIntf.([]interface{})
			}

			intList := []mc.InstanceArrayInterface{}
			for ii, intMapIntf := range intMapList {

				if !d.HasChange(fmt.Sprintf("instance_array.%d.interface.%d", i, ii)) {
					continue
				}

				intMap := intMapIntf.(map[string]interface{})

				//because we could have alternations to the interface index - to network map
				//we're retrieving the networks rather than relying on the exisitng network_id
				//locate network with label and get it's network id
				var networkID = 0
				for _, nMapIntf := range nMapList {
					nMap := nMapIntf.(map[string]interface{})
					if nMap["network_label"] == intMap["network_label"] {
						networkID = nMap["network_id"].(int)
					}
				}

				intf := mc.InstanceArrayInterface{
					InstanceArrayInterfaceIndex: intMap["interface_index"].(int),
					NetworkID:                   networkID,
				}

				intList = append(intList, intf)

				needsDeploy = true
			}

			ia.InstanceArrayInterfaces = intList

			bkeepDetachingDrives := d.Get("keep_detaching_drives").(bool)
			bSwapExistingInstancesHardware := false

			retIA, err := createOrUpdateInstanceArray(infrastructureID, ia, client, &bSwapExistingInstancesHardware, &bkeepDetachingDrives, nil, nil)
			if err != nil {
				return err
			}

			//update drive arrays
			daList := iaMap["drive_array"].([]interface{})

			for di, daMap := range daList {
				if !d.HasChange(fmt.Sprintf("instance_array.%d.drive_array.%d", i, di)) {
					continue
				}
				da := expandDriveArray(daMap.(map[string]interface{}))
				da.InstanceArrayID = retIA.InstanceArrayID
				if _, err := createOrUpdateDriveArray(infrastructureID, da, client); err != nil {
					return err
				}
				needsDeploy = true
			}

			needsDeploy = true
		}
	}

	d.Partial(false)

	if needsDeploy {
		if preventDeploy, ok := d.GetOkExists("prevent_deploy"); !ok || preventDeploy == false {
			if err := deployInfrastructure(infrastructureID, d, meta); err != nil {
				return err
			}

			if d.Get("await_deploy_finished").(bool) {
				return waitForInfrastructureFinished(infrastructureID, d, meta, d.Timeout(schema.TimeoutUpdate))
			}
		}
	}

	return resourceInfrastructureRead(d, meta)
}

//does not wait for a deploy to finish.
func resourceInfrastructureDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*mc.Client)

	infrastructureID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	if err := client.InfrastructureDelete(infrastructureID); err != nil {
		return err
	}

	if preventDeploy, ok := d.GetOkExists("prevent_deploy"); !ok || preventDeploy == false {
		if err := deployInfrastructure(infrastructureID, d, meta); err != nil {
			return err
		}
		if d.Get("await_delete_finished").(bool) {
			return waitForInfrastructureFinished(infrastructureID, d, meta, d.Timeout(schema.TimeoutUpdate))
		}
	}

	d.SetId("")
	return nil
}

func createOrUpdateInstanceArray(infrastructureID int, ia mc.InstanceArray, client *mc.Client, bSwapExistingInstancesHardware *bool, bKeepDetachingDrives *bool, objServerTypeMatches *[]mc.ServerType, arrInstancesToBeDeleted *[]int) (*mc.InstanceArray, error) {
	var instanceArrayID = ia.InstanceArrayID

	var iaToReturn *mc.InstanceArray

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

		retIA2, err2 := client.InstanceArrayEdit(retIA.InstanceArrayID, *retIA.InstanceArrayOperation, bSwapExistingInstancesHardware, bKeepDetachingDrives, objServerTypeMatches, arrInstancesToBeDeleted)
		if err2 != nil {
			return nil, err2
		}
		iaToReturn = retIA2
	}
	return iaToReturn, nil
}

func createOrUpdateDriveArray(infrastructureID int, da mc.DriveArray, client *mc.Client) (*mc.DriveArray, error) {
	var driveArrayToReturn *mc.DriveArray
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

//waitForInfrastructureFinished awaits for the "finished" status in the specified infrastructure
func waitForInfrastructureFinished(infrastructureID int, d *schema.ResourceData, meta interface{}, timeout time.Duration) error {

	client := meta.(*mc.Client)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			"not_started",
			"ongoing",
		},
		Target: []string{
			"finished",
		},
		Refresh: func() (interface{}, string, error) {
			log.Printf("calling InfrastructureGet(%d) ...", infrastructureID)
			resp, err := client.InfrastructureGet(infrastructureID)
			if err != nil {
				return 0, "", err
			}
			return resp, resp.InfrastructureOperation.InfrastructureDeployStatus, nil
		},
		Timeout:                   timeout,
		Delay:                     30 * time.Second,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	if _, err := createStateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for example instance (%s) to be created: %s", d.Id(), err)
	}

	return resourceInfrastructureRead(d, meta)
}

//deployInfrastructure starts a deploy
func deployInfrastructure(infrastructureID int, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mc.Client)

	shutDownOptions := mc.ShutdownOptions{
		HardShutdownAfterTimeout:   d.Get("hard_shutdown_after_timeout").(bool),
		AttemptSoftShutdown:        d.Get("attempt_soft_shutdown").(bool),
		SoftShutdownTimeoutSeconds: d.Get("soft_shutdown_timeout_seconds").(int),
	}

	return client.InfrastructureDeploy(
		infrastructureID, shutDownOptions,
		d.Get("allow_data_loss").(bool),
		d.Get("skip_ansible").(bool),
	)
}
