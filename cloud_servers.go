package goarubacloud

import "time"

const cloudSeverListPath = "GetServers"
const cloudSeverDetailsPath = "GetServerDetails"
const cloudSeverCreatePath = "SetEnqueueServerCreation"
const cloudSeverDeletePath = "SetEnqueueServerDeletion"

// CloudServersService is an interface for interfacing with the Cloud Server
// endpoints of the Arubacloud API
type CloudServersService interface {
	List() ([]CloudServer, *Response, error)
	Get(int) (*CloudServerDetails, *Response, error)
	Create(interface{}) (*CloudServer, *Response, error)
	Delete(int) (*Response, error)
}

// CloudServersServiceOp handles communication with the Cloud Server related methods of the
// Arubacloud API.
type CloudServersServiceOp struct {
	client *Client
}

var _ CloudServersService = &CloudServersServiceOp{}

// CloudServer represents a Arubacloud Cloud Server
type CloudServer struct {
	Busy                 bool
	CPUQuantity          int
	CompanyId            int
	DatacenterId         DataCenterRegion
	HDQuantity           int
	HDTotalSize          int
	HypervisorServerType int
	HypervisorType       HypervisorType
	Name                 string
	OSTemplateId         int
	RAMQuantity          int
	ServerId             int
	ServerStatus         ServerStatus
	UserId               int
}

type CloudServerDetails struct {
	ActiveJobs                []ActiveJob
	CPUQuantity               CPUQuantity
	CompanyId                 int
	ControlToolActivationDate string
	ControlToolInstalled      bool
	CreationDate              string
	DatacenterId              DataCenterRegion
	EasyCloudIPAddress        EasyCloudIPAddress
	EasyCloudPackageID        int
	HypervisorServerType      int
	HypervisorType            HypervisorType
	Name                      string
	NetworkAdapters           []NetworkAdapter
	Note                      string
	OSTemplate                OSTemplateDetails
	Parameters                []interface{}
	RAMQuantity               RAMQuantity
	RenewDateSmart            string
	ScheduledOperations       []ScheduledTask
	ServerId                  int
	ServerStatus              ServerStatus
	Snapshots                 []interface{}
	ToolsAvailable            bool
	UserId                    int
	VirtualDVDs               []interface{}
	VirtualDisks              []VirtualDisk
	VncPort                   int
}

type ActiveJob struct {
	JobId          int
	Status         int
	OperationName  string
	Progress       int
	ServerId       int
	ServerName     string
	CreationDate   string
	LastUpdateDate string
	LicenseId      interface{}
	ResourceId     int
	ResourceValue  interface{}
	UserId         int
	Username       string
}

type EasyCloudIPAddress struct {
	Value            string
	SubNetMask       string
	Gateway          string
	GatewayIPv6      string
	PrefixIPv6       int
	SubnetPrefixIPv6 string
	StartRangeIPv6   string
	EndRangeIPv6     string
	ServerId         int
	CompanyId        int
	ProductId        int
	ResourceId       int
	ResourceType     int
	UserId           int
	LoadBalancerID   interface{}
}

type CPUQuantity struct {
	CompanyId    int
	ProductId    int
	ResourceId1  int
	ResourceType int
	UserId       int
	Quantity     int
}

type RAMQuantity struct {
	CompanyId    int
	ProductId    int
	ResourceId   int
	ResourceType int
	UserId       int
	Quantity     int
}

type VirtualDisk struct {
	CompanyId    int
	ProductId    int
	ResourceId   int
	ResourceType int
	UserId       int
	CreationDate string
	Size         int
}

type ServerStatus int

const (
	CREATION_IN_PROGRESS ServerStatus = 1 + iota
	OFF
	ON
)

var server_status = [...]string{
	"CREATION IN PROGRESS",
	"OFF",
	"ON",
}

// String returns the name of the Datacenter.
func (m ServerStatus) String() string {
	return server_status[m-1]
}

type CloudServerSmartType int

const (
	SMALL CloudServerSmartType = 1 + iota
	MEDIUM
	LARGE
	EXTRALARGE
)

type NetworkAdapter struct {
	Id                 int
	NetworkAdapterType int
	IPAddresses        []IpAddress
	PublicIpAddresses  []PublicIpAddress
	MacAddress         string
	ServerId           int
	VLan               PurchasedVLAN
}

type IpAddress struct {
	Value            string
	Gateway          string
	SubNetMask       string
	GatewayIPv6      string
	PrefixIPv6       int
	SubnetPrefixIPv6 string
	StartRangeIPv6   string
	EndRangeIPv6     string
	ServerId         int
	CompanyId        int
	ProductId        int
	ResourceId       int
	ResourceType     int
	UserId           int
	LoadBalancerID   interface{}
}

type PublicIpAddress struct {
	PrimaryIPAddress          string
	PublicIpAddressResourceId string
}

// CloudServerCreateRequestPro represents a request to create a CloudServer PRO.
type CloudServerCreateRequestPro struct {
	Name                         string                         `json:"Name"`
	AdministratorPassword        string                         `json:"AdministratorPassword"`
	OSTemplateId                 int                            `json:"OSTemplateId"`
	Note                         string                         `json:"Note"`
	CPUQuantity                  int                            `json:"CPUQuantity"`
	RAMQuantity                  int                            `json:"RAMQuantity"`
	VirtualDisks                 []CloudServerCreateVirtualDisk `json:"VirtualDisks"`
	NetworkAdaptersConfiguration []NetworkAdapter               `json:"NetworkAdaptersConfiguration"`
}

// CloudServerCreateRequestSmart represents a request to create a CloudServer SMART.
type CloudServerCreateRequestSmart struct {
	Name                  string               `json:"Name"`
	AdministratorPassword string               `json:"AdministratorPassword"`
	OSTemplateId          int                  `json:"OSTemplateId"`
	CloudServerSmartType  CloudServerSmartType `json:"SmartVMWarePackageID"`
	Note                  string               `json:"Note"`
}

type hypervisorsRoot struct {
	Hypervisors []Hypervisor `json:"Value"`
}

type cloudServerDetailsRoot struct {
	CloudServerDetails *CloudServerDetails `json:"Value"`
}

type cloudServersRoot struct {
	CloudServers []CloudServer `json:"Value"`
}

// CloudServerCreateVirtualDisk identifies a VirtualDisk to attach for the create request.
type CloudServerCreateVirtualDisk struct {
	VirtualDiskType int
	Size            int
}

// List all CloudServers
func (s *CloudServersServiceOp) List() ([]CloudServer, *Response, error) {
	req, err := s.client.NewRequest(cloudSeverListPath, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(cloudServersRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.CloudServers, resp, err
}

// Get individual CloudServer
func (s *CloudServersServiceOp) Get(serverId int) (*CloudServerDetails, *Response, error) {
	req, err := s.client.NewRequest(cloudSeverDetailsPath, ServerIdCreate{ServerId: serverId})

	if err != nil {
		return nil, nil, err
	}

	root := new(cloudServerDetailsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.CloudServerDetails, resp, err
}

// Create cloudServer
func (s *CloudServersServiceOp) Create(createRequest interface{}) (*CloudServer, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	data := struct {
		Server interface{} `json:"Server"`
	}{createRequest}

	req, err := s.client.NewRequest(cloudSeverCreatePath, data)

	if err != nil {
		return nil, nil, err
	}

	root := new(CloudServer)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete CloudServer
func (s *CloudServersServiceOp) Delete(serverId int) (*Response, error) {
	serverDetails, resp, err := s.client.CloudServers.Get(serverId)
	if err != nil {
		return resp, err
	}

	if serverDetails.ServerStatus == ON {
		resp, err := s.client.CloudServerActions.PowerOff(serverId)
		if err != nil {
			return resp, err
		}
	}

	req, err := s.client.NewRequest(cloudSeverDeletePath, ServerIdCreate{ServerId: serverId})

	if err != nil {
		return nil, err
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

const (
	// activeFailure is the amount of times we can fail before deciding
	// the check for active is a total failure. This can help account
	// for servers randomly not answering.
	activeFailure = 4
)

// ServerStatus waits for a cloud servers status
func WaitForServerStatus(client *Client, serverId int, status ServerStatus) error {
	completed := false
	failCount := 0
	for !completed {
		server_details, _, err := client.CloudServers.Get(serverId)

		if err != nil {
			if failCount <= activeFailure {
				failCount++
				continue
			}
			return err
		}

		if server_details.ServerStatus != status {
			time.Sleep(10 * time.Second)
		} else {
			completed = true
		}
	}

	return nil
}
