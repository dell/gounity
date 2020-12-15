package types

import (
	"time"
)

//StoragePool Struct to capture the response of StoragePool response
type StoragePool struct {
	StoragePoolContent StoragePoolContent `json:"content"`
}

//StoragePoolContent Struct to capture the StoragePool Content properties
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

//PoolFastVP struct to capture fastvp property of pool
type PoolFastVP struct {
	Status            int  `json:"status"`
	RelocationRate    int  `json:"relocationRate"`
	Type              int  `json:"type"`
	IsScheduleEnabled bool `json:"isScheduleEnabled"`
}

//ListVolumes Struct to capture the response of StorageResource response
type ListVolumes struct {
	Volumes []Volume `json:"entries"`
}

//Volume struct to capture response of volume
type Volume struct {
	VolumeContent VolumeContent `json:"content"`
}

//VolumeContent struct to capture volume properties
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
	IsThinClone            bool                 `json:"isThinClone"`
	ParentSnap             ParentSnap           `json:"parentSnap,omitempty"`
	TieringPolicy          int                  `json:"tieringPolicy,omitempty"`
	ParentVolume           StorageResource      `json:"originalParentLun,omitempty"`
}

//Parent Snapshot to capture Source Snapshot Id
type ParentSnap struct {
	Id string `json:"id"`
}

//Pool struct to capture Pool Id
type Pool struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

//HostAccessResponse Struct to capture Host Access in Volume response
type HostAccessResponse struct {
	HostContent HostContent `json:"host"`
	HLU         int         `json:"hlu"`
}

//Link Struct to capture the link response
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

//BasicSystemInfo Struct to capture the BasicSystemInfo response
type BasicSystemInfo struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Links   []Link    `json:"links"`
	Entries []Entries `json:"entries"`
}

//Entries Struct to capture Entries contains a list of Links.
type Entries struct {
	Base    string    `json:"@base"`
	Updated time.Time `json:"updated"`
	Links   []Link    `json:"links"`
	Content Content   `json:"content"`
}

//Content Struct to capture the SystemInfo Content.
type Content struct {
	ID                 string `json:"id"`
	Model              string `json:"model"`
	Name               string `json:"name"`
	SoftwareVersion    string `json:"softwareVersion"`
	APIVersion         string `json:"apiVersion"`
	EarliestAPIVersion string `json:"earliestApiVersion"`
}

//Host struct to capture host object
type Host struct {
	HostContent HostContent `json:"content"`
}

//HostContent struct to capture host parameters
type HostContent struct {
	ID              string       `json:"id"`
	Name            string       `json:"name,omitempty"`
	Description     string       `json:"description,omitempty"`
	FcInitiators    []Initiators `json:"fcHostInitiators,omitempty"`
	IscsiInitiators []Initiators `json:"iscsiHostInitiators,omitempty"`
	IpPorts         []IpPorts    `json:"hostIPPorts,omitempty"`
	Address         string       `json:"address,omitempty"`
}

//FcInitiators struct to capture Initiator Id
type Initiators struct {
	Id string `json:"id"`
}

//IpPorts struct to capture IpPort Id
type IpPorts struct {
	Id      string `json:"id"`
	Address string `json:"address,omitempty"`
}

//HostIpPort struct to capture Host IP port object
type HostIpPort struct {
	HostIpContent HostContent `json:"content"`
}

//ListHostInitiator struct to capture host initiators
type ListHostInitiator struct {
	HostInitiator []HostInitiator `json:"entries"`
}

//HostInitiator struct to capture host initiator object
type HostInitiator struct {
	HostInitiatorContent HostInitiatorContent `json:"content"`
}

//HostInitiatorContent struct to capture host initiator parameters
type HostInitiatorContent struct {
	Id          string        `json:"id"`
	Health      HealthContent `json:"health"`
	Type        int           `json:"type"`
	InitiatorId string        `json:"InitiatorId"`
	IsIgnored   bool          `json:"isIgnored"`
	ParentHost  HostContent   `json:"parentHost"`
	Paths       []Path        `json:"paths"`
}

//Path struct to capture Path Id
type Path struct {
	Id string `json:"id"`
}

//Health struct to capture health status
type HealthContent struct {
	Value          int      `json:"value"`
	DescriptionIDs []string `json:"descriptionIds"`
	Descriptions   []string `json:"descriptions"`
}

//ListSnapshot struct to capture snapshot list
type ListSnapshot struct {
	Snapshots []Snapshot `json:"entries"`
}

//Snapshot struct to capture snapshot object
type Snapshot struct {
	SnapshotContent SnapshotContent `json:"content"`
}

//SnapshotContent struct to capture snapshot parameters
type SnapshotContent struct {
	ResourceId      string          `json:"id"`
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

//CopySnapshots struct to capture copy snapshot content
type CopySnapshots struct {
	CopySnapshotsContent CopySnapshotsContent `json:"content"`
}

//CopySnapshotsContent struct to capture copies list
type CopySnapshotsContent struct {
	Copies []StorageResource `json:"copies,omitempty"`
}

//StorageResource struct to capture storage resource Id
type StorageResource struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

//Struct to capture Storage Resource content
type StorageResourceParameters struct {
	StorageResourceContent StorageResourceContent `json:"content"`
}

type StorageResourceContent struct {
	Id         string          `json:"id"`
	Name       string          `json:"name,omitempty"`
	Filesystem StorageResource `json:"filesystem,omitempty"`
}

//IoLimitPolicy struct IO limit policy object
type IoLimitPolicy struct {
	IoLimitPolicyContent IoLimitPolicyContent `json:"content,omitempty"`
}

//IoLimitPolicyContent struct to capture IoLimitPolicyContent parameters
type IoLimitPolicyContent struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//Filesystem struct to capture filesystem object
type Filesystem struct {
	FileContent FileContent `json:"content"`
}

//FileContent struct to capture filesystem parameters
type FileContent struct {
	Id                     string  `json:"id"`
	Name                   string  `json:"name,omitempty"`
	SizeTotal              uint64  `json:"sizeTotal,omitempty"`
	Description            string  `json:"description,omitempty"`
	Type                   int     `json:"type,omitempty"`
	Format                 int     `json:"format,omitempty"`
	HostIOSize             int64   `json:"hostIOSize,omitempty"`
	TieringPolicy          uint64  `json:"tieringPolicy,omitempty"`
	IsThinEnabled          bool    `json:"isThinEnabled"`
	IsDataReductionEnabled bool    `json:"isDataReductionEnabled"`
	Pool                   Pool    `json:"pool,omitempty"`
	NASServer              Pool    `json:"nasServer,omitempty"`
	StorageResource        Pool    `json:"storageResource,omitempty"`
	NFSShare               []Share `json:"nfsShare,omitempty"`
	CIFSShare              []Pool  `json:"cifsShare,omitempty"`
}

//Share object to capture NFS Share object from FileContent
type Share struct {
	Id         string          `json:"id"`
	Name       string          `json:"name,omitempty"`
	Path       string          `json:"path,omitempty"`
	ParentSnap StorageResource `json:"snap,omitempty"`
}

//NFSShare struct to capture NFS Share object
type NFSShare struct {
	NFSShareContent NFSShareContent `json:"content"`
}

//NFSShareContent struct to capture NFS Share parameters
type NFSShareContent struct {
	Id                      string        `json:"id"`
	Name                    string        `json:"name,omitempty"`
	Filesystem              Pool          `json:"filesystem,omitempty"`
	ReadOnlyHosts           []HostContent `json:"readOnlyHosts,omitempty"`
	ReadWriteHosts          []HostContent `json:"readWriteHosts,omitempty"`
	ReadOnlyRootAccessHosts []HostContent `json:"readOnlyRootAccessHosts,omitempty"`
	RootAccessHosts         []HostContent `json:"rootAccessHosts,omitempty"`
	ExportPaths             []string      `json:"exportPaths,omitempty"`
}

//Struct to capture NAS Server object
type NASServer struct {
	NASServerContent NASServerContent `json:"content"`
}

type NASServerContent struct {
	Id        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	NFSServer NFSServer `json:"nfsServer,omitempty"`
}

type NFSServer struct {
	Id           string `json:"id"`
	Name         string `json:"name,omitempty"`
	NFSv3Enabled bool   `json:"nfsv3Enabled"`
	NFSv4Enabled bool   `json:"nfsv4Enabled"`
}

//ListIPInterfaces struct to capture snapshot list
type ListIPInterfaces struct {
	Entries []IPInterfaceEntries `json:"entries"`
}

//IPInterfaceEntries struct to capture IpInterface object
type IPInterfaceEntries struct {
	IPInterfaceContent IPInterfaceContent `json:"content"`
}

//IPInterfaceContent struct to capture IpInterface parameters
type IPInterfaceContent struct {
	ID        string `json:"id"`
	IPAddress string `json:"ipAddress"`
	Type      int    `json:"type"`
}

//LicenseInfo for features on Array
type LicenseInfo struct {
	LicenseInfoContent LicenseInfoContent `json:"content"`
}

//LicenseInfoContent for features on Array
type LicenseInfoContent struct {
	IsInstalled bool `json:"isInstalled"`
	IsValid     bool `json:"isValid"`
}

//HostInitiatorPath struct to capture host initiator path object
type HostInitiatorPath struct {
	HostInitiatorPathContent HostInitiatorPathContent `json:"content"`
}

//HostInitiatorPathContent struct to capture host initiator parameters
type HostInitiatorPathContent struct {
	FcPortID FcPortID `json:"fcPort"`
}

//FcPortID struct to capture FC port Id
type FcPortID struct {
	Id string `json:"id"`
}

//FcPort struct to capture FC port object
type FcPort struct {
	FcPortContent FcPortContent `json:"content"`
}

//FcPortContent struct to capture FC port Id
type FcPortContent struct {
	Wwn string `json:"wwn"`
}
