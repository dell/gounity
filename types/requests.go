/*
Copyright (c) 2019 Dell Corporation
All Rights Reserved
*/

package types

import (
	"fmt"
)

//ErrorContent Struct to capture the Error information.
type ErrorContent struct {
	Message        []ErrorMessage `json:"messages"`
	HTTPStatusCode int            `json:"httpStatusCode"`
	ErrorCode      int            `json:"errorCode"`
}

//ErrorMessage Struct to cature error message
type ErrorMessage struct {
	EnUS string `json:"en-US"`
}

//Error Struct to cature error
type Error struct {
	ErrorContent ErrorContent `json:"error"`
}

//Error function returns the error message.
func (e Error) Error() string {
	return fmt.Sprintf("%v", e.ErrorContent.Message)
}

// Struct to capture the StoragePool properties
//type StoragePool struct {
//	ID string `json:"id"`
//}

//FastVPParameters Struct to capture Tiering Policy for Create Volume
type FastVPParameters struct {
	TieringPolicy int `json:"tieringPolicy"`
}

//StoragePoolID Struct to capture Storage pool ID for Create Volume
type StoragePoolID struct {
	PoolID string `json:"id"`
}

//NasServerID Struct to capture Nas server ID for Create Volume
type NasServerID struct {
	NasServerID string `json:"id"`
}

//LunCreateParam Struct to capture the Lun create Params
type LunCreateParam struct {
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	LunParameters *LunParameters `json:"lunParameters"`
}

// ConsistencyGroupCreate create consistency group request
type ConsistencyGroupCreate struct {
	Name                   string                   `json:"name"`
	Description            string                   `json:"description,omitempty"`
	ReplicationParameters  *CGReplicationParameters `json:"replicationParameters,omitempty"`
	LunAdd                 *[]CGLunAdd              `json:"lunAdd,omitempty"`
	LunCreate              *[]CGLunCreate           `json:"lunCreate,omitempty"`
	SnapScheduleParameters *SnapScheduleParameters  `json:"snapScheduleParameters,omitempty"`
}

// ConsistencyGroupModify create consistency group request
type ConsistencyGroupModify struct {
	Name                   string                  `json:"name,omitempty"`
	Description            string                  `json:"description,omitempty"`
	LunAdd                 *[]CGLunAdd             `json:"lunAdd,omitempty"`
	LunCreate              *[]CGLunModify          `json:"lunCreate,omitempty"`
	SnapScheduleParameters *SnapScheduleParameters `json:"snapScheduleParameters,omitempty"`
}

// CGReplicationParameters to work with IsReplicationDestination parameter
type CGReplicationParameters struct {
	IsReplicationDestination bool `json:"isReplicationDestination,omitempty"`
}

// CGLunAdd structure
type CGLunAdd struct {
	Lun *CGLun `json:"lun,omitempty"`
}

// CGLun structure
type CGLun struct {
	ID string `json:"id,omitempty"`
}

// CGLunCreate structure
type CGLunCreate struct {
	Name      string             `json:"name"`
	LunParams *CGLunCreateParams `json:"lunParameters"`
}

// CGLunModify structure
type CGLunModify struct {
	Name      string             `json:"name"`
	LunParams *CGLunModifyParams `json:"lunParameters"`
}

// CGLunCreateParams structure
type CGLunCreateParams struct {
	Size uint64 `json:"size"`
	// (Applies only to create requests.)
	Pool *CGLunCretePool `json:"pool"`
}

// CGLunModifyParams structure
type CGLunModifyParams struct {
	Size uint64 `json:"size"`
}

// CGLunCretePool structure
type CGLunCretePool struct {
	ID string `json:"id"`
}

// SnapScheduleParameters structure
type SnapScheduleParameters struct {
	SnapSchedule         *SnapSchedule `json:"snapSchedule,omitempty"`
	IsSnapSchedulePaused bool          `json:"isSnapSchedulePaused,omitempty"`
}

// SnapSchedule structure
type SnapSchedule struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

//Tenants Struct to capture the Tenants
type Tenants struct {
	TenantID string `json:"id"`
}

//LunParameters Struct to capture the Lun properties
type LunParameters struct {
	Name                   string                 `json:"name,omitempty"`
	Size                   uint64                 `json:"size,omitempty"`
	IsThinEnabled          string                 `json:"isThinEnabled,omitempty"`
	StoragePool            *StoragePoolID         `json:"pool,omitempty"`
	IsDataReductionEnabled string                 `json:"isDataReductionEnabled,omitempty"`
	FastVPParameters       *FastVPParameters      `json:"fastVPParameters,omitempty"`
	HostAccess             *[]HostAccess          `json:"hostAccess,omitempty"`
	IoLimitParameters      *HostIoLimitParameters `json:"ioLimitParameters,omitempty"`
}

//FsCreateParam Struct to capture the Filesystem create Params
type FsCreateParam struct {
	Name                  string                 `json:"name"`
	Description           string                 `json:"description,omitempty"`
	FsParameters          *FsParameters          `json:"fsParameters"`
	ReplicationParameters *ReplicationParameters `json:"replicationParameters,omitempty"`
}

//FsParameters Struct to capture the File system properties
type FsParameters struct {
	Size                   uint64                 `json:"size,omitempty"`
	IsThinEnabled          string                 `json:"isThinEnabled,omitempty"`
	IsDataReductionEnabled string                 `json:"isDataReductionEnabled,omitempty"`
	SupportedProtocol      int                    `json:"supportedProtocols"`
	HostIOSize             int                    `json:"hostIOSize"`
	StoragePool            *StoragePoolID         `json:"pool,omitempty"`
	FastVPParameters       *FastVPParameters      `json:"fastVPParameters,omitempty"`
	HostAccess             *[]HostAccess          `json:"hostAccess,omitempty"`
	IoLimitParameters      *HostIoLimitParameters `json:"ioLimitParameters,omitempty"`
	NasServer              *NasServerID           `json:"nasServer"`
	FileEventSettings      FileEventSettings      `json:"fileEventSettings,omitempty"`
}

//ReplicationParameters Struct for setting replication destination parameters
type ReplicationParameters struct {
	IsReplicationDestination bool `json:"isReplicationDestination"`
}

//FsExpandParameters Struct to capture expand Filesystem parameters
type FsExpandParameters struct {
	Size uint64 `json:"size"`
}

//FsExpandModifyParam Struct to expand Filesystem
type FsExpandModifyParam struct {
	FsParameters *FsExpandParameters `json:"fsParameters"`
}

//FsModifyParameters Struct to modify Filesystem parameters
type FsModifyParameters struct {
	NFSShares   *[]NFSShareCreateParam `json:"nfsShareCreate,omitempty"`
	Description string                 `json:"description,omitempty"`
}

//NFSShareCreateParam Struct to capture NFS Share Create parameters
type NFSShareCreateParam struct {
	Name               string              `json:"name"`
	Path               string              `json:"path"`
	NFSShareParameters *NFSShareParameters `json:"nfsShareParameters,omitempty"`
}

//NFSShareCreateFromSnapParam Struct to capture create NFS share from snapshot parameters
type NFSShareCreateFromSnapParam struct {
	Name          string            `json:"name"`
	Path          string            `json:"path"`
	DefaultAccess string            `json:"defaultAccess,omitempty"`
	Snapshot      SnapshotIDContent `json:"snap"`
}

//NFSShareModify Struct to modify NFS Share parameters
type NFSShareModify struct {
	NFSSharesModifyContent *[]NFSShareModifyContent `json:"nfsShareModify,omitempty"`
}

//NFSShareCreateFromSnapModify Struct to modify NFS Share created from snapshot parameters
type NFSShareCreateFromSnapModify struct {
	DefaultAccess           string           `json:"defaultAccess,omitempty"`
	ReadOnlyHosts           *[]HostIDContent `json:"readOnlyHosts,omitempty"`
	ReadWriteHosts          *[]HostIDContent `json:"readWriteHosts,omitempty"`
	ReadOnlyRootAccessHosts *[]HostIDContent `json:"readOnlyRootAccessHosts,omitempty"`
	RootAccessHosts         *[]HostIDContent `json:"rootAccessHosts,omitempty"`
}

//NFSShareDelete Struct to modify NFS Share parameters
type NFSShareDelete struct {
	NFSSharesDeleteContent *[]NFSShareModifyContent `json:"nfsShareDelete,omitempty"`
}

//NFSShareModifyContent Struct to capture NFS Share modify content
type NFSShareModifyContent struct {
	NFSShare           *StorageResourceParam `json:"nfsShare,omitempty"`
	NFSShareParameters *NFSShareParameters   `json:"nfsShareParameters,omitempty"`
}

//NFSShareParameters Struct to capture NFS Share properties
type NFSShareParameters struct {
	DefaultAccess           string           `json:"defaultAccess,omitempty"`
	ReadOnlyHosts           *[]HostIDContent `json:"readOnlyHosts,omitempty"`
	ReadWriteHosts          *[]HostIDContent `json:"readWriteHosts,omitempty"`
	ReadOnlyRootAccessHosts *[]HostIDContent `json:"readOnlyRootAccessHosts,omitempty"`
	RootAccessHosts         *[]HostIDContent `json:"rootAccessHosts,omitempty"`
}

//FileEventSettings Struct to capture File event settings
type FileEventSettings struct {
	IsCIFSEnabled bool `json:"isCIFSEnabled"`
	IsNFSEnabled  bool `json:"isNFSEnabled"`
}

//LunExpandParameters to capture Lun expand parameters
type LunExpandParameters struct {
	Size uint64 `json:"size,omitempty"`
}

//LunHostAccessParameters to capture Lun Host Access parameters
type LunHostAccessParameters struct {
	HostAccess *[]HostAccess `json:"hostAccess,omitempty"`
}

//HostCreateParam Struct to capture Host Request
type HostCreateParam struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OsType      string   `json:"osType"`
	Tenant      *Tenants `json:"tenant,omitempty"`
}

//HostIDContent Struct to capture Host ID Content
type HostIDContent struct {
	ID string `json:"id"`
}

//HostIPPortCreateParam Struct to capture Host IP Pot Request
type HostIPPortCreateParam struct {
	HostIDContent *HostIDContent `json:"host"`
	Address       string         `json:"address"`
}

//HostInitiatorCreateParam Struct to capture Host Initiator create parameters
type HostInitiatorCreateParam struct {
	HostIDContent *HostIDContent `json:"host"`
	InitiatorType InitiatorType  `json:"initiatorType"`
	InitiatorWwn  string         `json:"initiatorWWNorIqn"`
}

//HostInitiatorModifyParam Struct to capture Host Initiator modify parameters
type HostInitiatorModifyParam struct {
	HostIDContent *HostIDContent `json:"host"`
}

//HostAccess Struct to capture Host access parameters
type HostAccess struct {
	HostIDContent *HostIDContent `json:"host"`
	AccessMask    string         `json:"accessMask,omitempty"`
}

//LunModifyParam Struct to capture Lun modify parameters
type LunModifyParam struct {
	LunParameters *LunParameters `json:"lunParameters"`
}

//LunExpandModifyParam Struct to capture Lun expand modify parameters
type LunExpandModifyParam struct {
	LunParameters *LunExpandParameters `json:"lunParameters"`
}

//LunHostAccessModifyParam Struct to capture Lun host access modify parameters
type LunHostAccessModifyParam struct {
	LunHostAccessParameters *LunHostAccessParameters `json:"lunParameters"`
}

//CreateSnapshotParam struct to capture create snapshot parameters
type CreateSnapshotParam struct {
	Name                 string                `json:"name,omitempty"`
	StorageResource      *StorageResourceParam `json:"storageResource,omitempty"`
	Description          string                `json:"description,omitempty"`
	RetentionDuration    uint64                `json:"retentionDuration,omitempty"`
	IsAutoDelete         bool                  `json:"isAutoDelete"`
	FilesystemAccessType int                   `json:"filesystemAccessType,omitempty"`
}

//CopySnapshot struct to capture Copy snapshot parameters
type CopySnapshot struct {
	Name  string `json:"copyName,omitempty"`
	Child bool   `json:"child"`
}

//StorageResourceParam struct to capture storage resource parameters
type StorageResourceParam struct {
	ID string `json:"id"`
}

//HostIoLimitParameters struct to capture HostIO Limit parameters
type HostIoLimitParameters struct {
	IoLimitPolicyParam *IoLimitPolicyParam `json:"ioLimitPolicy"`
}

//IoLimitPolicyParam struct to capture IOLimit Policy Parameters
type IoLimitPolicyParam struct {
	ID string `json:"id"`
}

//SnapshotIDContent struct to capture Snapshot ID Content
type SnapshotIDContent struct {
	ID string `json:"id"`
}

//CreateLunThinCloneParam struct to capture Create LUN thin clone Parameters
type CreateLunThinCloneParam struct {
	SnapIDContent *SnapshotIDContent `json:"snap"`
	Name          string             `json:"name"`
}

//InitiatorType is string Type
type InitiatorType string

// CreateReplicationSessionParam structure
type CreateReplicationSessionParam struct {
	Name                 string               `json:"name"`
	MaxTimeOutOfSync     string               `json:"maxTimeOutOfSync"`
	SrcResourceID        string               `json:"srcResourceId"`
	DstResourceID        string               `json:"dstResourceId"`
	RemoteSystem         *RemoteSystemContent `json:"remoteSystem,omitempty"`
	OverwriteDestination bool                 `json:"overwriteDestination,omitempty"`
}
