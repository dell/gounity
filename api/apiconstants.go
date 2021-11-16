package api

import "github.com/dell/gounity/types"

const (
	FCInitiatorType     types.InitiatorType = "1"
	ISCSCIInitiatorType types.InitiatorType = "2"
)

const (
	MaxResourceNameLength = 63

	AuthorizationHeader = "Authorization"
	XEmcRestClient      = "X-EMC-REST-CLIENT"
	// Base resource URIs
	unityRootApi  = "/api"
	unityApiTypes = unityRootApi + "/types"

	UnityApiInstanceTypeResources = unityApiTypes + "/%s" + "/instances"

	UnityApiGetTenantUri = UnityApiInstanceTypeResources + "?&compact=true&fields=%s"

	UnityApiInstanceTypeResourcesWithFields = UnityApiInstanceTypeResources + "?fields=%s"

	UnityApiInstancesUri = unityRootApi + "/instances"

	// GET RESOURCE BY RESOURCE ID {1}=type of resource, {2}=resource id
	UnityApiGetResourceUri = UnityApiInstancesUri + "/%s/%s"

	// GET RESOURCE BY RESOURCE NAME {1}=type of resource, {2}=name of the resource
	UnityApiGetResourceByNameUri = UnityApiInstancesUri + "/%s/name:%s"

	// GET RESOURCE BY RESOURCE ID {1}=type of resource, {2}=resource id, {3}=fields
	UnityApiGetResourceWithFieldsUri = UnityApiGetResourceUri + "?fields=%s"

	// GET RESOURCE BY RESOURCE NAME {1}=type of resource, {2}=name of the resource, {3}=fields
	UnityApiGetResourceByNameWithFieldsUri = UnityApiGetResourceByNameUri + "?fields=%s"

	// LOGIN resource URIs
	UnityApiLoginSessionInfoUri = unityApiTypes + "/loginSessionInfo"

	// BasicSystemInfo URI
	UnityApiBasicSysInfoUri = unityApiTypes + "/basicSystemInfo/instances"

	//StorageResource instance Action URI
	UnityApiStorageResourceInstanceActionUri = UnityApiInstancesUri + "/storageResource/%s/action"

	//Modify Volume URIs
	UnityApiModifyLunUri = UnityApiStorageResourceInstanceActionUri + "/modifyLun"

	//Create LUN Thin Clone
	UnityApiCreateLunThinCloneUri = UnityApiStorageResourceInstanceActionUri + "/createLunThinClone"

	// StorageResource resource URIs
	UnityApiModifyStorageResourceUri = UnityApiInstancesUri + "/storageResource/%s"
	// StorageResource Action resource URI
	UnityApiStorageResourceActionUri = unityApiTypes + "/storageResource/action/%s"

	// StorageResource Action resource URIs
	UnityModifyLunUri = UnityApiModifyStorageResourceUri + "/action/modifyLun"

	//Modify Filesystem URIs
	UnityModifyFilesystemUri = UnityApiModifyStorageResourceUri + "/action/modifyFilesystem"

	//Modify NFS Share URIs
	UnityModifyNFSShareUri = UnityApiGetResourceUri + "/action/modify"

	//Snapshot Action resource URIs
	UnityModifySnapshotUri = UnityApiGetResourceUri + "/action/modify"

	//Snapshot Copy Action
	UnityCopySnapshotUri = UnityApiGetResourceUri + "/action/copy"

	//Host Initiator URIs
	UnityListHostInitiatorsUri = unityApiTypes + "/hostInitiator/instances?fields="
	UnityModifyHostInitiators  = unityRootApi + "/instances/hostInitiator/%s/action/modify"

	//Unity Instance Filter
	UnityInstancesFilter = UnityApiInstanceTypeResources + "?filter=%s"

	//Unity Metrics
	UnityMetric              = "metric"
	UnityMetricQueryResult   = "metricQueryResult"
	UnityMetricRealTimeQuery = "metricRealTimeQuery"

	//Action types for URL's
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
)
