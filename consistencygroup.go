/*
Copyright (c) 2021 Dell Corporation
All Rights Reserved
*/

package gounity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"github.com/dell/gounity/util"
)

//Constants
const (
	ConsistencyGroupNameMaxLength             = 95
	ConsistencyGroupNotFoundErrorCode = "0x7d13005"
)


//ErrorConsistencyGroupNotFound stores ConsistencyGroup not found error
var ErrorConsistencyGroupNotFound = errors.New("Unable to find ConsistencyGroup")

//ConsistencyGroup structure
type ConsistencyGroup struct {
	client *Client 
}

//NewConsistencyGroup function returns ConsistencyGroup
func NewConsistencyGroup(client *Client) *ConsistencyGroup {
	return &ConsistencyGroup{client}
}

// GetConsistencyGroup query returns ConsistencyGroup by id
func (c *ConsistencyGroup) GetConsistencyGroup(ctx context.Context, id string) (*types.ConsistencyGroup, error) {
	cgResp  := &types.ConsistencyGroup{}
	log := util.GetRunIDLogger(ctx)

	err := c.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceWithFieldsURI, api.StorageResourceAction, id, ConsistencyGroupDisplayFields), nil, cgResp)

	if err != nil {
		if strings.Contains(err.Error(), ConsistencyGroupNotFoundErrorCode) {
			log.Debugf("Unable to get ConsistencyGroup by Id %s Error: %v", id, err)
			return nil, ErrorConsistencyGroupNotFound
		}
		return nil, err
	}
	
	return cgResp, nil
}

//GetConsistencyGroupByName - Find the GetConsistencyGroup by it's name. If the GetConsistencyGroup is not found, an error will be returned.
func (c *ConsistencyGroup) GetConsistencyGroupByName(ctx context.Context, cgName string) (*types.ConsistencyGroup, error) {
	if len(cgName) == 0 {
		return nil, fmt.Errorf("ConsistencyGroup Name shouldn't be empty")
	}
	cgResp  := &types.ConsistencyGroup{}
	err := c.client.executeWithRetryAuthenticate(ctx, http.MethodGet, fmt.Sprintf(api.UnityAPIGetResourceByNameWithFieldsURI, api.StorageResourceAction, cgName, ConsistencyGroupDisplayFields), nil, cgResp)
	if err != nil {
		return nil, fmt.Errorf("unable to find ConsistencyGroup by name %s", cgName)
	}

	return cgResp, nil
}

// CreateLun API create a ConsistencyGroup with the given arguments.
// Pre-validations: 1. Name is not empty.
//                  2. Length of the ConsistencyGroup name should not exceed ConsistencyGroupNameMaxLength characters
func (c *ConsistencyGroup) CreateConsistencyGroup(ctx context.Context, createParams *types.ConsistencyGroupCreate) (*types.ConsistencyGroup, error) {

	if createParams.Name == "" {
		return nil, errors.New("ConsistencyGroup name should not be empty")
	}

	if len(createParams.Name) > ConsistencyGroupNameMaxLength {
		return nil, fmt.Errorf("ConsistencyGroup name %s should not exceed %d characters", createParams.Name, ConsistencyGroupNameMaxLength)
	}

	cgResp := &types.ConsistencyGroup{}

	err := c.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityAPIStorageResourceActionURI, api.CreateCGAction), createParams, cgResp)
	if err != nil {
		return nil, err
	}

	return cgResp, nil
}

//DeleteConsistencyGroup - Delete ConsistencyGroup by ID. 
func (c *ConsistencyGroup) DeleteConsistencyGroup(ctx context.Context, cgID string) error {
	if len(cgID) == 0 {
		return errors.New("ConsistencyGroup Id cannot be empty")
	}

	_, err := c.GetConsistencyGroup(ctx, cgID)

	if err != nil {
		return err
	}

	deleteErr := c.client.executeWithRetryAuthenticate(ctx, http.MethodDelete, fmt.Sprintf(api.UnityAPIGetResourceURI, api.StorageResourceAction, cgID), nil, nil)

	if deleteErr != nil {
		return fmt.Errorf("delete ConsistencyGroup %s Failed. Error: %v", cgID, deleteErr)
	}

	return nil
}

//ModifyConsistencyGroup - Modify ConsistencyGroup by ID. 
func (c *ConsistencyGroup) ModifyConsistencyGroup(ctx context.Context, cgID string, modifyParams *types.ConsistencyGroupModify) error {

	err := c.client.executeWithRetryAuthenticate(ctx, http.MethodPost, fmt.Sprintf(api.UnityModifyCGURI, cgID), modifyParams, nil)
	if err != nil {
		return err
	}

	return nil
}

