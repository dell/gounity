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

// Struct to capture the Lun properties
type LunParameters struct {
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
	IsThinEnabled          bool                   `json:"isThinEnabled"`
	IsDataReductionEnabled bool                   `json:"isDataReductionEnabled"`
	SupportedProtocol      int                    `json:"supportedProtocols"`
	HostIOSize             int                    `json:"hostIOSize"`
	IsCacheDisabled        bool                   `json:"isCacheDisabled"`
	StoragePool            *StoragePoolID         `json:"pool,omitempty"`
	FastVPParameters       *FastVPParameters      `json:"fastVPParameters,omitempty"`
	HostAccess             *[]HostAccess          `json:"hostAccess,omitempty"`
	IoLimitParameters      *HostIoLimitParameters `json:"ioLimitParameters,omitempty"`
	NasServer              *NasServerID           `json:"nasServer"`
	FileEventSettings      FileEventSettings      `json:"fileEventSettings,omitempty"`
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

// Struct to capture the FileSystem create Params
type FileSystemRequestParam struct {
	Name                 string                `json:"name"`
	FileSystemParameters *FileSystemParameters `json:"filesystemParameters"`
}

// Struct to capture the FileSystem properties
type FileSystemParameters struct {
	Size                   string `json:"size"`
	IsThinEnabled          bool   `json:"isThinEnabled"`
	StoragePoolID          string `json:"pool"`
	Name                   string `json:"name"`
	IsDataReductionEnabled bool   `json:"isDataReductionEnabled"`
	NasServer              string `json:"nasServer"`
}

// Struct to capture Host Request
type HostCreateParam struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OsType      string `json:"osType"`
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
	AccessMask    string         `json:"accessMask"`
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
	Name              string                `json:"name,omitempty"`
	StorageResource   *StorageResourceParam `json:"storageResource,omitempty"`
	Description       string                `json:"description,omitempty"`
	RetentionDuration uint64                `json:"retentionDuration,omitempty"`
	IsAutoDelete      bool                  `json:"isAutoDelete"`
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
