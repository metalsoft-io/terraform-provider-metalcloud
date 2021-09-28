package metalcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
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

			"infrastructure_label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter_name": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_DATACENTER", nil),
				Optional:    true,
			},
			"infrastructure_custom_variables": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
			},
			"instance_array": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceInstanceArray(),
				Set:      instanceArrayResourceHash,
			},
			"network": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceNetwork(),
				Set:      networkResourceHash,
			},
			"prevent_deploy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"shared_drive": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceSharedDrive(),
				Set:      sharedDriveResourceHash,
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
			"instance_array_additional_wan_ipv4_json": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_array_custom_variables": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
			},
			"instance_custom_variables": {
				Type:     schema.TypeList,
				Elem:     instanceCustomVariableResource(),
				Optional: true,
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
				Set:      firewallRuleResourceHash,
			},
			"drive_array": {
				Type: schema.TypeSet,
				// Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceDriveArray(),
				Set:      driveArrayResourceHash,
			},

			"interface": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceInstanceArrayInterface(),
				Set:      interfaceResourceHash,
			},
			"instances": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//	Elem:     resourceInstanceArrayInstances(),
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

func instanceCustomVariableResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_index": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"custom_variables": &schema.Schema{
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Required: true,
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
			"network_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
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

func resourceSharedDrive() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"shared_drive_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"shared_drive_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"shared_drive_size_mbytes": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"shared_drive_storage_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared_drive_has_gfs": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"shared_drive_attached_instance_arrays": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	if cvIntf, ok := d.GetOkExists("infrastructure_custom_variables"); ok {
		cv := make(map[string]string)

		for k, v := range cvIntf.(map[string]interface{}) {
			cv[k] = v.(string)
		}

		infrastructure.InfrastructureCustomVariables = cv
	}

	createdInfra, err := client.InfrastructureCreate(infrastructure)
	if err != nil || createdInfra == nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", createdInfra.InfrastructureID))

	//create networks
	//there are some networks already created by infrastructure create so we will be updating their labels only
	retNetworks, err := client.Networks(createdInfra.InfrastructureID)
	if err != nil {
		err1 := resourceInfrastructureRead(d, meta)
		if err1 != nil {
			return err1
		}
		return err
	}

	nCreatedNetworksMap := map[string]mc.Network{}
	if networks, ok := d.GetOk("network"); ok {
		for _, resN := range networks.(*schema.Set).List() {

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
							err1 := resourceInfrastructureRead(d, meta)
							if err1 != nil {
								return err1
							}
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
					err1 := resourceInfrastructureRead(d, meta)
					if err1 != nil {
						return err1
					}
					return err
				}
			}
			nCreatedNetworksMap[network.NetworkLabel] = *network
		}
	}

	iaInfraMap := make(map[string]mc.InstanceArray)

	//create instance arrays (and their drives)
	if instanceArrays, ok := d.GetOkExists("instance_array"); ok {
		for _, resIA := range instanceArrays.(*schema.Set).List() {
			iaMap := resIA.(map[string]interface{})
			//populate instance array
			ia := expandInstanceArray(iaMap)

			//create the instance array
			iaCreated, err := client.InstanceArrayCreate(createdInfra.InfrastructureID, ia)
			if err != nil {
				err1 := resourceInfrastructureRead(d, meta)
				if err1 != nil {
					return err1
				}
				return err
			}

			iaInfraMap[ia.InstanceArrayLabel] = *iaCreated

			cvList := iaMap["instance_custom_variables"].([]interface{})
			instanceList, err := client.InstanceArrayInstances(iaCreated.InstanceArrayID)
			if err != nil {
				return err
			}

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

			for _, icvIntf := range cvList {
				icv := icvIntf.(map[string]interface{})
				cvIntf := icv["custom_variables"].(map[string]interface{})
				instance_custom_variables := make(map[string]string)
				for k, v := range cvIntf {
					instance_custom_variables[k] = v.(string)
				}
				instance_index := icv["instance_index"].(int)
				if instance_index < nInstances {
					instance := instances[instance_index]
					instance.InstanceOperation.InstanceCustomVariables = instance_custom_variables
					_, err := client.InstanceEdit(instance.InstanceID, instance.InstanceOperation)
					if err != nil {
						return err
					}
				}
			}

			//add interfaces to networks

			intListIntf := d.Get(fmt.Sprintf("instance_array.%d.interface", instanceArrayResourceHash(resIA)))

			for _, intIntf := range intListIntf.(*schema.Set).List() {
				intIntfMap := intIntf.(map[string]interface{})
				networkLabel := intIntfMap["network_label"].(string)

				//look for a network with the given label
				if n, ok := nCreatedNetworksMap[networkLabel]; ok {

					log.Printf("found network %s with id %d, associating with interface with index %d", networkLabel, n.NetworkID, intIntfMap["interface_index"].(int))

					_, err = client.InstanceArrayInterfaceAttachNetwork(iaCreated.InstanceArrayID, intIntfMap["interface_index"].(int), n.NetworkID)
					if err != nil {
						err1 := resourceInfrastructureRead(d, meta)
						if err1 != nil {
							return err1
						}
						return err
					}

				} else {
					return fmt.Errorf("could not find network with label %s", networkLabel)
				}
			}

			//create drive arrays
			daListIntf := d.Get(fmt.Sprintf("instance_array.%d.drive_array", instanceArrayResourceHash(resIA)))

			for _, daMapIntf := range daListIntf.(*schema.Set).List() {
				da := expandDriveArray(daMapIntf.(map[string]interface{}))

				da.InstanceArrayID = iaCreated.InstanceArrayID

				_, err := client.DriveArrayCreate(createdInfra.InfrastructureID, da)
				if err != nil {
					err1 := resourceInfrastructureRead(d, meta)
					if err1 != nil {
						return err1
					}
					return err
				}
			}
		}
	}

	//create shared drives
	if sharedDrives, ok := d.GetOkExists("shared_drive"); ok {
		for _, resSD := range sharedDrives.(*schema.Set).List() {
			//populate shared drive
			retInstanceArrays, err := client.InstanceArrays(createdInfra.InfrastructureID)
			if err != nil {
				return err
			}

			sdMap := resSD.(map[string]interface{})
			sdMap["infrastructure_instance_arrays_planned"] = *retInstanceArrays
			sdMap["infrastructure_instance_arrays_existing"] = iaInfraMap
			sd := expandSharedDrive(sdMap)
			//create shared drive
			_, err = client.SharedDriveCreate(createdInfra.InfrastructureID, sd)
			if err != nil {
				err1 := resourceInfrastructureRead(d, meta)
				if err1 != nil {
					return err1
				}
				return err
			}
		}
	}

	log.Printf("current state object after create (before read):%v", d.Get("instance_array"))

	if preventDeploy := d.Get("prevent_deploy"); preventDeploy == false {
		if err := deployInfrastructure(createdInfra.InfrastructureID, d, meta); err != nil {
			err1 := resourceInfrastructureRead(d, meta)
			if err1 != nil {
				return err1
			}
			return err
		}

		if d.Get("await_deploy_finished").(bool) {
			return waitForInfrastructureFinished(createdInfra.InfrastructureID, d, meta, d.Timeout(schema.TimeoutCreate), DEPLOY_STATUS_FINISHED)
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

	switch infrastructure.InfrastructureCustomVariables.(type) {
	case []interface{}:
		err := d.Set("infrastructure_custom_variables", make(map[string]string))
		if err != nil {
			return fmt.Errorf("error setting infrastructure custom variables %s", err)
		}
	default:
		icv := make(map[string]string)

		for k, v := range infrastructure.InfrastructureCustomVariables.(map[string]interface{}) {
			icv[k] = v.(string)
		}
		err := d.Set("infrastructure_custom_variables", icv)

		if err != nil {
			return fmt.Errorf("error setting infrastructure custom variables %s", err)
		}
	}

	//discover networks, we use the original network id so that we preserve the order of elements
	retNetworks, err := client.Networks(infrastructureID)
	if err != nil {
		return err
	}

	nSet := schema.NewSet(networkResourceHash, []interface{}{})
	nListByID := make(map[int]mc.Network, len(*retNetworks))

	if networks, ok := d.GetOk("network"); ok {
		for _, network := range networks.(*schema.Set).List() {
			networkMap := network.(map[string]interface{})
			if networkLabel, ok := networkMap["network_label"]; ok {
				n, ok := (*retNetworks)[networkLabel.(string)]
				if !ok {
					continue
				}
				nSet.Add(flattenNetwork(n))
				nListByID[n.NetworkID] = n
			}
		}
	}

	err = d.Set("network", nSet)
	if err != nil {
		return err
	}

	//discover drive arrays
	retInstanceArrays, err := client.InstanceArrays(infrastructureID)
	if err != nil {
		return err
	}

	retDriveArrays, err := client.DriveArrays(infrastructureID)
	if err != nil {
		return err
	}

	iaSet := schema.NewSet(instanceArrayResourceHash, []interface{}{})

	if instanceArrays, ok := d.GetOk("instance_array"); ok {
		for _, iaIntf := range instanceArrays.(*schema.Set).List() {
			instance_array := iaIntf.(map[string]interface{})

			if instanceArrayLabel, ok := instance_array["instance_array_label"]; ok {
				log.Printf("Retriving existing instance array with label :%s", instanceArrayLabel.(string))

				//locate the instance array
				ia, ok := (*retInstanceArrays)[fmt.Sprintf("%s.vanilla", instanceArrayLabel)]

				if !ok {
					continue
				}

				iaMap := flattenInstanceArray(ia)

				retInstances, err := client.InstanceArrayInstances(ia.InstanceArrayID)
				if err != nil {
					return err
				}

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

				bytes, err := json.Marshal(retInstances)
				if err != nil {
					return fmt.Errorf("error serializing instances array: %s", err)
				}

				iaMap["instances"] = string(bytes)

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

				if len(customVars) > 0 {
					iaMap["instance_custom_variables"] = customVars

				}

				//get the drive arrays of the current instance array
				daSet := schema.NewSet(driveArrayResourceHash, []interface{}{})

				if driveArrays, ok := instance_array["drive_array"]; ok {

					for _, daIntf := range driveArrays.(*schema.Set).List() {
						drive_array := daIntf.(map[string]interface{})

						driveArrayLabel := drive_array["drive_array_label"]

						da, ok := (*retDriveArrays)[fmt.Sprintf("%s.vanilla", driveArrayLabel)]
						if !ok {
							continue
						}
						daSet.Add(flattenDriveArray(da))

					}
				}

				log.Printf("daList=%d", daSet.List())
				iaMap["drive_array"] = daSet

				interfacesRes := ia.InstanceArrayInterfaces
				intfSet := schema.NewSet(interfaceResourceHash, []interface{}{})

				//iterate over interfaces
				if interfaces, ok := instance_array["interface"]; ok {

					for _, iIntf := range interfaces.(*schema.Set).List() {
						iaInterface := iIntf.(map[string]interface{})
						interfaceIndex := iaInterface["interface_index"]

						log.Printf("Looking for interface with index %d among instance arrays's interfaces %v", interfaceIndex, interfaces)

						//locate interface with index in returned data
						for _, intf := range interfacesRes {
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

									intfSet.Add(intfMap)

									log.Printf("Appended interface %v", intfMap)

								} else {
									return fmt.Errorf("somehow an interface with id %d is connected to an inexistent network with id %d", intf.InstanceArrayInterfaceID, intf.NetworkID)
								}

							}
						}

					}
				}

				iaMap["interface"] = intfSet

				log.Printf("Appending instance array %v\n", iaMap)
				//finally append the instance array map to the list of instance arrays
				iaSet.Add(iaMap)

			}
		}
	}

	if err := d.Set("instance_array", iaSet); err != nil {
		return fmt.Errorf("error setting instance_array: %s", err)
	}

	//discover shared drives
	retSharedDrives, err := client.SharedDrives(infrastructureID)
	if err != nil {
		return err
	}

	sdSet := schema.NewSet(sharedDriveResourceHash, []interface{}{})

	if sharedDrives, ok := d.GetOk("shared_drive"); ok {
		for _, sd := range sharedDrives.(*schema.Set).List() {
			sdMap := sd.(map[string]interface{})
			if sdLabel, ok := sdMap["shared_drive_label"]; ok {
				sd, ok := (*retSharedDrives)[sdLabel.(string)]
				if !ok {
					continue
				}

				sdAttIAs := []interface{}{}
				for _, value := range sd.SharedDriveAttachedInstanceArrays {
					for _, ia := range *retInstanceArrays {
						if value == ia.InstanceArrayID {
							sdAttIAs = append(sdAttIAs, ia.InstanceArrayLabel)
						}
					}
				}

				sdResult := flattenSharedDrive(sd, sdAttIAs)
				sdSet.Add(sdResult)
			}
		}
	}

	log.Printf("shared drives:%v", sdSet.List())

	d.Set("shared_drive", sdSet)

	return nil
}

//resourceInfrastructureUpdate applies changes on the serverside
//attempts to merge serverside changes into the current state
func resourceInfrastructureUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	client := meta.(*mc.Client)

	infrastructureID, err := strconv.Atoi(d.Id())
	if err != nil {
		err1 := resourceInfrastructureRead(d, meta)
		if err1 != nil {
			return err1
		}
		return err
	}

	needsDeploy := false

	if d.HasChange("infrastructure_label") || d.HasChange("datacenter_name") || d.HasChange("infrastructure_custom_variables") {
		infrastructure, err := client.InfrastructureGet(int(infrastructureID))
		if err != nil {
			err1 := resourceInfrastructureRead(d, meta)
			if err1 != nil {
				return err1
			}
			return err
		}

		operation := infrastructure.InfrastructureOperation
		operation.InfrastructureLabel = d.Get("infrastructure_label").(string)
		operation.DatacenterName = d.Get("datacenter_name").(string)

		if cvIntf, ok := d.GetOkExists("infrastructure_custom_variables"); ok {
			cv := make(map[string]string)

			for k, v := range cvIntf.(map[string]interface{}) {
				cv[k] = v.(string)
			}

			operation.InfrastructureCustomVariables = cv
		}

		if _, err = client.InfrastructureEdit(int(infrastructureID), operation); err != nil {
			err1 := resourceInfrastructureRead(d, meta)
			if err1 != nil {
				return err1
			}
			return err
		}

		needsDeploy = true
	}

	retInstanceArrays, err := client.InstanceArrays(infrastructureID)
	if err != nil {
		err1 := resourceInfrastructureRead(d, meta)
		if err1 != nil {
			return err1
		}
		return err
	}

	retDriveArraysMap, err := client.DriveArrays(infrastructureID)
	if err != nil {
		err1 := resourceInfrastructureRead(d, meta)
		if err1 != nil {
			return err1
		}
		return err
	}

	retSharedDrivesMap, err := client.SharedDrives(infrastructureID)
	if err != nil {
		err1 := resourceInfrastructureRead(d, meta)
		if err1 != nil {
			return err1
		}
		return err
	}

	retNetworksMap, err := client.Networks(infrastructureID)
	if err != nil {
		err1 := resourceInfrastructureRead(d, meta)
		if err1 != nil {
			return err1
		}
		return err
	}

	iaInfraMap := make(map[string]mc.InstanceArray)
	stateInstanceArrayMap := make(map[int]*mc.InstanceArray)
	stateDriveArrayMap := make(map[int]*mc.DriveArray)
	stateNetworksMap := make(map[int]*mc.Network)

	if d.HasChange("network") {
		nMapByLabel := d.Get("network").(*schema.Set).List()

		for _, nMapIntf := range nMapByLabel {
			nMap := nMapIntf.(map[string]interface{})
			n := expandNetwork(nMap)
			stateNetworksMap[n.NetworkID] = &n

			if _, err := createOrUpdateNetwork(infrastructureID, n, client); err != nil {
				err1 := resourceInfrastructureRead(d, meta)
				if err1 != nil {
					return err1
				}
				return err
			}
			needsDeploy = true
		}

		for _, v := range *retNetworksMap {
			if _, ok := stateNetworksMap[v.NetworkID]; !ok && v.NetworkType == NETWORK_TYPE_LAN {
				needsDeploy = true
				err := deleteNetwork(&v, client)
				if err != nil {
					err1 := resourceInfrastructureRead(d, meta)
					if err1 != nil {
						return err1
					}
					return err
				}
			}
		}
	}

	if d.HasChange("instance_array") {
		//take each instance array and apply changes
		currentInstanceArraysMap := d.Get("instance_array").(*schema.Set).List()
		labelList := map[string]int{}

		for _, iaMapIntf := range currentInstanceArraysMap {
			iaMap := iaMapIntf.(map[string]interface{})

			label, _ := iaMap["instance_array_label"]
			labelList[label.(string)] = len(label.(string))
		}

		for _, iaMapIntf := range currentInstanceArraysMap {
			iaMap := iaMapIntf.(map[string]interface{})

			if _, ok := iaMap["instance_array_label"].(string); !ok {
				continue
			}

			if ok := len(iaMap["instance_array_label"].(string)); ok == 0 {
				continue
			}

			retIA := &mc.InstanceArray{}
			ia := expandInstanceArray(iaMap)

			if iaRes, ok := (*retInstanceArrays)[fmt.Sprintf("%s.vanilla", ia.InstanceArrayLabel)]; ok {
				ia.InstanceArrayID = iaRes.InstanceArrayID
				stateInstanceArrayMap[ia.InstanceArrayID] = &ia
			}

			//update interfaces
			intMapList := iaMap["interface"].(*schema.Set).List()

			var nMapList []interface{}

			if nMapListIntf, ok := d.GetOkExists("networks"); ok {
				nMapList = nMapListIntf.(*schema.Set).List()
			}

			intList := []mc.InstanceArrayInterface{}
			for _, intMapIntf := range intMapList {
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
					intf := mc.InstanceArrayInterface{
						InstanceArrayInterfaceIndex: intMap["interface_index"].(int),
						NetworkID:                   networkID,
					}

					intList = append(intList, intf)

					needsDeploy = true
				}
			}

			ia.InstanceArrayInterfaces = intList

			bkeepDetachingDrives := d.Get("keep_detaching_drives").(bool)
			bSwapExistingInstancesHardware := false

			retIA, err = createOrUpdateInstanceArray(infrastructureID, ia, client, &bSwapExistingInstancesHardware, &bkeepDetachingDrives, nil, nil)

			if err != nil {
				err1 := resourceInfrastructureRead(d, meta)
				if err1 != nil {
					return err1
				}
				return err
			}

			cvList := iaMap["instance_custom_variables"].([]interface{})
			instanceList, err := client.InstanceArrayInstances(ia.InstanceArrayID)
			if err != nil {
				return err
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

			currentCVLabelList := make(map[string]int, len(*instanceList))

			for _, icvIntf := range cvList {
				icv := icvIntf.(map[string]interface{})
				cvIntf := icv["custom_variables"].(map[string]interface{})
				instance_custom_variables := make(map[string]string)
				for k, v := range cvIntf {
					instance_custom_variables[k] = v.(string)
				}
				instance_index := icv["instance_index"].(int)
				if instance_index < nInstances {
					instance := instances[instance_index]
					currentCVLabelList[instance.InstanceLabel] = instance.InstanceID
					instance.InstanceOperation.InstanceCustomVariables = instance_custom_variables
					_, err := client.InstanceEdit(instance.InstanceID, instance.InstanceOperation)
					if err != nil {
						return err
					}
				}
			}

			for _, instance := range *instanceList {
				if _, ok := currentCVLabelList[instance.InstanceLabel]; !ok {
					instance.InstanceOperation.InstanceCustomVariables = make(map[string]string)
					_, err := client.InstanceEdit(instance.InstanceID, instance.InstanceOperation)
					if err != nil {
						return err
					}
				}
			}

			iaInfraMap[ia.InstanceArrayLabel] = *retIA

			//update drive arrays
			daList := iaMap["drive_array"].(*schema.Set).List()

			for _, daMapIntf := range daList {
				daMap := daMapIntf.(map[string]interface{})
				da := expandDriveArray(daMap)
				if daRes, ok := (*retDriveArraysMap)[fmt.Sprintf("%s.vanilla", da.DriveArrayLabel)]; ok {
					da.DriveArrayID = daRes.DriveArrayID
					stateDriveArrayMap[da.DriveArrayID] = &da
				}

				if ia.InstanceArrayID != 0 {
					da.InstanceArrayID = ia.InstanceArrayID
				} else {
					da.InstanceArrayID = retIA.InstanceArrayID
				}

				if _, err := createOrUpdateDriveArray(infrastructureID, da, client); err != nil {
					err1 := resourceInfrastructureRead(d, meta)
					if err1 != nil {
						return err1
					}
					return err
				}
				needsDeploy = true
			}

			needsDeploy = true
		}
	}

	stateSharedDriveMap := make(map[string]*mc.SharedDrive)

	if d.HasChange("shared_drive") {
		//update shared drives
		sdList := d.Get("shared_drive").(*schema.Set).List()

		for _, sdMapIntf := range sdList {
			sdMap := sdMapIntf.(map[string]interface{})
			sdMap["infrastructure_instance_arrays_planned"] = *retInstanceArrays
			sdMap["infrastructure_instance_arrays_existing"] = iaInfraMap
			sd := expandSharedDrive(sdMap)

			if sdRes, ok := (*retSharedDrivesMap)[sd.SharedDriveLabel]; ok {
				sd.SharedDriveID = sdRes.SharedDriveID
				stateSharedDriveMap[sd.SharedDriveLabel] = &sd
			}

			if _, err := createOrUpdateSharedDrive(infrastructureID, sd, client); err != nil {
				err1 := resourceInfrastructureRead(d, meta)
				if err1 != nil {
					return err1
				}
				return err
			}
			needsDeploy = true
		}

		for k, v := range *retSharedDrivesMap {
			if _, ok := stateSharedDriveMap[k]; !ok {
				needsDeploy = true
				err := deleteSharedDrive(&v, client)
				if err != nil {
					err1 := resourceInfrastructureRead(d, meta)
					if err1 != nil {
						return err1
					}
					return err
				}
			}
		}
	}

	for _, v := range *retDriveArraysMap {
		if _, ok := stateDriveArrayMap[v.DriveArrayID]; !ok {
			needsDeploy = true
			err := deleteDriveArray(&v, client)
			if err != nil {
				err1 := resourceInfrastructureRead(d, meta)
				if err1 != nil {
					return err1
				}
				return err
			}
		}
	}

	for _, v := range *retInstanceArrays {
		if _, ok := stateInstanceArrayMap[v.InstanceArrayID]; !ok {
			needsDeploy = true
			err := deleteInstanceArray(&v, client)
			if err != nil {
				err1 := resourceInfrastructureRead(d, meta)
				if err1 != nil {
					return err1
				}
				return err
			}
		}
	}

	d.Partial(false)

	if needsDeploy {
		if preventDeploy := d.Get("prevent_deploy"); preventDeploy == false {
			if err := deployInfrastructure(infrastructureID, d, meta); err != nil {
				err1 := resourceInfrastructureRead(d, meta)
				if err1 != nil {
					return err1
				}
				return err
			}

			if d.Get("await_deploy_finished").(bool) {
				return waitForInfrastructureFinished(infrastructureID, d, meta, d.Timeout(schema.TimeoutUpdate), DEPLOY_STATUS_FINISHED)
			}
		}
	}
	err = resourceInfrastructureRead(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func deleteInstanceArray(ia *mc.InstanceArray, client *mc.Client) error {
	err := client.InstanceArrayDelete(ia.InstanceArrayID)
	if err != nil {
		return err
	}
	return nil
}

func deleteSharedDrive(sd *mc.SharedDrive, client *mc.Client) error {
	err := client.SharedDriveDelete(sd.SharedDriveID)
	if err != nil {
		return err
	}
	return nil
}

func deleteDriveArray(da *mc.DriveArray, client *mc.Client) error {
	err := client.DriveArrayDelete(da.DriveArrayID)
	if err != nil {
		return err
	}
	return nil
}

func deleteNetwork(n *mc.Network, client *mc.Client) error {
	err := client.NetworkDelete(n.NetworkID)
	if err != nil {
		return err
	}
	return nil
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

	if preventDeploy := d.Get("prevent_deploy"); preventDeploy == false {
		if err := deployInfrastructure(infrastructureID, d, meta); err != nil {
			return err
		}
		if d.Get("await_delete_finished").(bool) {
			return waitForInfrastructureFinished(infrastructureID, d, meta, d.Timeout(schema.TimeoutUpdate), DEPLOY_STATUS_DELETED)
		}
	}

	d.SetId("")
	return nil
}
func createOrUpdateNetwork(infrastructureID int, n mc.Network, client *mc.Client) (*mc.Network, error) {
	var networkID = n.NetworkID
	var nToReturn *mc.Network

	if networkID == 0 {
		networks, err := client.Networks(infrastructureID)
		if err != nil {
			return nil, err
		}
		bExists := false
		for _, v := range *networks {
			if v.NetworkType == n.NetworkType &&
				(v.NetworkType == NETWORK_TYPE_SAN ||
					v.NetworkType == NETWORK_TYPE_WAN) {
				v.NetworkOperation.NetworkLabel = n.NetworkLabel
				nToReturn, err = client.NetworkEdit(networkID, *v.NetworkOperation)
				if err != nil {
					return nil, err
				}
				bExists = true

			}
		}

		if !bExists {
			nToReturn, err = client.NetworkCreate(infrastructureID, n)
			if err != nil {
				return nil, err
			}
		}

		return nToReturn, nil
	}

	retN, err := client.NetworkGet(networkID)
	if err != nil {
		return nil, err
	}

	retN.NetworkOperation.NetworkLabel = n.NetworkLabel
	retN, err = client.NetworkEdit(networkID, *retN.NetworkOperation)
	if err != nil {
		return nil, err
	}

	return retN, nil
}

func networkResourceHash(v interface{}) int {
	var buf bytes.Buffer
	n := v.(map[string]interface{})

	network_label := n["network_label"].(string)
	network_type := n["network_type"].(string)

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(network_label)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(network_type)))

	return hash(buf.String())
}

func sharedDriveResourceHash(v interface{}) int {
	var buf bytes.Buffer

	sd := v.(map[string]interface{})

	shared_drive_label := sd["shared_drive_label"].(string)
	shared_drive_storage_type := sd["shared_drive_storage_type"].(string)
	shared_drive_size_mbytes := strconv.Itoa(sd["shared_drive_size_mbytes"].(int))
	iaList := sd["shared_drive_attached_instance_arrays"].([]interface{})

	shared_drive_attached_instance_arrays := make([]string, len(iaList))

	for _, iaLabel := range iaList {
		shared_drive_attached_instance_arrays = append(shared_drive_attached_instance_arrays, iaLabel.(string))
	}

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(shared_drive_label)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(shared_drive_storage_type)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(shared_drive_size_mbytes)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(strings.Join(shared_drive_attached_instance_arrays, "-"))))

	return hash(buf.String())
}

func instanceArrayResourceHash(v interface{}) int {
	var buf bytes.Buffer
	ia := v.(map[string]interface{})

	instance_array_label := ia["instance_array_label"].(string)
	instance_array_boot_method := ia["instance_array_boot_method"].(string)
	instance_array_instance_count := strconv.Itoa(ia["instance_array_instance_count"].(int))
	instance_array_ram_gbytes := strconv.Itoa(ia["instance_array_ram_gbytes"].(int))
	instance_array_processor_count := strconv.Itoa(ia["instance_array_processor_count"].(int))
	instance_array_processor_core_mhz := strconv.Itoa(ia["instance_array_processor_core_mhz"].(int))
	instance_array_disk_count := strconv.Itoa(ia["instance_array_disk_count"].(int))
	instance_array_disk_size_mbytes := strconv.Itoa(ia["instance_array_disk_size_mbytes"].(int))
	volume_template_id := strconv.Itoa(ia["volume_template_id"].(int))
	instance_array_firewall_managed := strconv.FormatBool(ia["instance_array_firewall_managed"].(bool))
	instance_array_custom_variables, _ := json.Marshal(ia["instance_array_custom_variables"])

	var instance_custom_variables []byte

	if ia["instance_custom_variables"] != nil {
		for _, iaIntf := range ia["instance_custom_variables"].([]interface{}) {
			iacv := iaIntf.(map[string]interface{})

			cv := make(map[string]string)
			custom_variables := iacv["custom_variables"].(map[string]interface{})

			for k, v := range custom_variables {
				cv[k] = v.(string)
			}
			cv["index"] = strconv.Itoa(iacv["instance_index"].(int))

			instance_custom_variables, _ = json.Marshal(cv)
		}
	}

	drive_arrays := ia["drive_array"].(*schema.Set).List()
	for _, driveArray := range drive_arrays {
		buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(driveArrayToString(driveArray))))
	}

	interfaces := ia["interface"].(*schema.Set).List()
	for _, intf := range interfaces {
		buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(interfaceToString(intf))))
	}

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_label)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_boot_method)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_instance_count)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_ram_gbytes)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_processor_count)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_processor_core_mhz)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_disk_count)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_disk_size_mbytes)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(volume_template_id)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_firewall_managed)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(string(instance_array_custom_variables))))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(string(instance_custom_variables))))

	return hash(buf.String())

}

func driveArrayToString(v interface{}) string {
	var buf bytes.Buffer

	da := v.(map[string]interface{})
	drive_array_label := da["drive_array_label"].(string)
	drive_array_storage_type := da["drive_array_storage_type"].(string)
	drive_size_mbytes_default := strconv.Itoa(da["drive_size_mbytes_default"].(int))
	volume_template_id := strconv.Itoa(da["volume_template_id"].(int))

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(drive_array_label)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(drive_array_storage_type)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(drive_size_mbytes_default)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(volume_template_id)))

	return buf.String()
}

func driveArrayResourceHash(v interface{}) int {
	return hash(driveArrayToString(v))
}

func interfaceToString(v interface{}) string {
	var buf bytes.Buffer

	i := v.(map[string]interface{})

	instance_array_interface_label := i["instance_array_interface_label"].(string)
	instance_array_interface_service_status := i["instance_array_interface_service_status"].(string)
	instance_array_interface_index := strconv.Itoa(i["instance_array_interface_index"].(int))
	network_id := strconv.Itoa(i["network_id"].(int))

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_interface_label)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_interface_service_status)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(instance_array_interface_index)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(network_id)))

	return buf.String()
}

func interfaceResourceHash(v interface{}) int {
	return hash(interfaceToString(v))
}

func firewallRuleResourceHash(v interface{}) int {
	var buf bytes.Buffer
	fr := v.(map[string]interface{})

	firewall_rule_description := fr["firewall_rule_description"].(string)
	firewall_rule_source_ip_address_range_start := fr["firewall_rule_source_ip_address_range_start"].(string)
	firewall_rule_source_ip_address_range_end := fr["firewall_rule_source_ip_address_range_end"].(string)
	firewall_rule_destination_ip_address_range_start := fr["firewall_rule_destination_ip_address_range_start"].(string)
	firewall_rule_destination_ip_address_range_end := fr["firewall_rule_destination_ip_address_range_end"].(string)
	firewall_rule_protocol := fr["firewall_rule_protocol"].(string)
	firewall_rule_ip_address_type := fr["firewall_rule_ip_address_type"].(string)
	firewall_rule_port_range_start := strconv.Itoa(fr["firewall_rule_port_range_start"].(int))
	firewall_rule_port_range_end := strconv.Itoa(fr["firewall_rule_port_range_end"].(int))
	firewall_rule_enabled := strconv.FormatBool(fr["firewall_rule_enabled"].(bool))

	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_description)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_source_ip_address_range_start)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_source_ip_address_range_end)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_destination_ip_address_range_start)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_destination_ip_address_range_end)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_protocol)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_ip_address_type)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_port_range_start)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_port_range_end)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(firewall_rule_enabled)))

	return hash(buf.String())
}

func hash(v string) int {
	hash := crc32.ChecksumIEEE([]byte(v))

	return int(hash)

}

func createOrUpdateInstanceArray(infrastructureID int, ia mc.InstanceArray, client *mc.Client, bSwapExistingInstancesHardware *bool, bKeepDetachingDrives *bool, objServerTypeMatches *mc.ServerTypeMatches, arrInstancesToBeDeleted *[]int) (*mc.InstanceArray, error) {
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

func createOrUpdateSharedDrive(infrastructureID int, sd mc.SharedDrive, client *mc.Client) (*mc.SharedDrive, error) {
	var sharedDriveToReturn *mc.SharedDrive
	if sd.SharedDriveID == 0 {
		retSD, err := client.SharedDriveCreate(infrastructureID, sd)
		if err != nil {
			return nil, err
		}
		sharedDriveToReturn = retSD
	} else {
		retSD, err := client.SharedDriveGet(sd.SharedDriveID)
		if err != nil {
			return nil, err
		}

		copySharedDriveToOperation(sd, &retSD.SharedDriveOperation)

		retSD, err = client.SharedDriveEdit(sd.SharedDriveID, *&retSD.SharedDriveOperation)
		if err != nil {
			return nil, err
		}
		sharedDriveToReturn = retSD
	}

	return sharedDriveToReturn, nil
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
func waitForInfrastructureFinished(infrastructureID int, d *schema.ResourceData, meta interface{}, timeout time.Duration, targetStatus string) error {

	client := meta.(*mc.Client)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			DEPLOY_STATUS_NOT_STARTED,
			DEPLOY_STATUS_ONGOING,
		},
		Target: []string{
			targetStatus,
		},
		Refresh: func() (interface{}, string, error) {
			log.Printf("calling InfrastructureGet(%d) ...", infrastructureID)
			resp, err := client.InfrastructureGet(infrastructureID)
			if err != nil {
				if targetStatus == DEPLOY_STATUS_DELETED {
					return 0, targetStatus, nil
				}
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

const DEPLOY_STATUS_FINISHED = "finished"
const DEPLOY_STATUS_ONGOING = "ongoing"
const DEPLOY_STATUS_DELETED = "deleted"
const DEPLOY_STATUS_NOT_STARTED = "not_started"
const NETWORK_TYPE_LAN = "lan"
const NETWORK_TYPE_SAN = "san"
const NETWORK_TYPE_WAN = "wan"
