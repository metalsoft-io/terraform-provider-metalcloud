package metalcloud

// Common constants
const (
	fieldInfrastructureId = "infrastructure_id"
	fieldCustomVariables  = "custom_variables"
	fieldNetworkId        = "network_id"
	fieldNetworkProfileId = "network_profile_id"
)

// Extension constants
const (
	fieldExtensionId    = "extension_id"
	fieldExtensionLabel = "extension_label"

	fieldExtensionInstanceId     = "extension_instance_id"
	fieldExtensionInstanceLabel  = "extension_instance_label"
	fieldExtensionInstanceInput  = "extension_instance_input"
	fieldExtensionInstanceOutput = "extension_instance_output"

	//Extension Statuses
	extensionStatus_Draft    = /*ExtensionStatus*/ "draft"
	extensionStatus_Active   = /*ExtensionStatus*/ "active"
	extensionStatus_Archived = /*ExtensionStatus*/ "archived"
	//Extension Input Types
	extensionInputType_String     = /*ExtensionInputType*/ "ExtensionInputString"
	extensionInputType_Integer    = /*ExtensionInputType*/ "ExtensionInputInteger"
	extensionInputType_ServerType = /*ExtensionInputType*/ "ExtensionInputServerType"
	extensionInputType_OSTemplate = /*ExtensionInputType*/ "ExtensionInputOsTemplate"
	extensionInputType_Boolean    = /*ExtensionInputType*/ "ExtensionInputBoolean"
)

// VM constants
const (
	fieldVmTypeId    = "vm_type_id"
	fieldVmTypeLabel = "vm_type_label"
	fieldVmCpuCores  = "vm_type_cpu_cores"
	fieldVmRamGbytes = "vm_type_ram_gbytes"

	fieldVmInstanceGroupId              = "vm_instance_group_id"
	fieldVmInstanceGroupLabel           = "vm_instance_group_label"
	fieldVmInstanceGroupInstanceCount   = "vm_instance_group_instance_count"
	fieldVmInstanceGroupDiskSizeGbytes  = "vm_instance_group_disk_size_gbytes"
	fieldVmInstanceGroupTemplateId      = "vm_instance_group_template_id"
	fieldVmInstanceCustomVariables      = "vm_instance_custom_variables"
	fieldVmInstanceGroupInterfaces      = "vm_instance_group_interfaces"
	fieldVmInstanceGroupNetworkProfiles = "vm_instance_group_network_profiles"
	fieldVmInstanceIndex                = "vm_instance_index"
	fieldVmInterfaceIndex               = "vm_interface_index"
)
