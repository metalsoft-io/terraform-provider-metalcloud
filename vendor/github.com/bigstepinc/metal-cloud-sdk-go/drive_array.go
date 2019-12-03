package metalcloud

import "fmt"

//DriveArray represents a collection of identical drives
type DriveArray struct {
	DriveArrayID                      int    `json:"drive_array_id,omitempty"`
	DriveArrayLabel                   string `json:"drive_array_label,omitempty"`
	VolumeTemplateID                  int    `json:"volume_template_id,omitempty"`
	DriveArrayStorageType             string `json:"drive_array_storage_type,omitempty"`
	DriveSizeMBytesDefault            int    `json:"drive_size_mbytes_default,omitempty"`
	InstanceArrayID                   int    `json:"instance_array_id,omitempty"`
	InfrastructureID                  int    `json:"infrastructure_id,omitempty"`
	DriveArrayServiceStatus           string `json:"drive_array_service_status,omitempty"`
	DriveArrayCount                   int    `json:"drive_array_count,omitempty"`
	DriveArrayExpandWithInstanceArray bool   `json:"drive_array_expand_with_instance_array,omitempty"`

	DriveArrayOperation *DriveArrayOperation `json:"drive_array_operation,omitempty"`
}

//DriveArrayOperation defines changes to be applied to a DriveArray
type DriveArrayOperation struct {
	DriveArrayID                      int    `json:"drive_array_id,omitempty"`
	DriveArrayLabel                   string `json:"drive_array_label,omitempty"`
	VolumeTemplateID                  int    `json:"volume_template_id,omitempty"`
	DriveArrayStorageType             string `json:"drive_array_storage_type,omitempty"`
	DriveSizeMBytesDefault            int    `json:"drive_size_mbytes_default,omitempty"`
	InstanceArrayID                   int    `json:"instance_array_id,omitempty"`
	InfrastructureID                  int    `json:"infrastructure_id,omitempty"`
	DriveArrayCount                   int    `json:"drive_array_count,omitempty"`
	DriveArrayExpandWithInstanceArray bool   `json:"drive_array_expand_with_instance_array,omitempty"`

	DriveArrayDeployType   string `json:"drive_array_deploy_type,omitempty"`
	DriveArrayDeployStatus string `json:"drive_array_deploy_status,omitempty"`
	DriveArrayChangeID     int    `json:"drive_array_change_id,omitempty"`
}

//DriveArrays retrieves the list of drives arrays of an infrastructure
func (c *Client) DriveArrays(infrastructureID ID) (*map[string]DriveArray, error) {

	if err := checkID(infrastructureID); err != nil {
		return nil, err
	}

	res, err := c.rpcClient.Call(
		"drive_arrays",
		infrastructureID)

	if err != nil {
		return nil, err
	}

	_, ok := res.Result.([]interface{})
	if ok {
		var m = map[string]DriveArray{}
		return &m, nil
	}

	var createdObject map[string]DriveArray

	err2 := res.GetObject(&createdObject)
	if err2 != nil {
		return nil, err2
	}

	return &createdObject, nil
}

//DriveArrayGet retrieves a DriveArray object with specified ids
func (c *Client) DriveArrayGet(driveArrayID ID) (*DriveArray, error) {

	var createdObject DriveArray

	if err := checkID(driveArrayID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(
		&createdObject,
		"drive_array_get",
		driveArrayID)

	if err != nil {

		return nil, err
	}

	return &createdObject, nil
}

//DriveArrayCreate creates a drive array. Requires deploy.
func (c *Client) DriveArrayCreate(infrastructureID ID, driveArray DriveArray) (*DriveArray, error) {
	var createdObject DriveArray

	if err := checkID(infrastructureID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(
		&createdObject,
		"drive_array_create",
		infrastructureID,
		driveArray)

	if err != nil {

		return nil, err
	}

	return &createdObject, nil
}

//DriveArrayEdit alters a deployed drive array. Requires deploy.
func (c *Client) DriveArrayEdit(driveArrayID ID, driveArrayOperation DriveArrayOperation) (*DriveArray, error) {
	var createdObject DriveArray

	if err := checkID(driveArrayID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(
		&createdObject,
		"drive_array_edit",
		driveArrayID,
		driveArrayOperation)

	if err != nil {

		return nil, err
	}

	return &createdObject, nil
}

//DriveArrayDelete deletes a Drive Array with specified id
func (c *Client) DriveArrayDelete(driveArrayID ID) error {

	if err := checkID(driveArrayID); err != nil {
		return err
	}

	resp, err := c.rpcClient.Call(
		"drive_array_delete",
		driveArrayID)

	if err != nil {

		return err
	}

	if resp.Error != nil {
		return fmt.Errorf(resp.Error.Message)
	}

	return nil
}
