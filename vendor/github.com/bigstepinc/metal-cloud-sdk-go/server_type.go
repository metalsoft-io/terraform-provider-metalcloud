package metalcloud

//ServerType struct defines a server type
type ServerType struct {
	ServerTypeID                   int    `json:"server_type_id,omitempty"`
	ServerNetworkTotalCapacityMBps int    `json:"server_network_total_capacity_mbps,omitempty"`
	ServerTypeName                 string `json:"server_type_name,omitempty"`
	ServerTypeDisplayName          string `json:"server_type_display_name,omitempty"`
	ServerTypeLabel                string `json:"server_type_label,omitempty"`
	ServerProcessorCoreCount       int    `json:"server_processor_core_count,omitempty"`
	ServerProcessorCoreMHz         int    `json:"server_processor_core_mhz,omitempty"`
	ServerProcessorCount           int    `json:"server_processor_count,omitempty"`
	ServerRAMGbytes                int    `json:"server_ram_gbytes,omitempty"`
	ServerDiskCount                int    `json:"server_disk_count,omitempty"`
	ServerDiskType                 string `json:"server_disk_type,omitempty"`
	ServerDiskSizeMBytes           int    `json:"server_disk_size_mbytes,omitempty"`
	ServerProcessorNamesJSON       string `json:"server_processor_names_json,omitempty"`
	ServerProcessorName            string `json:"server_processor_name,omitempty"`
	ServerClass                    string `json:"server_class,omitempty"`
	ServerTypeIsExperimental       bool   `json:"server_type_is_experimental,omitempty"`
	ServerCount                    int    `json:"server_count,omitempty"`
}

//HardwareConfiguration holds the desired hardware configuration when trying to find available servers for provisioning.
type HardwareConfiguration struct {
	InstanceArrayRAMGbytes          int      `json:"instance_array_ram_gbytes,omitempty"`
	InstanceArrayProcessorCount     int      `json:"instance_array_processor_count,omitempty"`
	InstanceArrayProcessorCoreMHZ   int      `json:"instance_array_processor_core_mhz,omitempty"`
	InstanceArrayProcessorCoreCount int      `json:"instance_array_processor_core_count,omitempty"`
	InstanceArrayDiskCount          int      `json:"instance_array_disk_count,omitempty"`
	InstanceArrayDiskSizeMBytes     int      `json:"instance_array_disk_size_mbytes,omitempty"`
	InstanceArrayTotalMhz           int      `json:"instance_array_total_mhz,omitempty"`
	InstanceArrayDiskTypes          []string `json:"instance_array_disk_types,omitempty"`
	InstanceArrayInstanceCount      int      `json:"instance_array_instance_count,omitempty"`
}

//ServerTypeGet retrieves a server type by id
func (c *Client) ServerTypeGet(serverTypeID ID) (*ServerType, error) {
	var createdObject ServerType

	if err := checkID(serverTypeID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(
		&createdObject,
		"server_type_get",
		serverTypeID)

	if err != nil {
		return nil, err
	}

	return &createdObject, nil
}

//ServerTypesMatches matches available servers with a certain Instance's configuration, using the properties specified in the objHardwareConfiguration object, and returns the number of compatible servers for each server_type_id.
func (c *Client) ServerTypesMatches(infrastructureID ID, hardwareConfiguration HardwareConfiguration, instanceArrayID *int, bAllowServerSwap bool) (*map[string]ServerType, error) {

	if err := checkID(infrastructureID); err != nil {
		return nil, err
	}

	res, err := c.rpcClient.Call(
		"server_type_matches",
		infrastructureID)

	if err != nil {
		return nil, err
	}

	_, ok := res.Result.([]interface{})
	if ok {
		var m = map[string]ServerType{}
		return &m, nil
	}

	var createdObject map[string]ServerType

	err2 := res.GetObject(&createdObject)
	if err2 != nil {
		return nil, err2
	}

	return &createdObject, nil
}

//ServerTypesMatchHardwareConfiguration Retrieves a list of server types that match the provided hardware configuration. The function does not check for availability, only compatibility, so physical servers associated with the returned server types might be unavailable.
func (c *Client) ServerTypesMatchHardwareConfiguration(datacenterName string, hardwareConfiguration HardwareConfiguration) (*map[int]ServerType, error) {

	res, err := c.rpcClient.Call(
		"server_types_match_hardware_configuration",
		datacenterName,
		hardwareConfiguration)

	if err != nil {
		return nil, err
	}

	_, ok := res.Result.([]interface{})
	if ok {
		var m = map[int]ServerType{}
		return &m, nil
	}

	var createdObject map[int]ServerType

	err2 := res.GetObject(&createdObject)
	if err2 != nil {
		return nil, err2
	}

	return &createdObject, nil
}

//ServerTypeDatacenter retrieves all the server type IDs for servers found in a specified Datacenter
func (c *Client) ServerTypeDatacenter(datacenterName string) (*[]int, error) {

	res, err := c.rpcClient.Call(
		"server_types_datacenter",
		datacenterName)

	if err != nil {
		return nil, err
	}

	_, ok := res.Result.([]interface{})
	if ok {
		var m = []int{}
		return &m, nil
	}

	var createdObject []int

	err2 := res.GetObject(&createdObject)
	if err2 != nil {
		return nil, err2
	}

	return &createdObject, nil
}

//ServerTypes retrieves all ServerType objects from the database.
func (c *Client) ServerTypes(datacenterName string, bOnlyAvailable bool) (*map[int]ServerType, error) {

	res, err := c.rpcClient.Call(
		"server_types",
		datacenterName,
		bOnlyAvailable)

	if err != nil {
		return nil, err
	}

	_, ok := res.Result.([]interface{})
	if ok {
		var m = map[int]ServerType{}
		return &m, nil
	}

	var createdObject map[int]ServerType

	err2 := res.GetObject(&createdObject)
	if err2 != nil {
		return nil, err2
	}

	return &createdObject, nil
}
