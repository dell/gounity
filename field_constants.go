package gounity

const (
	//To Display the Volume fields
	LunDisplayFields = "id,name,description,type,wwn,sizeTotal,sizeUsed,sizeAllocated,hostAccess,pool,tieringPolicy,ioLimitPolicy,isThinEnabled,isDataReductionEnabled,isThinClone,parentSnap,originalParentLun?fields,health"

	//To Display the File System fields
	FileSystemDisplayFields = "id,name,description,type,sizeTotal,isThinEnabled,isDataReductionEnabled,pool,nasServer,storageResource,nfsShare?fields,cifsShare,tieringPolicy,hostIOSize,health"

	//To Display Storage Resource fields
	StorageResourceDisplayFields = "id,name,filesystem"

	//To Display Tenants fields
	TenantDisplayFields = "id,name"

	//To Display the NFS Share fields
	NFSShareDisplayfields = "id,name,filesystem,readOnlyHosts,readWriteHosts,readOnlyRootAccessHosts,rootAccessHosts,exportPaths"

	//To Display the NAS Server fields
	NasServerDisplayfields = "id,name,nfsServer?fields"

	//To Display the Snapshot fields
	SnapshotDisplayFields = "id,name,description,storageResource?,lun,creationTime,expirationTime,lastRefreshTime,state,size,isAutoDelete,accessType,parentSnap"

	//To Display the HostInitiator fields
	HostInitiatorsDisplayFields = "id,health,type,initiatorId,isIgnored,parentHost,paths"

	//To Display the HostIpPort fields
	HostIpPortDisplayFields = "id,address"

	//To Display License Info fields
	LicenseInfoDisplayFields = "isInstalled,isValid"

	//To Display the HostInitiatorPath fields
	HostInitiatorPathDisplayFields = "fcPort"

	//To Display the FC Port fields
	FcPortDisplayFields = "wwn"

	//Host IO Limit display fields
	HostIOLimitFields = "id,name,description"

	//Iscsi IP Interface display fields
	IscsiIPFields = "id,ipAddress,type"

	//Host Display fields
	HostfieldsToQuery = "id,name,description,fcHostInitiators,iscsiHostInitiators,hostIPPorts?fields"

	//Find Storage Pool fields
	StoragePoolFields = "id,name,description,sizeFree,sizeTotal,sizeUsed,sizeSubscribed,hasDataReductionEnabledLuns,hasDataReductionEnabledFs,isFASTCacheEnabled,type,isAllFlash,poolFastVP"
)
