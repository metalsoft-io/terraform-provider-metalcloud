package metalcloud

import (
	"fmt"
)

//Infrastructure - the main infrastructure object
type Infrastructure struct {
	InfrastructureID               int                     `json:"infrastructure_id,omitempty"`
	InfrastructureLabel            string                  `json:"infrastructure_label"`
	DatacenterName                 string                  `json:"datacenter_name"`
	InfrastructureSubdomain        string                  `json:"infrastructure_subdomain,omitempty"`
	UserIDowner                    int                     `json:"user_id_owner,omitempty"`
	UserEmailOwner                 string                  `json:"user_email_owner,omitempty"`
	InfrastructureTouchUnixtime    string                  `json:"infrastructure_touch_unixtime,omitempty"`
	InfrastructureServiceStatus    string                  `json:"infrastructure_service_status,omitempty"`
	InfrastructureCreatedTimestamp string                  `json:"infrastructure_created_timestamp,omitempty"`
	InfrastructureUpdatedTimestamp string                  `json:"infrastructure_updated_timestamp,omitempty"`
	InfrastructureChangeID         int                     `json:"infrastructure_change_id,omitempty"`
	InfrastructureDeployID         int                     `json:"infrastructure_deploy_id,omitempty"`
	InfrastructureDesignIsLocked   bool                    `json:"infrastructure_design_is_locked,omitempty"`
	InfrastructureOperation        InfrastructureOperation `json:"infrastructure_operation,omitempty"`
}

//InfrastructureOperation - object with alternations to be applied
type InfrastructureOperation struct {
	InfrastructureID               int    `json:"infrastructure_id,omitempty"`
	InfrastructureLabel            string `json:"infrastructure_label"`
	DatacenterName                 string `json:"datacenter_name"`
	InfrastructureDeployStatus     string `json:"infrastructure_deploy_status,omitempty"`
	InfrastructureDeployType       string `json:"infrastructure_deploy_type,omitempty"`
	InfrastructureSubdomain        string `json:"infrastructure_subdomain,omitempty"`
	UserIDOwner                    int    `json:"user_id_owner,omitempty"`
	InfrastructureUpdatedTimestamp string `json:"infrastructure_updated_timestamp,omitempty"`
	InfrastructureChangeID         int    `json:"infrastructure_change_id,omitempty"`
	InfrastructureDeployID         int    `json:"infrastructure_deploy_id,omitempty"`
}

//ShutdownOptions controls how the deploy engine handles running instances
type ShutdownOptions struct {
	HardShutdownAfterTimeout   bool `json:"hard_shutdown_after_timeout"`
	AttemptSoftShutdown        bool `json:"attempt_soft_shutdown"`
	SoftShutdownTimeoutSeconds int  `json:"soft_shutdown_timeout_seconds"`
}

//InfrastructureCreate creates an infrastructure
func (c *Client) InfrastructureCreate(infrastructure Infrastructure) (*Infrastructure, error) {
	var createdObject Infrastructure

	err := c.rpcClient.CallFor(
		&createdObject,
		"infrastructure_create",
		c.user,
		infrastructure)

	if err != nil {
		return nil, err
	}

	return &createdObject, nil
}

//InfrastructureEdit alters an infrastructure
func (c *Client) InfrastructureEdit(infrastructureID ID, infrastructureOperation InfrastructureOperation) (*Infrastructure, error) {
	var createdObject Infrastructure

	if err := checkID(infrastructureID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(
		&createdObject,
		"infrastructure_edit",
		infrastructureID,
		infrastructureOperation)

	if err != nil {
		return nil, err
	}

	return &createdObject, nil
}

//InfrastructureDelete deletes an infrastructure and all associated elements. Requires deploy
func (c *Client) InfrastructureDelete(infrastructureID ID) error {

	if err := checkID(infrastructureID); err != nil {
		return err
	}

	resp, err := c.rpcClient.Call("infrastructure_delete", infrastructureID)

	if err != nil {
		return err
	}

	if resp.Error != nil {
		return fmt.Errorf(resp.Error.Message)
	}

	return nil
}

//InfrastructureOperationCancel reverts (undos) alterations done before deploy
func (c *Client) InfrastructureOperationCancel(infrastructureID ID) error {

	if err := checkID(infrastructureID); err != nil {
		return err
	}

	resp, err := c.rpcClient.Call(
		"infrastructure_operation_cancel",
		infrastructureID)

	if err != nil {
		return err
	}

	if resp.Error != nil {
		return fmt.Errorf(resp.Error.Message)
	}

	return nil
}

//InfrastructureDeploy initiates a deploy operation that will apply all registered changes for the respective infrastructure
func (c *Client) InfrastructureDeploy(infrastructureID ID, shutdownOptions ShutdownOptions, allowDataLoss bool, skipAnsible bool) error {

	if err := checkID(infrastructureID); err != nil {
		return err
	}

	resp, err := c.rpcClient.Call(
		"infrastructure_deploy",
		infrastructureID,
		shutdownOptions,
		nil,
		allowDataLoss,
		skipAnsible,
	)

	if err != nil {
		return err
	}

	if resp.Error != nil {
		return fmt.Errorf(resp.Error.Message)
	}

	return nil
}

//Infrastructures returns a list of infrastructures
func (c *Client) Infrastructures() (*map[string]Infrastructure, error) {

	res, err := c.rpcClient.Call(
		"infrastructures",
		c.user)

	if err != nil {
		return nil, err
	}

	_, ok := res.Result.([]interface{})
	if ok {
		var m = map[string]Infrastructure{}
		return &m, nil
	}

	var createdObject map[string]Infrastructure

	err2 := res.GetObject(&createdObject)
	if err2 != nil {
		return nil, err2
	}

	return &createdObject, nil
}

//InfrastructureGet returns a specific infrastructure
func (c *Client) InfrastructureGet(infrastructureID ID) (*Infrastructure, error) {
	var infrastructure Infrastructure

	if err := checkID(infrastructureID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(&infrastructure, "infrastructure_get", infrastructureID)

	if err != nil {
		return nil, err
	}

	return &infrastructure, nil
}

//InfrastructureUserLimits returns user metadata
func (c *Client) InfrastructureUserLimits(infrastructureID ID) (*map[string]interface{}, error) {
	var userLimits map[string]interface{}

	if err := checkID(infrastructureID); err != nil {
		return nil, err
	}

	err := c.rpcClient.CallFor(&userLimits, "infrastructure_user_limits", infrastructureID)

	if err != nil {
		return nil, err
	}

	return &userLimits, nil
}
