package api

import "github.com/dell/gounity/types"

//InitiatorType constants
const (
	FCInitiatorType     types.InitiatorType = "1"
	ISCSCIInitiatorType types.InitiatorType = "2"
)

//constants and URIs
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

	//UnityAPIGetResourceURI GETS RESOURCE BY RESOURCE ID {1}=type of resource, {2}=resource id
	UnityAPIGetResourceURI = UnityAPIInstancesURI + "/%s/%s"

	//UnityAPIGetResourceByNameURI GETS RESOURCE BY RESOURCE NAME {1}=type of resource, {2}=name of the resource
	UnityAPIGetResourceByNameURI = UnityAPIInstancesURI + "/%s/name:%s"

	//UnityAPIGetResourceWithFieldsURI GETS RESOURCE BY RESOURCE ID {1}=type of resource, {2}=resource id, {3}=fields
	UnityAPIGetResourceWithFieldsURI = UnityAPIGetResourceURI + "?fields=%s"

	//UnityAPIGetResourceByNameWithFieldsURI GETS RESOURCE BY RESOURCE NAME {1}=type of resource, {2}=name of the resource, {3}=fields
	UnityAPIGetResourceByNameWithFieldsURI = UnityAPIGetResourceByNameURI + "?fields=%s"

	//UnityAPILoginSessionInfoURI LOGINS resource URIs
	UnityAPILoginSessionInfoURI = unityAPITypes + "/loginSessionInfo"

	//UnityAPIBasicSysInfoURI gets BasicSystemInfo URI
	UnityAPIBasicSysInfoURI = unityAPITypes + "/basicSystemInfo/instances"

	//UnityAPIStorageResourceInstanceActionURI gets StorageResource instance Action URI
	UnityAPIStorageResourceInstanceActionURI = UnityAPIInstancesURI + "/storageResource/%s/action"

	//UnityAPIModifyLunURI Modify Volume URIs
	UnityAPIModifyLunURI = UnityAPIStorageResourceInstanceActionURI + "/modifyLun"

	//UnityAPICreateLunThinCloneURI Create LUN Thin Clone
	UnityAPICreateLunThinCloneURI = UnityAPIStorageResourceInstanceActionURI + "/createLunThinClone"

	//UnityAPIModifyStorageResourceURI gets StorageResource resource URIs
	UnityAPIModifyStorageResourceURI = UnityAPIInstancesURI + "/storageResource/%s"
	//UnityAPIStorageResourceActionURI gets StorageResource Action resource URI
	UnityAPIStorageResourceActionURI = unityAPITypes + "/storageResource/action/%s"

	//UnityModifyLunURI Modify Lun URIs
	UnityModifyLunURI = UnityAPIModifyStorageResourceURI + "/action/modifyLun"

	//UnityModifyCGURI Modify Consistency Group URIs
	UnityModifyCGURI = UnityAPIModifyStorageResourceURI + "/action/modifyConsistencyGroup"

	//UnityModifyFilesystemURI Modify Filesystem URIs
	UnityModifyFilesystemURI = UnityAPIModifyStorageResourceURI + "/action/modifyFilesystem"

	//UnityModifyNFSShareURI Modify NFS Share URIs
	UnityModifyNFSShareURI = UnityAPIGetResourceURI + "/action/modify"

	//UnityModifySnapshotURI Snapshot Action resource URIs
	UnityModifySnapshotURI = UnityAPIGetResourceURI + "/action/modify"

	//UnityCopySnapshotURI does Snapshot Copy Action
	UnityCopySnapshotURI = UnityAPIGetResourceURI + "/action/copy"

	//UnityListHostInitiatorsURI gets Host Initiator URIs
	UnityListHostInitiatorsURI = unityAPITypes + "/hostInitiator/instances?fields="
	UnityModifyHostInitiators  = unityRootAPI + "/instances/hostInitiator/%s/action/modify"

	//UnityInstancesFilter does Unity Instance Filter
	UnityInstancesFilter = UnityAPIInstanceTypeResources + "?filter=%s"

	//UnityInstancesFilter does Unity Instance Filter with fields
	UnityInstancesFilterWithFields = UnityAPIInstanceTypeResourcesWithFields + "&filter=%s"

	UnityMetric              = "metric"
	UnityMetricQueryResult   = "metricQueryResult"
	UnityMetricRealTimeQuery = "metricRealTimeQuery"

	//Action types for URL's

	LunAction                = "lun"
	CreateLunAction          = "createLun"
	FileSystemAction         = "filesystem"
	CreateFSAction           = "createFilesystem"
	NfsShareAction           = "nfsShare"
	StorageResourceAction    = "storageResource"
	HostAction               = "host"
	IPInterface              = "ipInterface"
	SnapAction               = "snap"
	PoolAction               = "pool"
	IOLimitPolicy            = "ioLimitPolicy"
	LicenseAction            = "license"
	HostInitiatorPathAction  = "hostInitiatorPath"
	HostInitiatorAction      = "hostInitiator"
	HostIPPortAction         = "hostIPPort"
	NasServerAction          = "nasServer"
	TenantAction             = "tenant"
	RemoteSystemAction       = "remoteSystem"
	ReplicationSessionAction = "replicationSession"
	CreateCGAction           = "createConsistencyGroup"
	ModifyCGAction           = "modifyConsistencyGroup/%s"
)
