package goarubacloud

import (
	"fmt"
	"time"
)

const cloudSeverListPath = "GetServers"
const cloudSeverDetailsPath = "GetServerDetails"
const cloudSeverCreatePath = "SetEnqueueServerCreation"
const cloudSeverDeletePath = "SetEnqueueServerDeletion"

// CloudServersService is an interface for interfacing with the Cloud Server
// endpoints of the Arubacloud API
type CloudServersService interface {
	List() ([]CloudServer, *Response, error)
	Get(int) (*CloudServerDetails, *Response, error)
	Create(CloudServerCreator) (*CloudServer, *Response, error)
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
	PublicIpAddressResourceId int
}

type CloudServerCreator interface {
	GetServerName() string
	GetRequest() interface{}
}

type CloudServerProCreator interface {
	AddVirtualDisk(int) error
	AddPublicIp(int) error
	SetCPUQuantity(int) error
	SetRAMQuantity(int) error
	SetNote(string) error
	GetServerName() string
	GetRequest() interface{}
}

func NewCloudServerProCreateRequest(name string, admin_password string, os_template_id int) CloudServerProCreator {
	createRequest := &cloudServerCreateRequestPro{
		Name:                         name,
		OSTemplateId:                 os_template_id,
		AdministratorPassword:        admin_password,
		CPUQuantity:                  1,
		RAMQuantity:                  1,
		Note:                         "",
		VirtualDisks:                 []CloudServerCreateVirtualDisk{},
		NetworkAdaptersConfiguration: []NetworkAdapterCreateConfiguration{},
	}
	return createRequest
}

func NewCloudServerSmartCreateRequest(smart_server_type CloudServerSmartType, name string, admin_password string, os_template_id int) CloudServerCreator {
	return &cloudServerCreateRequestSmart{
		Name:         name,
		OSTemplateId: os_template_id,
		Note:         "",
		AdministratorPassword: admin_password,
		CloudServerSmartType:  smart_server_type,
	}
}

// cloudServerCreateRequestPro represents a request to create a CloudServer PRO.
type cloudServerCreateRequestPro struct {
	Name                         string
	AdministratorPassword        string
	OSTemplateId                 int
	Note                         string
	CPUQuantity                  int
	RAMQuantity                  int
	VirtualDisks                 []CloudServerCreateVirtualDisk
	NetworkAdaptersConfiguration []NetworkAdapterCreateConfiguration
}

func (r *cloudServerCreateRequestPro) AddVirtualDisk(size int) error {
	if len(r.VirtualDisks) == 4 {
		return NewArgError("operation", "max disk count is 4")
	}
	if size < 10 || size > 500 {
		return NewArgError("size", "MaxSize per Disk: 500 GB. MinSize per Disk 10 GB")
	}
	if size%10 != 0 {
		return NewArgError("size", "disk size must be a multiple of 10")

	}
	r.VirtualDisks = append(r.VirtualDisks, CloudServerCreateVirtualDisk{
		VirtualDiskType: len(r.VirtualDisks),
		Size:            size,
	})
	return nil
}

func (r *cloudServerCreateRequestPro) AddPublicIp(resourceId int) error {
	if resourceId == 0 {
		return NewArgError("resourceId", "it must be > 0")
	}
	publicIpAddress := PublicIpAddress{PrimaryIPAddress: "true", PublicIpAddressResourceId: resourceId}

	r.NetworkAdaptersConfiguration = append(r.NetworkAdaptersConfiguration, NetworkAdapterCreateConfiguration{
		NetworkAdapterType: len(r.NetworkAdaptersConfiguration),
		PublicIpAddresses:  []PublicIpAddress{publicIpAddress},
	})

	return nil
}

func (r *cloudServerCreateRequestPro) SetCPUQuantity(cpu_quantity int) error {
	if cpu_quantity < 1 {
		return NewArgError("cpu_quantity", "it must be > 1")
	}

	r.CPUQuantity = cpu_quantity
	return nil
}

func (r *cloudServerCreateRequestPro) SetRAMQuantity(ram_quantity int) error {
	if ram_quantity < 1 {
		return NewArgError("ram_quantity", "it must be > 1")
	}

	r.RAMQuantity = ram_quantity
	return nil
}

func (r *cloudServerCreateRequestPro) SetNote(note string) error {
	if len(note) > 4096 {
		return NewArgError("note", "it is too long")
	}
	r.Note = note
	return nil
}

func (r *cloudServerCreateRequestPro) GetServerName() string {
	return r.Name
}

func (r *cloudServerCreateRequestPro) GetRequest() interface{} {
	if len(r.VirtualDisks) == 0 {
		r.VirtualDisks = append(r.VirtualDisks, CloudServerCreateVirtualDisk{
			VirtualDiskType: 0,
			Size:            10,
		})
	}

	return r
}

// cloudServerCreateRequestSmart represents a request to create a CloudServer SMART.
type cloudServerCreateRequestSmart struct {
	Name                  string               `json:"Name"`
	AdministratorPassword string               `json:"AdministratorPassword"`
	OSTemplateId          int                  `json:"OSTemplateId"`
	CloudServerSmartType  CloudServerSmartType `json:"SmartVMWarePackageID"`
	Note                  string               `json:"Note"`
}

func (r *cloudServerCreateRequestSmart) GetServerName() string {
	return r.Name
}

func (r *cloudServerCreateRequestSmart) GetRequest() interface{} {
	return r
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

type NetworkAdapterCreateConfiguration struct {
	NetworkAdapterType int
	PublicIpAddresses  []PublicIpAddress
}

func (s *CloudServerDetails) GetPublicIpAddress() (string, error) {
	if s.EasyCloudPackageID != 0 {
		return s.EasyCloudIPAddress.Value, nil
	} else {
		for _, network_adapter := range s.NetworkAdapters {
			for _, public_ip := range network_adapter.IPAddresses {
				return public_ip.Value, nil
			}
		}
		return "", fmt.Errorf("Server %d doesn't have any public IPs assigned to it", s.ServerId)
	}
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
func (s *CloudServersServiceOp) Create(requestCreator CloudServerCreator) (*CloudServer, *Response, error) {
	if requestCreator == nil {
		return nil, nil, NewArgError("requestCreator", "cannot be nil")
	}

	if requestCreator.GetRequest() == nil {
		return nil, nil, NewArgError("request", "cannot be nil")
	}

	data := struct {
		Server interface{} `json:"Server"`
	}{requestCreator.GetRequest()}

	req, err := s.client.NewRequest(cloudSeverCreatePath, data)

	if err != nil {
		return nil, nil, err
	}

	root := new(CloudServer)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	server, err := WaitForServerWithName(s.client, requestCreator.GetServerName())
	if err != nil {
		return nil, nil, err
	}

	return server, resp, err
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

	err = WaitForServerStatus(s.client, serverId, OFF)
	if err != nil {
		return nil, err
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
	// maxRetries is the amount of times we can fail before deciding
	// the check is a total failure.
	maxRetries = 10
)

// WaitForServerStatus waits for a cloud servers status
func WaitForServerStatus(client *Client, serverId int, status ServerStatus) error {
	completed := false
	failCount := 0
	for !completed {
		server_details, _, err := client.CloudServers.Get(serverId)

		if err != nil {
			if failCount <= maxRetries {
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

// WaitForServerWithName waits for a server with specified name appears in the list
func WaitForServerWithName(client *Client, serverName string) (*CloudServer, error) {
	completed := false
	failCount := 0
	var server *CloudServer
	for !completed {
		servers, _, err := client.CloudServers.List()

		if err != nil {
			if failCount <= maxRetries {
				failCount++
				continue
			}
			return nil, err
		}

		for _, serverItem := range servers {
			if serverItem.Name == serverName {
				completed = true
				server = &serverItem
			}
		}

		time.Sleep(5 * time.Second)
	}

	return server, nil
}
