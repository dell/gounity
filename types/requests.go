/*
Copyright (c) 2019 Dell EMC Corporation
All Rights Reserved
*/
package types

import "fmt"

// Struct to capture the Error information.
type ErrorContent struct {
	Message        []ErrorMessage `json:"messages"`
	HTTPStatusCode int            `json:"httpStatusCode"`
	ErrorCode      int            `json:"errorCode"`
}

type ErrorMessage struct {
	EnUS string `json:"en-US"`
}

type Error struct {
	ErrorContent ErrorContent `json:"error"`
}

// Return the error message from the given Error object.
func (e Error) Error() string {
	return fmt.Sprintf("%v", e.ErrorContent.Message)
}

// Struct to capture the StoragePool properties
//type StoragePool struct {
//	ID string `json:"id"`
//}

//Struct to capture Tiering Policy for Create Volume
type FastVPParameters struct {
	TieringPolicy int `json:"tieringPolicy"`
}

//Struct to capture Storage pool Id for Create Volume
type StoragePoolID struct {
	PoolId string `json:"id"`
}

//Struct to capture Nas server Id for Create Volume
type NasServerID struct {
	NasServerID string `json:"id"`
}

// Struct to capture the Lun create Params
type LunCreateParam struct {
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	LunParameters *LunParameters `json:"lunParameters"`
}

// Struct to capture the Tenants
type Tenants struct {
	TenantId string `json:"id"`
}

// Struct to capture the Lun properties
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

// Struct to capture the Filesystem create Params
type FsCreateParam struct {
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	FsParameters *FsParameters `json:"fsParameters"`
}

// Struct to capture the File system properties
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

//Struct to capture expand Filesystem parameters
type FsExpandParameters struct {
	Size uint64 `json:"size"`
}

//Struct to expand Filesystem
type FsExpandModifyParam struct {
	FsParameters *FsExpandParameters `json:"fsParameters"`
}

// Struct to modify Filesystem parameters
type FsModifyParameters struct {
	NFSShares   *[]NFSShareCreateParam `json:"nfsShareCreate,omitempty"`
	Description string                 `json:"description,omitempty"`
}

// Struct to capture NFS Share Create parameters
type NFSShareCreateParam struct {
	Name               string              `json:"name"`
	Path               string              `json:"path"`
	NFSShareParameters *NFSShareParameters `json:"nfsShareParameters,omitempty"`
}

type NFSShareCreateFromSnapParam struct {
	Name          string            `json:"name"`
	Path          string            `json:"path"`
	DefaultAccess string            `json:"defaultAccess,omitempty"`
	Snapshot      SnapshotIdContent `json:"snap"`
}

//Struct to modify NFS Share parameters
type NFSShareModify struct {
	NFSSharesModifyContent *[]NFSShareModifyContent `json:"nfsShareModify,omitempty"`
}

type NFSShareCreateFromSnapModify struct {
	DefaultAccess           string           `json:"defaultAccess,omitempty"`
	ReadOnlyHosts           *[]HostIdContent `json:"readOnlyHosts,omitempty"`
	ReadWriteHosts          *[]HostIdContent `json:"readWriteHosts,omitempty"`
	ReadOnlyRootAccessHosts *[]HostIdContent `json:"readOnlyRootAccessHosts,omitempty"`
	RootAccessHosts         *[]HostIdContent `json:"rootAccessHosts,omitempty"`
}

//Struct to modify NFS Share parameters
type NFSShareDelete struct {
	NFSSharesDeleteContent *[]NFSShareModifyContent `json:"nfsShareDelete,omitempty"`
}

//Struct to capture NFS Share modify content
type NFSShareModifyContent struct {
	NFSShare           *StorageResourceParam `json:"nfsShare,omitempty"`
	NFSShareParameters *NFSShareParameters   `json:"nfsShareParameters,omitempty"`
}

// Struct to capture NFS Share properties
type NFSShareParameters struct {
	DefaultAccess           string           `json:"defaultAccess,omitempty"`
	ReadOnlyHosts           *[]HostIdContent `json:"readOnlyHosts,omitempty"`
	ReadWriteHosts          *[]HostIdContent `json:"readWriteHosts,omitempty"`
	ReadOnlyRootAccessHosts *[]HostIdContent `json:"readOnlyRootAccessHosts,omitempty"`
	RootAccessHosts         *[]HostIdContent `json:"rootAccessHosts,omitempty"`
}

// Struct to capture File event settings
type FileEventSettings struct {
	IsCIFSEnabled bool `json:"isCIFSEnabled"`
	IsNFSEnabled  bool `json:"isNFSEnabled"`
}

type LunExpandParameters struct {
	Size uint64 `json:"size,omitempty"`
}
type LunHostAccessParameters struct {
	HostAccess *[]HostAccess `json:"hostAccess,omitempty"`
}

// Struct to capture Host Request
type HostCreateParam struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OsType      string   `json:"osType"`
	Tenant      *Tenants `json:"tenant,omitempty"`
}

type HostIdContent struct {
	ID string `json:"id"`
}

// Struct to capture Host Ip Pot Request
type HostIpPortCreateParam struct {
	HostIdContent *HostIdContent `json:"host"`
	Address       string         `json:"address"`
}

type HostInitiatorCreateParam struct {
	HostIdContent *HostIdContent `json:"host"`
	InitiatorType InitiatorType  `json:"initiatorType"`
	InitiatorWwn  string         `json:"initiatorWWNorIqn"`
}

type HostInitiatorModifyParam struct {
	HostIdContent *HostIdContent `json:"host"`
}

type HostAccess struct {
	HostIdContent *HostIdContent `json:"host"`
	AccessMask    string         `json:"accessMask,omitempty"`
}

type LunModifyParam struct {
	LunParameters *LunParameters `json:"lunParameters"`
}

type LunExpandModifyParam struct {
	LunParameters *LunExpandParameters `json:"lunParameters"`
}

type LunHostAccessModifyParam struct {
	LunHostAccessParameters *LunHostAccessParameters `json:"lunParameters"`
}

type CreateSnapshotParam struct {
	Name                 string                `json:"name,omitempty"`
	StorageResource      *StorageResourceParam `json:"storageResource,omitempty"`
	Description          string                `json:"description,omitempty"`
	RetentionDuration    uint64                `json:"retentionDuration,omitempty"`
	IsAutoDelete         bool                  `json:"isAutoDelete"`
	FilesystemAccessType int                   `json:"filesystemAccessType,omitempty"`
}

type CopySnapshot struct {
	Name  string `json:"copyName,omitempty"`
	Child bool   `json:"child"`
}

type StorageResourceParam struct {
	ID string `json:"id"`
}

type HostIoLimitParameters struct {
	IoLimitPolicyParam *IoLimitPolicyParam `json:"ioLimitPolicy"`
}

type IoLimitPolicyParam struct {
	Id string `json:"id"`
}

type SnapshotIdContent struct {
	Id string `json:"id"`
}

type CreateLunThinCloneParam struct {
	SnapIdContent *SnapshotIdContent `json:"snap"`
	Name          string             `json:"name"`
}

type InitiatorType string
