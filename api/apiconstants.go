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

package api

import "github.com/dell/gounity/types"

// InitiatorType constants
const (
	FCInitiatorType     types.InitiatorType = "1"
	ISCSCIInitiatorType types.InitiatorType = "2"
)

// constants and URIs
const (
	MaxResourceNameLength = 63

	AuthorizationHeader = "Authorization"
	XEmcRestClient      = "X-EMC-REST-CLIENT"
	// Base resource URIs
	unityRootAPI  = "/api"
	unityAPITypes = unityRootAPI + "/types"

	UnityAPIInstanceTypeResources = unityAPITypes + "/%s" + "/instances"

	UnityAPIGetTenantURI = UnityAPIInstanceTypeResources + "?&compact=true&fields=%s"

	UnityAPIInstanceTypeResourcesWithFields = UnityAPIInstanceTypeResources + "?fields=%s"

	UnityAPIInstancesURI = unityRootAPI + "/instances"

	// UnityAPIGetResourceURI GETS RESOURCE BY RESOURCE ID {1}=type of resource, {2}=resource id
	UnityAPIGetResourceURI = UnityAPIInstancesURI + "/%s/%s"

	// UnityAPIGetResourceByNameURI GETS RESOURCE BY RESOURCE NAME {1}=type of resource, {2}=name of the resource
	UnityAPIGetResourceByNameURI = UnityAPIInstancesURI + "/%s/name:%s"

	// UnityAPIGetResourceWithFieldsURI GETS RESOURCE BY RESOURCE ID {1}=type of resource, {2}=resource id, {3}=fields
	UnityAPIGetResourceWithFieldsURI = UnityAPIGetResourceURI + "?fields=%s"

	// UnityAPIGetResourceByNameWithFieldsURI GETS RESOURCE BY RESOURCE NAME {1}=type of resource, {2}=name of the resource, {3}=fields
	UnityAPIGetResourceByNameWithFieldsURI = UnityAPIGetResourceByNameURI + "?fields=%s"

	// UnityAPILoginSessionInfoURI LOGINS resource URIs
	UnityAPILoginSessionInfoURI = unityAPITypes + "/loginSessionInfo"

	// UnityAPIBasicSysInfoURI gets BasicSystemInfo URI
	UnityAPIBasicSysInfoURI = unityAPITypes + "/basicSystemInfo/instances"

	// UnityAPIStorageResourceInstanceActionURI gets StorageResource instance Action URI
	UnityAPIStorageResourceInstanceActionURI = UnityAPIInstancesURI + "/storageResource/%s/action"

	// UnityAPIModifyLunURI Modify Volume URIs
	UnityAPIModifyLunURI = UnityAPIStorageResourceInstanceActionURI + "/modifyLun"

	// UnityAPICreateLunThinCloneURI Create LUN Thin Clone
	UnityAPICreateLunThinCloneURI = UnityAPIStorageResourceInstanceActionURI + "/createLunThinClone"

	// UnityAPIModifyStorageResourceURI gets StorageResource resource URIs
	UnityAPIModifyStorageResourceURI = UnityAPIInstancesURI + "/storageResource/%s"
	// UnityAPIStorageResourceActionURI gets StorageResource Action resource URI
	UnityAPIStorageResourceActionURI = unityAPITypes + "/storageResource/action/%s"

	// UnityModifyLunURI Modify Lun URIs
	UnityModifyLunURI = UnityAPIModifyStorageResourceURI + "/action/modifyLun"

	// UnityModifyFilesystemURI Modify Filesystem URIs
	UnityModifyFilesystemURI = UnityAPIModifyStorageResourceURI + "/action/modifyFilesystem"

	// UnityModifyNFSShareURI Modify NFS Share URIs
	UnityModifyNFSShareURI = UnityAPIGetResourceURI + "/action/modify"

	// UnityModifySnapshotURI Snapshot Action resource URIs
	UnityModifySnapshotURI = UnityAPIGetResourceURI + "/action/modify"

	// UnityCopySnapshotURI does Snapshot Copy Action
	UnityCopySnapshotURI = UnityAPIGetResourceURI + "/action/copy"

	// UnityAPIGetMaxVolumeSize gets the maximum volume size of an array {1}=unique identifier of the systemLimit instance, {2}=fields
	UnityAPIGetMaxVolumeSize = UnityAPIInstancesURI + "/systemLimit/%s?fields=%s"

	// UnityListHostInitiatorsURI gets Host Initiator URIs
	UnityListHostInitiatorsURI = unityAPITypes + "/hostInitiator/instances?fields="
	UnityModifyHostInitiators  = unityRootAPI + "/instances/hostInitiator/%s/action/modify"

	// UnityInstancesFilter does Unity Instance Filter
	UnityInstancesFilter = UnityAPIInstanceTypeResources + "?filter=%s"

	UnityMetric              = "metric"
	UnityMetricQueryResult   = "metricQueryResult"
	UnityMetricRealTimeQuery = "metricRealTimeQuery"

	// UnitySystemCapacity is used to get capacity metrics for Unity XT
	UnitySystemCapacity = "systemCapacity"

	// Action types for URL's

	LunAction               = "lun"
	CreateLunAction         = "createLun"
	FileSystemAction        = "filesystem"
	CreateFSAction          = "createFilesystem"
	NfsShareAction          = "nfsShare"
	StorageResourceAction   = "storageResource"
	HostAction              = "host"
	IPInterface             = "ipInterface"
	SnapAction              = "snap"
	PoolAction              = "pool"
	IOLimitPolicy           = "ioLimitPolicy"
	LicenseAction           = "license"
	HostInitiatorPathAction = "hostInitiatorPath"
	HostInitiatorAction     = "hostInitiator"
	HostIPPortAction        = "hostIPPort"
	NasServerAction         = "nasServer"
	TenantAction            = "tenant"
	UnityNFSServer          = "nfsServer"
	UnityNFS3AndNFS4Enabled = "nfsv3Enabled,nfsv4Enabled"
)
