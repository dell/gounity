/*
 Copyright © 2019 Dell Inc. or its subsidiaries. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
      http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package gounity

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/dell/gounity/util"

	"github.com/dell/gounity/api"
	"github.com/dell/gounity/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emcCsrfToken = "EMC-CSRF-TOKEN" // #nosec G101
)

var (
	accHeader string
	conHeader string

	debug, _    = strconv.ParseBool(os.Getenv("GOUNITY_DEBUG"))
	showHTTP, _ = strconv.ParseBool(os.Getenv("GOUNITY_SHOWHTTP"))
	errNoLink   = errors.New("error: problem finding link")
)

// UnityClient interface for Unity Client
type UnityClient interface {
	Authenticate(ctx context.Context, configConnect *ConfigConnect) error
	BasicSystemInfo(ctx context.Context, configConnect *ConfigConnect) error
	GetToken() string
	SetToken(token string)
	CreateFilesystem(ctx context.Context, name string, storagepool string, description string, nasServer string, size uint64, tieringPolicy int, hostIOSize int, supportedProtocol int, isThinEnabled bool, isDataReductionEnabled bool) (*types.Filesystem, error)
	CreateNFSShare(ctx context.Context, name string, path string, filesystemID string, nfsShareDefaultAccess NFSShareDefaultAccess) (*types.Filesystem, error)
	CreateNFSShareFromSnapshot(ctx context.Context, name string, path string, snapshotID string, nfsShareDefaultAccess NFSShareDefaultAccess) (*types.NFSShare, error)
	DeleteFilesystem(ctx context.Context, filesystemID string) error
	DeleteNFSShare(ctx context.Context, filesystemID string, nfsShareID string) error
	DeleteNFSShareCreatedFromSnapshot(ctx context.Context, nfsShareID string) error
	ExpandFilesystem(ctx context.Context, filesystemID string, newSize uint64) error
	FindFilesystemByID(ctx context.Context, filesystemID string) (*types.Filesystem, error)
	FindFilesystemByName(ctx context.Context, filesystemName string) (*types.Filesystem, error)
	FindNASServerByID(ctx context.Context, nasServerID string) (*types.NASServer, error)
	FindNFSShareByID(ctx context.Context, nfsShareID string) (*types.NFSShare, error)
	FindNFSShareByName(ctx context.Context, nfsSharename string) (*types.NFSShare, error)
	GetFilesystemIDFromResID(ctx context.Context, filesystemResID string) (string, error)
	ModifyNFSShareCreatedFromSnapshotHostAccess(ctx context.Context, nfsShareID string, hostIDs []string, accessType AccessType) error
	ModifyNFSShareHostAccess(ctx context.Context, filesystemID string, nfsShareID string, hostIDs []string, accessType AccessType) error
	FindHostByName(ctx context.Context, hostName string) (*types.Host, error)
	CreateHost(ctx context.Context, hostName string, tenantID string) (*types.Host, error)
	DeleteHost(ctx context.Context, hostName string) error
	CreateHostIPPort(ctx context.Context, hostID, ip string) (*types.HostIPPort, error)
	FindHostIPPortByID(ctx context.Context, hostIPID string) (*types.HostIPPort, error)
	ListHostInitiators(ctx context.Context) ([]types.HostInitiator, error)
	FindHostInitiatorByName(ctx context.Context, wwnOrIqn string) (*types.HostInitiator, error)
	FindHostInitiatorByID(ctx context.Context, wwnOrIqn string) (*types.HostInitiator, error)
	CreateHostInitiator(ctx context.Context, hostID, wwnOrIqn string, initiatorType types.InitiatorType) (*types.HostInitiator, error)
	ModifyHostInitiator(ctx context.Context, hostID string, initiator *types.HostInitiator) (*types.HostInitiator, error)
	ModifyHostInitiatorByID(ctx context.Context, hostID, initiatorID string) (*types.HostInitiator, error)
	FindHostInitiatorPathByID(ctx context.Context, initiatorPathID string) (*types.HostInitiatorPath, error)
	FindFcPortByID(ctx context.Context, fcPortID string) (*types.FcPort, error)
	FindTenants(ctx context.Context) (*types.TenantInfo, error)
	ListIscsiIPInterfaces(ctx context.Context) ([]types.IPInterfaceEntries, error)
	CreateRealTimeMetricsQuery(ctx context.Context, metricPaths []string, interval int) (*types.MetricQueryCreateResponse, error)
	DeleteRealTimeMetricsQuery(ctx context.Context, queryID int) error
	GetAllRealTimeMetricPaths(ctx context.Context) error
	GetCapacity(ctx context.Context) (*types.SystemCapacityMetricsQueryResult, error)
	GetMetricsCollection(ctx context.Context, queryID int) (*types.MetricQueryResult, error)
	CopySnapshot(ctx context.Context, sourceSnapshotID string, name string) (*types.Snapshot, error)
	CreateSnapshot(ctx context.Context, storageResourceID string, snapshotName string, description string, retentionDuration string) (*types.Snapshot, error)
	CreateSnapshotWithFsAccesType(ctx context.Context, storageResourceID string, snapshotName string, _ string, retentionDuration string, filesystemAccessType FilesystemAccessType) (*types.Snapshot, error)
	DeleteFilesystemAsSnapshot(ctx context.Context, snapshotID string, sourceFs *types.Filesystem) error
	DeleteSnapshot(ctx context.Context, snapshotID string) error
	FindSnapshotByID(ctx context.Context, snapshotID string) (*types.Snapshot, error)
	FindSnapshotByName(ctx context.Context, snapshotName string) (*types.Snapshot, error)
	ListSnapshots(ctx context.Context, startToken int, maxEntries int, sourceVolumeID string, snapshotID string) ([]types.Snapshot, int, error)
	ModifySnapshot(ctx context.Context, snapshotID string, description string, retentionDuration string) error
	ModifySnapshotAutoDeleteParameter(ctx context.Context, snapshotID string) error
	FindStoragePoolByName(ctx context.Context, poolName string) (*types.StoragePool, error)
	FindStoragePoolByID(ctx context.Context, poolID string) (*types.StoragePool, error)
	CreateCloneFromVolume(ctx context.Context, name string, volID string) (*types.Volume, error)
	CreateLun(ctx context.Context, name string, poolID string, description string, size uint64, fastVPTieringPolicy int, hostIOLimitID string, isThinEnabled bool, isDataReductionEnabled bool) (*types.Volume, error)
	CreteLunThinClone(ctx context.Context, name string, snapID string, volID string) (*types.Volume, error)
	DeleteVolume(ctx context.Context, volumeID string) error
	ExpandVolume(ctx context.Context, volumeID string, newSize uint64) error
	ExportVolume(ctx context.Context, volID string, hostID string) error
	FindHostIOLimitByName(ctx context.Context, hostIoPolicyName string) (*types.IoLimitPolicy, error)
	FindVolumeByID(ctx context.Context, volID string) (*types.Volume, error)
	FindVolumeByName(ctx context.Context, volName string) (*types.Volume, error)
	GetMaxVolumeSize(ctx context.Context, systemLimitID string) (*types.MaxVolumSizeInfo, error)
	ListVolumes(ctx context.Context, startToken int, maxEntries int) ([]types.Volume, int, error)
	ModifyVolumeExport(ctx context.Context, volID string, hostIDList []string) error
	RenameVolume(ctx context.Context, newName string, volID string) error
	UnexportVolume(ctx context.Context, volID string) error
	GetAllNFSServers(ctx context.Context) (*types.NFSServersResponse, error)
}

// UnityClientImpl Struct holds the configuration & REST Client.
type UnityClientImpl struct {
	configConnect *ConfigConnect
	api           api.Client
	loginMutex    sync.Mutex
}

// ConfigConnect Struct holds the endpoint & credential info.
type ConfigConnect struct {
	Endpoint string
	Username string
	Password string
	Insecure bool
}

// BasicSystemInfo make a REST API call [/basicSystemInfo/instances] to Unity to check if array is responding.
func (c *UnityClientImpl) BasicSystemInfo(ctx context.Context, configConnect *ConfigConnect) error {
	log := util.GetRunIDLogger(ctx)
	log.Debug("Executing BasicSystemInfo REST client")
	c.configConnect = configConnect
	headers := make(map[string]string, 3)
	headers[api.XEmcRestClient] = "true"
	headers[api.HeaderKeyContentType] = api.HeaderValContentTypeJSON
	resp, err := c.api.DoAndGetResponseBody(ctx, http.MethodGet, api.UnityAPIBasicSysInfoURI, headers, nil)
	if err != nil {
		return fmt.Errorf("Error getting BasicSystemInfo: %v", err)
	}

	if resp != nil {
		log.Debugf("BasicSystemInfo response code: %d", resp.StatusCode)
		if err != nil {
			log.Errorf("Reading BasicSystemInfo response body error:%v", err)
		}

		defer resp.Body.Close()

		switch {
		case resp.StatusCode >= 200 && resp.StatusCode <= 299:
			{
				log.Debug("Getting BasicSystemInfo details successful")
			}
		default:
			return fmt.Errorf("Get BaicSystemInfo error. Response: %v", c.api.ParseJSONError(ctx, resp))
		}

	} else {
		log.Errorf("Getting BasicSystenInfo details faile")
	}
	return nil
}

// Authenticate make a REST API call [/loginSessionInfo] to Unity to get authenticate the given credentials.
// The response contains the EMC-CSRF-TOKEN and the client caches it for further communication.
func (c *UnityClientImpl) Authenticate(ctx context.Context, configConnect *ConfigConnect) error {
	c.loginMutex.Lock()
	defer c.loginMutex.Unlock()
	log := util.GetRunIDLogger(ctx)
	log.Debug("Executing Authenticate REST client")
	c.configConnect = configConnect
	c.api.SetToken("")
	headers := make(map[string]string, 3)
	headers[api.AuthorizationHeader] = "Basic " + basicAuth(configConnect.Username, configConnect.Password)
	headers[api.XEmcRestClient] = "true"
	headers[api.HeaderKeyContentType] = api.HeaderValContentTypeJSON
	resp, err := c.api.DoAndGetResponseBody(ctx, http.MethodGet, api.UnityAPILoginSessionInfoURI, headers, nil)
	if err != nil {
		return fmt.Errorf("authentication error: %v", err)
	}

	if resp != nil {
		log.Debugf("Authentication response code: %d", resp.StatusCode)
		if err != nil {
			log.Errorf("Reading Authentication response body error:%v", err)
		}

		defer resp.Body.Close()

		switch {
		case resp.StatusCode >= 200 && resp.StatusCode <= 299:
			{
				log.Debug("Authentication successful")
			}
		case resp.StatusCode == 401:
			{
				return status.Errorf(codes.Unauthenticated, "Authentication failed. Unable to login to Unity. Verify username and password.")
			}
		default:
			return fmt.Errorf("authenticate error. Response: %v", c.api.ParseJSONError(ctx, resp))
		}

		c.api.SetToken(resp.Header.Get(emcCsrfToken))
	} else {
		log.Errorf("Authenticate error: Nil response received")
	}
	return nil
}

// basicAuth converts the given username & password to Base64 encoded string.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// GetJSONWithRetry method responsible to make the given API call to Unity REST API Server.
// In case if the given EMC-CSRF-TOKEN becomes invalid, retries the same operation after performing authentication.
func (c *UnityClientImpl) executeWithRetryAuthenticate(ctx context.Context, method, uri string, body, resp interface{}) error {
	log := util.GetRunIDLogger(ctx)
	headers := make(map[string]string, 2)
	headers[api.HeaderKeyAccept] = accHeader
	headers[api.HeaderKeyContentType] = conHeader
	headers[api.XEmcRestClient] = "true"
	log.Debug("Invoking REST API server info Method: ", method, ", URI: ", uri)
	err := c.api.DoWithHeaders(ctx, method, uri, headers, body, resp)
	if err == nil {
		log.Debug("Execution successful on Method: ", method, ", URI: ", uri)
		return nil
	}
	// check if we need to authenticate
	if e, ok := err.(*types.Error); ok {
		log.Debugf("Error in response. Method:%s URI:%s Error: %v JSON Error: %+v", method, uri, err, e)
		if e.ErrorContent.HTTPStatusCode == 401 {
			log.Debug("need to re-authenticate")
			// Authenticate then try again
			if err := c.Authenticate(ctx, c.configConnect); err != nil {
				return fmt.Errorf("authentication failure due to: %v", err)
			}
			log.Debug("Authentication success")
			return c.api.DoWithHeaders(ctx, method, uri, headers, body, resp)
		}
	} else {
		log.Error("Error is not a type of \"*types.Error\". Error:", err)
	}
	log.WithError(err).Error("failed to invoke Unity REST API server")

	return err
}

// SetToken function sets token
func (c *UnityClientImpl) SetToken(token string) {
	c.api.SetToken(token)
}

// GetToken function gets token
func (c *UnityClientImpl) GetToken() string {
	return c.api.GetToken()
}

// NewClient initialize the new REST Client with default options.
func NewClient(ctx context.Context) (UnityClient, error) {
	insecure, _ := strconv.ParseBool(os.Getenv("GOUNITY_INSECURE"))
	return NewClientWithArgs(ctx, os.Getenv("GOUNITY_ENDPOINT"), insecure)
}

// NewClientWithArgs initialize the new REST Client with the given arguments.
func NewClientWithArgs(ctx context.Context, endpoint string, insecure bool) (UnityClient, error) {
	log := util.GetRunIDLogger(ctx)
	if util.ShowHTTP {
		util.Debug = true
	}

	fields := map[string]interface{}{
		"endpoint": endpoint,
		"insecure": insecure,
		"debug":    util.Debug,
		"showHTTP": util.ShowHTTP,
	}

	log.WithFields(fields).Debug("unity client init")

	if endpoint == "" {
		log.WithFields(fields).Error("endpoint is required")
		return nil, withFields(fields, "endpoint is required")
	}

	opts := api.ClientOptions{
		Insecure: insecure,
		ShowHTTP: util.ShowHTTP,
	}

	ac, err := api.New(ctx, endpoint, opts, util.Debug)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP client %v", err)
	}

	client := &UnityClientImpl{
		api:           ac,
		configConnect: &ConfigConnect{},
	}
	conHeader = api.HeaderValContentTypeJSON
	return client, nil
}

func withFields(fields map[string]interface{}, message string) error {
	return withFieldsE(fields, message, nil)
}

func withFieldsE(fields map[string]interface{}, message string, inner error) error {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if inner != nil {
		fields["inner"] = inner
	}

	x := 0
	l := len(fields)

	var b bytes.Buffer
	for k, v := range fields {
		if x < l-1 {
			b.WriteString(fmt.Sprintf("%s=%v,", k, v))
		} else {
			b.WriteString(fmt.Sprintf("%s=%v", k, v))
		}
		x = x + 1
	}

	return fmt.Errorf("%s %s", message, b.String())
}

func (c *UnityClientImpl) getAPI() api.Client {
	return c.api
}
