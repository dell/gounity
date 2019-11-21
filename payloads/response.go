package payloads

import (
	"time"
)

//Commented StoragePool as it has a conflict in requests.go
// Struct to capture the response of StoragePool response
type StoragePool struct {
	StoragePoolContent StoragePoolContent `json:"content"`
}

// Struct to capture the StoragePool Content properties
type StoragePoolContent struct {
	ID                        string     `json:"id"`
	Name                      string     `json:"name"`
	Description               string     `json:"description"`
	FreeCapacity              uint64     `json:"sizeFree"`
	TotalCapacity             uint64     `json:"sizeTotal"`
	UsedCapacity              uint64     `json:"sizeUsed"`
	SubscribedCapacity        uint64     `json:"sizeSubscribed"`
	HasCompressionEnabledLuns bool       `json:"hasCompressionEnabledLuns"`
	HasCompressionEnabledFs   bool       `json:"hasCompressionEnabledFs"`
	IsFASTCacheEnabled        bool       `json:"isFASTCacheEnabled"`
	Type                      int8       `json:"type"`
	IsAllFlash                bool       `json:"isAllFlash"`
	PoolFastVP                PoolFastVP `json:"poolFastVP"`
}

type PoolFastVP struct {
	Status            int  `json:"status"`
	RelocationRate    int  `json:"relocationRate"`
	Type              int  `json:"type"`
	IsScheduleEnabled bool `json:"isScheduleEnabled"`
}

// Struct to capture the response of StorageResource response
type ListVolumes struct {
	Volumes []Volume `json:"entries"`
}

type Volume struct {
	VolumeContent VolumeContent `json:"content"`
}

type VolumeContent struct {
	ResourceId             string               `json:"id"`
	Name                   string               `json:"name,omitempty"`
	Description            string               `json:"description,omitempty"`
	Type                   int                  `json:"type,omitempty"`
	SizeTotal              uint64               `json:"sizeTotal,omitempty"`
	SizeUsed               uint64               `json:"sizeUsed,omitempty"`
	SizeAllocated          uint64               `json:"sizeAllocated,omitempty"`
	HostAccessResponse     []HostAccessResponse `json:"hostAccess,omitempty"`
	Wwn                    string               `json:"wwn,omitempty"`
	Pool                   Pool                 `json:"pool,omitempty"`
	IsThinEnabled          bool                 `json:"isThinEnabled"`
	IsDataReductionEnabled bool                 `json:"isDataReductionEnabled"`
	IoLimitPolicyContent   IoLimitPolicyContent `json:"ioLimitPolicy,omitempty"`
}

type Pool struct {
	Id string `json:"id"`
}

// Struct to capture Host Access in Volume response
type HostAccessResponse struct {
	HostContent HostContent `json:"host"`
}

// Struct to capture the link response
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// Struct to capture the BasicSystemInfo response
type BasicSystemInfo struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Links   []Link    `json:"links"`
	Entries []Entries `json:"entries"`
}

// Struct to capture Entries contains a list of Links.
type Entries struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Links   []Link    `json:"links"`
	Content Content   `json:"content"`
}

// Struct to capture the SystemInfo Content.
type Content struct {
	ID                 string `json:"id"`
	Model              string `json:"model"`
	Name               string `json:"name"`
	SoftwareVersion    string `json:"softwareVersion"`
	APIVersion         string `json:"apiVersion"`
	EarliestAPIVersion string `json:"earliestApiVersion"`
}

type Host struct {
	HostContent HostContent `json:"content"`
}

type HostContent struct {
	ID           string         `json:"id"`
	Name         string         `json:"name,omitempty"`
	Description  string         `json:"description,omitempty"`
	FcInitiators []FcInitiators `json:"fcHostInitiators,omitempty"`
}

type FcInitiators struct {
	Id string `json:"id"`
}

type HostIpPort struct {
	HostContent HostContent `json:"content"`
}

type ListHostInitiator struct {
	HostInitiator []HostInitiator `json:"entries"`
}

type HostInitiator struct {
	HostInitiatorContent HostInitiatorContent `json:"content"`
}

type HostInitiatorContent struct {
	Id          string      `json:"id"`
	Health      string      `json:"string"`
	Type        int         `json:"type"`
	InitiatorId string      `json:"InitiatorId"`
	IsIgnored   bool        `json:"isIgnored"`
	ParentHost  HostContent `json:"parentHost"`
}

type ListSnapshot struct {
	Snapshots []Snapshot `json:"entries"`
}

type Snapshot struct {
	SnapshotContent SnapshotContent `json:"content"`
}

type SnapshotContent struct {
	ResourceId      string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	Lun             Lun       `json:"lun,omitempty"`
	CreationTime    time.Time `json:"creationTime,omitempty"`
	ExpirationTime  time.Time `json:"expirationTime,omitempty"`
	LastRefreshTime time.Time `json:"lastRefreshTime,omitempty"`
	State           int       `json:"state,omitempty"`
	Size            int64     `json:"size"`
}

type StorageResource struct {
	Id string `json:"id"`
}

type Lun struct {
	Id string `json:"id"`
}

type IoLimitPolicy struct {
	IoLimitPolicyContent IoLimitPolicyContent `json:"content,omitempty"`
}

type IoLimitPolicyContent struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
