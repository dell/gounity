package api

const (
	MaxResourceNameLength = 63

	AuthorizationHeader = "Authorization"
	XEmcRestClient      = "X-EMC-REST-CLIENT"
	// Base resource URIs
	unityRootApi  = "/api"
	unityApiTypes = unityRootApi + "/types"

	UnityApiInstanceTypeResources = unityApiTypes + "/%s" + "/instances"

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

	// StorageResource resource URIs
	UnityApiModifyStorageResourceUri = UnityApiInstancesUri + "/storageResource/%s"
	// StorageResource Action resource URI
	UnityApiStorageResourceActionUri = unityApiTypes + "/storageResource/action"

	// LUN resource URIs
	UnityApiCreateLunUri = UnityApiStorageResourceActionUri + "/createLun"

	// StorageResource Action resource URIs
	UnityModifyLunUri = UnityApiModifyStorageResourceUri + "/action/modifyLun"

	//To Display the Volume fields
	LunDisplayFields = "id,name,description,type,wwn,sizeTotal,sizeUsed,sizeAllocated,hostAccess,pool,ioLimitPolicy,isThinEnabled,isDataReductionEnabled"

	//To Display the Snapshot fields
	SnapshotDisplayFields = "id,name,description,storageResource,lun,creationTime,expirationTime,lastRefreshTime,state,size"

	//To Display the HostInitiator fields
	HostInitiatorsDisplayFields = "id,health,type,initiatorId,isIgnored,parentHost"
)
