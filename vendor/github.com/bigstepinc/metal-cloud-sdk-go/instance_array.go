package metalcloud

import "fmt"

//InstanceArray object describes a collection of identical instances
type InstanceArray struct {
	InstanceArrayID                 int                      `json:"instance_array_id,omitempty"`
	InstanceArrayLabel              string                   `json:"instance_array_label,omitempty"`
	InstanceArraySubdomain          string                   `json:"instance_array_subdomain,omitempty"`
	InstanceArrayBootMethod         string                   `json:"instance_array_boot_method,omitempty"`
	InstanceArrayInstanceCount      int                      `json:"instance_array_instance_count,omitempty"`
	InstanceArrayRAMGbytes          int                      `json:"instance_array_ram_gbytes,omitempty"`
	InstanceArrayProcessorCount     int                      `json:"instance_array_processor_count,omitempty"`
	InstanceArrayProcessorCoreMHZ   int                      `json:"instance_array_processor_core_mhz,omitempty"`
	InstanceArrayProcessorCoreCount int                      `json:"instance_array_processor_core_count,omitempty"`
	InstanceArrayDiskCount          int                      `json:"instance_array_disk_count,omitempty"`
	InstanceArrayDiskSizeMBytes     int                      `json:"instance_array_disk_size_mbytes,omitempty"`
	InstanceArrayDiskTypes          []string                 `json:"instance_array_disk_types,omitempty"`
	InfrastructureID                int                      `json:"infrastructure_id,omitempty"`
	InstanceArrayServiceStatus      string                   `json:"instance_array_service_status,omitempty"`
	InstanceArrayInterfaces         []InstanceArrayInterface `json:"instance_array_interfaces,omitempty"`
	ClusterID                       int                      `json:"cluster_id,omitempty"`
	ClusterRoleGroup                string                   `json:"cluster_role_group,omitempty"`
	InstanceArrayFirewallManaged    bool                     `json:"instance_array_firewall_managed,omitempty"`
	InstanceArrayFirewallRules      []FirewallRule           `json:"instance_array_firewall_rules,omitempty"`
	VolumeTemplateID                int                      `json:"volume_template_id,omitempty"`
	InstanceArrayOperation          *InstanceArrayOperation  `json:"instance_array_operation,omitempty"`
}

//InstanceArrayOperation object describes the changes that will be applied to an instance array
type InstanceArrayOperation struct {
	InstanceArrayID                 int                               `json:"instance_array_id,omitempty"`
	InstanceArrayLabel              string                            `json:"instance_array_label,omitempty"`
	InstanceArraySubdomain          string                            `json:"instance_array_subdomain,omitempty"`
	InstanceArrayBootMethod         string                            `json:"instance_array_boot_method,omitempty"`
	InstanceArrayInstanceCount      int                               `json:"instance_array_instance_count,omitempty"`
	InstanceArrayRAMGbytes          int                               `json:"instance_array_ram_gbytes,omitempty"`
	InstanceArrayProcessorCount     int                               `json:"instance_array_processor_count,omitempty"`
	InstanceArrayProcessorCoreMHZ   int                               `json:"instance_array_processor_core_mhz,omitempty"`
	InstanceArrayProcessorCoreCount int                               `json:"instance_array_processor_core_count,omitempty"`
	InstanceArrayDiskCount          int                               `json:"instance_array_disk_count,omitempty"`
	InstanceArrayDiskSizeMBytes     int                               `json:"instance_array_disk_size_mbytes,omitempty"`
	InstanceArrayDiskTypes          []string                          `json:"instance_array_disk_types,omitempty"`
	InstanceArrayServiceStatus      string                            `json:"instance_array_service_status,omitempty"`
	InstanceArrayInterfaces         []InstanceArrayInterfaceOperation `json:"instance_array_interfaces,omitempty"`
	ClusterID                       int                               `json:"cluster_id,omitempty"`
	ClusterRoleGroup                string                            `json:"cluster_role_group,omitempty"`
	InstanceArrayFirewallManaged    bool                              `json:"instance_array_firewall_managed,omitempty"`
	InstanceArrayFirewallRules      []FirewallRule                    `json:"instance_array_firewall_rules,omitempty"`
	VolumeTemplateID                int                               `json:"volume_template_id,omitempty"`
	InstanceArrayDeployType         string                            `json:"instance_array_deploy_type,omitempty"`
	InstanceArrayDeployStatus       string                            `json:"instance_array_deploy_status,omitempty"`
	InstanceArrayChangeID           int                               `json:"instance_array_change_id,omitempty"`
}

//FirewallRule describes a firewall rule that is to be applied on all instances of an array
type FirewallRule struct {
	FirewallRuleDescription                    string `json:"firewall_rule_description,omitempty"`
	FirewallRulePortRangeStart                 int    `json:"firewall_rule_port_range_start,omitempty"`
	FirewallRulePortRangeEnd                   int    `json:"firewall_rule_port_range_end,omitempty"`
	FirewallRuleSourceIPAddressRangeStart      string `json:"firewall_rule_source_ip_address_range_start,omitempty"`
	FirewallRuleSourceIPAddressRangeEnd        string `json:"firewall_rule_source_ip_address_range_end,omitempty"`
	FirewallRuleDestinationIPAddressRangeStart string `json:"firewall_rule_destination_ip_address_range_start,omitempty"`
	FirewallRuleDestinationIPAddressRangeEnd   string `json:"firewall_rule_destination_ip_address_range_end,omitempty"`
	FirewallRuleProtocol                       string `json:"firewall_rule_protocol,omitempty"`
	FirewallRuleIPAddressType                  string `json:"firewall_rule_ip_address_type,omitempty"`
	FirewallRuleEnabled                        bool   `json:"firewall_rule_enabled,omitempty"`
}

//InstanceArrayInterface describes a network interface of the array.
//It's properties will be applied to all InstanceInterfaces of the array's instances.
type InstanceArrayInterface struct {
	InstanceArrayInterfaceLabel            string                           `json:"instance_array_interface_label,omitempty"`
	InstanceArrayInterfaceSubdomain        string                           `json:"instance_array_interface_subdomain,omitempty"`
	InstanceArrayInterfaceID               int                              `json:"instance_array_interface_id,omitempty"`
	InstanceArrayID                        int                              `json:"instance_array_id,omitempty"`
	NetworkID                              int                              `json:"network_id,omitempty"`
	InstanceArrayInterfaceLAGGIndexes      []interface{}                    `json:"instance_array_interface_lagg_indexes,omitempty"`
	InstanceArrayInterfaceIndex            int                              `json:"instance_array_interface_index,omitempty"`
	InstanceArrayInterfaceServiceStatus    string                           `json:"instance_array_interface_service_status,omitempty"`
	InstanceArrayInterfaceCreatedTimestamp string                           `json:"instance_array_interface_created_timestamp,omitempty"`
	InstanceArrayInterfaceUpdatedTimestamp string                           `json:"instance_array_interface_updated_timestamp,omitempty"`
	InstanceArrayInterfaceOperation        *InstanceArrayInterfaceOperation `json:"instance_array_interface_operation,omitempty"`
}

//InstanceArrayInterfaceOperation describes changes to a network array interface
type InstanceArrayInterfaceOperation struct {
	InstanceArrayInterfaceLabel            string        `json:"instance_array_interface_label,omitempty"`
	InstanceArrayInterfaceSubdomain        string        `json:"instance_array_interface_subdomain,omitempty"`
	InstanceArrayInterfaceID               int           `json:"instance_array_interface_id,omitempty"`
	InstanceArrayID                        int           `json:"instance_array_id,omitempty"`
	NetworkID                              int           `json:"network_id,omitempty"`
	InstanceArrayInterfaceLAGGIndexes      []interface{} `json:"instance_array_interface_lagg_indexes,omitempty"`
	InstanceArrayInterfaceIndex            int           `json:"instance_array_interface_index,omitempty"`
	InstanceArrayInterfaceServiceStatus    string        `json:"instance_array_interface_service_status,omitempty"`
	InstanceArrayInterfaceCreatedTimestamp string        `json:"instance_array_interface_created_timestamp,omitempty"`
	InstanceArrayInterfaceUpdatedTimestamp string        `json:"instance_array_interface_updated_timestamp,omitempty"`
	InstanceArrayInterfaceChangeID         int           `json:"instance_array_interface_change_id,omitempty"`
}

//InstanceArrayGet returns an InstanceArray with specified id
func (c *Client) InstanceArrayGet(instanceArrayID ID) (*InstanceArray, error) {
	var createdObject InstanceArray

	if err := checkID(instanceArrayID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(
		&createdObject,
		"instance_array_get",
		instanceArrayID)

	if err != nil {
		return nil, err
	}

	return &createdObject, nil
}

//InstanceArrays returns list of instance arrays of specified infrastructure
func (c *Client) InstanceArrays(infrastructureID ID) (*map[string]InstanceArray, error) {

	if err := checkID(infrastructureID); err != nil {
		return nil, err
	}

	res, err := c.rpcClient.Call(
		"instance_arrays",
		infrastructureID)

	if err != nil {
		return nil, err
	}

	_, ok := res.Result.([]interface{})
	if ok {
		var m = map[string]InstanceArray{}
		return &m, nil
	}

	var createdObject map[string]InstanceArray

	err2 := res.GetObject(&createdObject)
	if err2 != nil {
		return nil, err2
	}

	return &createdObject, nil
}

//InstanceArrayCreate creates an instance array (colletion of identical instances). Requires Deploy.
func (c *Client) InstanceArrayCreate(infrastructureID ID, instanceArray InstanceArray) (*InstanceArray, error) {
	var createdObject InstanceArray

	if err := checkID(infrastructureID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(
		&createdObject,
		"instance_array_create",
		infrastructureID,
		instanceArray)

	if err != nil {
		return nil, err
	}

	return &createdObject, nil
}

//InstanceArrayEdit alterns a deployed instance array. Requires deploy.
func (c *Client) InstanceArrayEdit(instanceArrayID ID, instanceArrayOperation InstanceArrayOperation, bSwapExistingInstancesHardware *bool, bKeepDetachingDrives *bool, objServerTypeMatches *[]ServerType, arrInstancesToBeDeleted *[]int) (*InstanceArray, error) {
	var createdObject InstanceArray

	if err := checkID(instanceArrayID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(
		&createdObject,
		"instance_array_edit",
		instanceArrayID,
		instanceArrayOperation,
		bSwapExistingInstancesHardware,
		bKeepDetachingDrives,
		objServerTypeMatches,
		arrInstancesToBeDeleted)

	if err != nil {
		return nil, err
	}

	return &createdObject, nil
}

//InstanceArrayDelete deletes an instance array. Requires deploy.
func (c *Client) InstanceArrayDelete(instanceArrayID ID) error {

	if err := checkID(instanceArrayID); err != nil {
		return err
	}

	resp, err := c.rpcClient.Call(
		"instance_array_delete",
		instanceArrayID)

	if err != nil {
		return err
	}

	if resp.Error != nil {
		return fmt.Errorf(resp.Error.Message)
	}

	return nil
}

//InstanceArrayInterfaceAttachNetwork attaches an InstanceArrayInterface to a Network
func (c *Client) InstanceArrayInterfaceAttachNetwork(instanceArrayID int, instanceArrayInterfaceIndex int, networkID int) (*InstanceArray, error) {
	var createdObject InstanceArray

	err := c.rpcClient.CallFor(
		&createdObject,
		"instance_array_interface_attach_network",
		instanceArrayID,
		instanceArrayInterfaceIndex,
		networkID)

	if err != nil {
		return nil, err
	}

	return &createdObject, nil
}

//InstanceArrayInterfaceDetach detaches an InstanceArrayInterface from any Network element that is attached to.
func (c *Client) InstanceArrayInterfaceDetach(instanceArrayID int, instanceArrayInterfaceIndex int) (*InstanceArray, error) {
	var createdObject InstanceArray

	err := c.rpcClient.CallFor(
		&createdObject,
		"instance_array_interface_detach",
		instanceArrayID,
		instanceArrayInterfaceIndex)

	if err != nil {
		return nil, err
	}

	return &createdObject, nil
}
