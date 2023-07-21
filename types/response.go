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

package types

import (
	"time"
)

// StoragePool Struct to capture the response of StoragePool response
type StoragePool struct {
	StoragePoolContent StoragePoolContent `json:"content"`
}

// StoragePoolContent Struct to capture the StoragePool Content properties
type StoragePoolContent struct {
	ID                          string     `json:"id"`
	Name                        string     `json:"name"`
	Description                 string     `json:"description"`
	FreeCapacity                uint64     `json:"sizeFree"`
	TotalCapacity               uint64     `json:"sizeTotal"`
	UsedCapacity                uint64     `json:"sizeUsed"`
	SubscribedCapacity          uint64     `json:"sizeSubscribed"`
	HasDataReductionEnabledLuns bool       `json:"hasDataReductionEnabledLuns"`
	HasDataReductionEnabledFs   bool       `json:"hasDataReductionEnabledFs"`
	IsFASTCacheEnabled          bool       `json:"isFASTCacheEnabled"`
	Type                        int8       `json:"type"`
	IsAllFlash                  bool       `json:"isAllFlash"`
	PoolFastVP                  PoolFastVP `json:"poolFastVP"`
}

// PoolFastVP struct to capture fastvp property of pool
type PoolFastVP struct {
	Status            int  `json:"status"`
	RelocationRate    int  `json:"relocationRate"`
	Type              int  `json:"type"`
	IsScheduleEnabled bool `json:"isScheduleEnabled"`
}

// ListVolumes Struct to capture the response of StorageResource response
type ListVolumes struct {
	Volumes []Volume `json:"entries"`
}

// Volume struct to capture response of volume
type Volume struct {
	VolumeContent VolumeContent `json:"content"`
}

// VolumeContent struct to capture volume properties
type VolumeContent struct {
	ResourceID             string               `json:"id"`
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
	IsThinClone            bool                 `json:"isThinClone"`
	ParentSnap             ParentSnap           `json:"parentSnap,omitempty"`
	TieringPolicy          int                  `json:"tieringPolicy,omitempty"`
	ParentVolume           StorageResource      `json:"originalParentLun,omitempty"`
	Health                 HealthContent        `json:"health,omitempty"`
}

// ParentSnap to capture Source Snapshot ID
type ParentSnap struct {
	ID string `json:"id"`
}

// Pool struct to capture Pool Id
type Pool struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// HostAccessResponse Struct to capture Host Access in Volume response
type HostAccessResponse struct {
	HostContent HostContent `json:"host"`
	HLU         int         `json:"hlu"`
}

// Link Struct to capture the link response
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// TenantInfo Struct to capture the Tenant Info
type TenantInfo struct {
	Entries []TenantEntry `json:"entries"`
}

// TenantEntry Struct to capture the Tenant Entry
type TenantEntry struct {
	Content TenantContent `json:"content"`
}

// TenantContent Struct to capture the Tenant Content
type TenantContent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BasicSystemInfo Struct to capture the BasicSystemInfo response
type BasicSystemInfo struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Links   []Link    `json:"links"`
	Entries []Entries `json:"entries"`
}

// Entries Struct to capture Entries contains a list of Links.
type Entries struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Links   []Link    `json:"links"`
	Content Content   `json:"content"`
}

// Content Struct to capture the SystemInfo Content.
type Content struct {
	ID                 string `json:"id"`
	Model              string `json:"model"`
	Name               string `json:"name"`
	SoftwareVersion    string `json:"softwareVersion"`
	APIVersion         string `json:"apiVersion"`
	EarliestAPIVersion string `json:"earliestApiVersion"`
}

// Host struct to capture host object
type Host struct {
	HostContent HostContent `json:"content"`
}

// HostContent struct to capture host parameters
type HostContent struct {
	ID              string       `json:"id"`
	Name            string       `json:"name,omitempty"`
	Description     string       `json:"description,omitempty"`
	FcInitiators    []Initiators `json:"fcHostInitiators,omitempty"`
	IscsiInitiators []Initiators `json:"iscsiHostInitiators,omitempty"`
	IPPorts         []IPPorts    `json:"hostIPPorts,omitempty"`
	Address         string       `json:"address,omitempty"`
}

// Initiators struct to capture Initiator ID
type Initiators struct {
	ID string `json:"id"`
}

// IPPorts struct to capture IpPort ID
type IPPorts struct {
	ID      string `json:"id"`
	Address string `json:"address,omitempty"`
}

// HostIPPort struct to capture Host IP port object
type HostIPPort struct {
	HostIPContent HostContent `json:"content"`
}

// ListHostInitiator struct to capture host initiators
type ListHostInitiator struct {
	HostInitiator []HostInitiator `json:"entries"`
}

// HostInitiator struct to capture host initiator object
type HostInitiator struct {
	HostInitiatorContent HostInitiatorContent `json:"content"`
}

// HostInitiatorContent struct to capture host initiator parameters
type HostInitiatorContent struct {
	ID          string        `json:"id"`
	Health      HealthContent `json:"health"`
	Type        int           `json:"type"`
	InitiatorID string        `json:"InitiatorId"`
	IsIgnored   bool          `json:"isIgnored"`
	ParentHost  HostContent   `json:"parentHost"`
	Paths       []Path        `json:"paths"`
}

// Path struct to capture Path ID
type Path struct {
	ID string `json:"id"`
}

// HealthContent to capture health status
type HealthContent struct {
	Value          int      `json:"value"`
	DescriptionIDs []string `json:"descriptionIds"`
	Descriptions   []string `json:"descriptions"`
}

// ListSnapshot struct to capture snapshot list
type ListSnapshot struct {
	Snapshots []Snapshot `json:"entries"`
}

// Snapshot struct to capture snapshot object
type Snapshot struct {
	SnapshotContent SnapshotContent `json:"content"`
}

// SnapshotContent struct to capture snapshot parameters
type SnapshotContent struct {
	ResourceID      string          `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description,omitempty"`
	StorageResource StorageResource `json:"storageResource,omitempty"`
	CreationTime    time.Time       `json:"creationTime,omitempty"`
	ExpirationTime  time.Time       `json:"expirationTime,omitempty"`
	LastRefreshTime time.Time       `json:"lastRefreshTime,omitempty"`
	State           int             `json:"state,omitempty"`
	Size            int64           `json:"size"`
	IsAutoDelete    bool            `json:"isAutoDelete"`
	AccessType      int             `json:"accessType,omitempty"`
	ParentSnap      StorageResource `json:"parentSnap,omitempty"`
}

// CopySnapshots struct to capture copy snapshot content
type CopySnapshots struct {
	CopySnapshotsContent CopySnapshotsContent `json:"content"`
}

// CopySnapshotsContent struct to capture copies list
type CopySnapshotsContent struct {
	Copies []StorageResource `json:"copies,omitempty"`
}

// StorageResource struct to capture storage resource ID
type StorageResource struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// StorageResourceParameters struct to capture Storage Resource content
type StorageResourceParameters struct {
	StorageResourceContent StorageResourceContent `json:"content"`
}

// StorageResourceContent struct to capture Storage Resource content
type StorageResourceContent struct {
	ID         string          `json:"id"`
	Name       string          `json:"name,omitempty"`
	Filesystem StorageResource `json:"filesystem,omitempty"`
}

// IoLimitPolicy struct IO limit policy object
type IoLimitPolicy struct {
	IoLimitPolicyContent IoLimitPolicyContent `json:"content,omitempty"`
}

// IoLimitPolicyContent struct to capture IoLimitPolicyContent parameters
type IoLimitPolicyContent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Filesystem struct to capture filesystem object
type Filesystem struct {
	FileContent FileContent `json:"content"`
}

// FileContent struct to capture filesystem parameters
type FileContent struct {
	ID                     string        `json:"id"`
	Name                   string        `json:"name,omitempty"`
	SizeTotal              uint64        `json:"sizeTotal,omitempty"`
	Description            string        `json:"description,omitempty"`
	Type                   int           `json:"type,omitempty"`
	Format                 int           `json:"format,omitempty"`
	HostIOSize             int64         `json:"hostIOSize,omitempty"`
	TieringPolicy          uint64        `json:"tieringPolicy,omitempty"`
	IsThinEnabled          bool          `json:"isThinEnabled"`
	IsDataReductionEnabled bool          `json:"isDataReductionEnabled"`
	Pool                   Pool          `json:"pool,omitempty"`
	NASServer              Pool          `json:"nasServer,omitempty"`
	StorageResource        Pool          `json:"storageResource,omitempty"`
	NFSShare               []Share       `json:"nfsShare,omitempty"`
	CIFSShare              []Pool        `json:"cifsShare,omitempty"`
	Health                 HealthContent `json:"health,omitempty"`
}

// Share object to capture NFS Share object from FileContent
type Share struct {
	ID         string          `json:"id"`
	Name       string          `json:"name,omitempty"`
	Path       string          `json:"path,omitempty"`
	ParentSnap StorageResource `json:"snap,omitempty"`
}

// NFSShare struct to capture NFS Share object
type NFSShare struct {
	NFSShareContent NFSShareContent `json:"content"`
}

// NFSShareContent struct to capture NFS Share parameters
type NFSShareContent struct {
	ID                      string        `json:"id"`
	Name                    string        `json:"name,omitempty"`
	Filesystem              Pool          `json:"filesystem,omitempty"`
	ReadOnlyHosts           []HostContent `json:"readOnlyHosts,omitempty"`
	ReadWriteHosts          []HostContent `json:"readWriteHosts,omitempty"`
	ReadOnlyRootAccessHosts []HostContent `json:"readOnlyRootAccessHosts,omitempty"`
	RootAccessHosts         []HostContent `json:"rootAccessHosts,omitempty"`
	ExportPaths             []string      `json:"exportPaths,omitempty"`
}

// NASServer struct to capture NAS Server object
type NASServer struct {
	NASServerContent NASServerContent `json:"content"`
}

// NASServerContent struct to capture NAS Server object
type NASServerContent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	NFSServer NFSServer `json:"nfsServer,omitempty"`
}

// NFSServer struct to capture NFS Server object
type NFSServer struct {
	ID           string `json:"id"`
	Name         string `json:"name,omitempty"`
	NFSv3Enabled bool   `json:"nfsv3Enabled"`
	NFSv4Enabled bool   `json:"nfsv4Enabled"`
}

// ListIPInterfaces struct to capture snapshot list
type ListIPInterfaces struct {
	Entries []IPInterfaceEntries `json:"entries"`
}

// IPInterfaceEntries struct to capture IpInterface object
type IPInterfaceEntries struct {
	IPInterfaceContent IPInterfaceContent `json:"content"`
}

// IPInterfaceContent struct to capture IpInterface parameters
type IPInterfaceContent struct {
	ID        string `json:"id"`
	IPAddress string `json:"ipAddress"`
	Type      int    `json:"type"`
}

// LicenseInfo for features on Array
type LicenseInfo struct {
	LicenseInfoContent LicenseInfoContent `json:"content"`
}

// LicenseInfoContent for features on Array
type LicenseInfoContent struct {
	IsInstalled bool `json:"isInstalled"`
	IsValid     bool `json:"isValid"`
}

// HostInitiatorPath struct to capture host initiator path object
type HostInitiatorPath struct {
	HostInitiatorPathContent HostInitiatorPathContent `json:"content"`
}

// HostInitiatorPathContent struct to capture host initiator parameters
type HostInitiatorPathContent struct {
	FcPortID FcPortID `json:"fcPort"`
}

// FcPortID struct to capture FC port ID
type FcPortID struct {
	ID string `json:"id"`
}

// FcPort struct to capture FC port object
type FcPort struct {
	FcPortContent FcPortContent `json:"content"`
}

// FcPortContent struct to capture FC port ID
type FcPortContent struct {
	Wwn string `json:"wwn"`
}

// MetricRealTimeQuery is body of a request to create a MetricCollection query
type MetricRealTimeQuery struct {
	Paths    []string `json:"paths"`
	Interval int      `json:"interval"`
}

// MetricQueryResponseContent is part of response to creating a MetricCollection query
type MetricQueryResponseContent struct {
	MaximumSamples int      `json:"maximumSamples"`
	Expiration     string   `json:"expiration"`
	Interval       int      `json:"interval"`
	Paths          []string `json:"paths"`
	ID             int      `json:"id"`
}

// MetricQueryCreateResponse a response from creating a MetricCollection query
type MetricQueryCreateResponse struct {
	Base    string                     `json:"base"`
	Updated string                     `json:"updated"`
	Content MetricQueryResponseContent `json:"content"`
}

// MetricResult is part of response of a MetricCollection query
type MetricResult struct {
	QueryID   int                    `json:"queryId"`
	Path      string                 `json:"path"`
	Timestamp string                 `json:"timestamp"`
	Values    map[string]interface{} `json:"values"`
}

// MetricResultEntry is part of response of a MetricCollection query
type MetricResultEntry struct {
	Base    string       `json:"base"`
	Updated string       `json:"updated"`
	Content MetricResult `json:"content"`
}

// MetricQueryResult is response from querying a MetricCollection
type MetricQueryResult struct {
	Base    string              `json:"base"`
	Updated string              `json:"updated"`
	Entries []MetricResultEntry `json:"entries"`
}

// MetricContent is part of the response from /api/types/metric/instances
type MetricContent struct {
	ID int `json:"id"`
}

// MetricEntries is part of the response from /api/types/metric/instances
type MetricEntries struct {
	Cnt MetricContent `json:"content"`
}

// MetricPaths comes from response from /api/types/metric/instances
type MetricPaths struct {
	Entries []MetricEntries `json:"entries"`
}

// MetricInfo has all the details of instance of a Unity metric
type MetricInfo struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Path                  string `json:"path"`
	Product               int    `json:"product"`
	Type                  int    `json:"type"`
	Description           string `json:"description"`
	IsHistoricalAvailable bool   `json:"isHistoricalAvailable"`
	IsRealtimeAvailable   bool   `json:"isRealtimeAvailable"`
	Unit                  int    `json:"unit"`
	UnitDisplayString     string `json:"unitDisplayString"`
	Visibility            int    `json:"visibility"`
}

// MetricInstance describes an instance of Unity metric
type MetricInstance struct {
	Base    string     `json:"base"`
	Updated string     `json:"updated"`
	Content MetricInfo `json:"content"`
}
