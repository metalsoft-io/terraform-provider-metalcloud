package metalcloud

//type SchemaField string
//type ExtensionStatus string
//type ExtensionInputType string

const (
	//Schema Elements
	InfrastructureId = /*SchemaField*/ "infrastructure_id"

	ExtensionId    = /*SchemaField*/ "extension_id"
	ExtensionLabel = /*SchemaField*/ "extension_label"

	ExtensionInstanceId     = /*SchemaField*/ "extension_instance_id"
	ExtensionInstanceLabel  = /*SchemaField*/ "extension_instance_label"
	ExtensionInstanceInput  = /*SchemaField*/ "extension_instance_input"
	ExtensionInstanceOutput = /*SchemaField*/ "extension_instance_output"

	//Extension Statuses
	ExtensionStatus_Draft    = /*ExtensionStatus*/ "draft"
	ExtensionStatus_Active   = /*ExtensionStatus*/ "active"
	ExtensionStatus_Archived = /*ExtensionStatus*/ "archived"
	//Extension Input Types
	ExtensionInputType_String     = /*ExtensionInputType*/ "ExtensionInputString"
	ExtensionInputType_Integer    = /*ExtensionInputType*/ "ExtensionInputInteger"
	ExtensionInputType_ServerType = /*ExtensionInputType*/ "ExtensionInputServerType"
	ExtensionInputType_OSTemplate = /*ExtensionInputType*/ "ExtensionInputOsTemplate"
	ExtensionInputType_Boolean    = /*ExtensionInputType*/ "ExtensionInputBoolean"
)
