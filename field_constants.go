/*
 Copyright © 2020 Dell Inc. or its subsidiaries. All Rights Reserved.

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

const (
	//LunDisplayFields to display the Volume fields
	LunDisplayFields = "id,name,description,type,wwn,sizeTotal,sizeUsed,sizeAllocated,hostAccess,pool,tieringPolicy,ioLimitPolicy,isThinEnabled,isDataReductionEnabled,isThinClone,parentSnap,originalParentLun?fields,health"

	//FileSystemDisplayFields to display the File System fields
	FileSystemDisplayFields = "id,name,description,type,sizeTotal,isThinEnabled,isDataReductionEnabled,pool,nasServer,storageResource,nfsShare?fields,cifsShare,tieringPolicy,hostIOSize,health"

	//StorageResourceDisplayFields to display Storage Resource fields
	StorageResourceDisplayFields = "id,name,filesystem"

	//TenantDisplayFields to display Tenants fields
	TenantDisplayFields = "id,name"

	//NFSShareDisplayfields to display the NFS Share fields
	NFSShareDisplayfields = "id,name,filesystem,readOnlyHosts,readWriteHosts,readOnlyRootAccessHosts,rootAccessHosts,exportPaths"

	//NasServerDisplayfields to display the NAS Server fields
	NasServerDisplayfields = "id,name,nfsServer?fields"

	//SnapshotDisplayFields to display the Snapshot fields
	SnapshotDisplayFields = "id,name,description,storageResource?,lun,creationTime,expirationTime,lastRefreshTime,state,size,isAutoDelete,accessType,parentSnap"

	//HostInitiatorsDisplayFields to display the HostInitiator fields
	HostInitiatorsDisplayFields = "id,health,type,initiatorId,isIgnored,parentHost,paths"

	//HostIPPortDisplayFields to display the HostIPPort fields
	HostIPPortDisplayFields = "id,address"

	//LicenseInfoDisplayFields to display License Info fields
	LicenseInfoDisplayFields = "isInstalled,isValid"

	//HostInitiatorPathDisplayFields to display the HostInitiatorPath fields
	HostInitiatorPathDisplayFields = "fcPort"

	//FcPortDisplayFields to display the FC Port fields
	FcPortDisplayFields = "wwn"

	//HostIOLimitFields to display host IO limit fields
	HostIOLimitFields = "id,name,description"

	//IscsiIPFields to display Iscsi IP fields
	IscsiIPFields = "id,ipAddress,type"

	//HostfieldsToQuery to display host fields
	HostfieldsToQuery = "id,name,description,fcHostInitiators,iscsiHostInitiators,hostIPPorts?fields"

	//StoragePoolFields to display Storage Pool fields
	StoragePoolFields = "id,name,description,sizeFree,sizeTotal,sizeUsed,sizeSubscribed,hasDataReductionEnabledLuns,hasDataReductionEnabledFs,isFASTCacheEnabled,type,isAllFlash,poolFastVP"

	// SystemCapacityFields to display system capacity details
	SystemCapacityFields = "id,sizeFree,sizeTotal,sizeUsed,sizePreallocated,sizeSubscribed,totalLogicalSize"
)
