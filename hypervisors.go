package goarubacloud

import (
	"fmt"
	"log"
)

const hypervisorsPath = "GetHypervisors"

// HypervisorsService is an interface for interfacing with the Hypervisors
// endpoints of the Arubacloud API
type HypervisorsService interface {
	GetHypervisors() ([]Hypervisor, *Response, error)
	FindOsTemplate(HypervisorType, string) (*OSTemplate, error)
}

// HypervisorsServiceOp handles communication with the Hypervisor related methods of the
// Arubacloud API.
type HypervisorsServiceOp struct {
	client *Client
}

var _ HypervisorsService = &HypervisorsServiceOp{}

func (s *HypervisorsServiceOp) GetHypervisors() ([]Hypervisor, *Response, error) {

	req, err := s.client.NewRequest(hypervisorsPath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(hypervisorsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Hypervisors, resp, err
}

func (s *HypervisorsServiceOp)FindOsTemplate(hypervisorType HypervisorType, name string)(*OSTemplate, error)  {
	hypervisors, _, err := s.client.Hypervisors.GetHypervisors()
	if err != nil {
		log.Println("[ERROR] Unable to fetch hypervisors: ", err)
		return nil, err
	}

	var hypervisor *Hypervisor

	for _, hypervisorItem := range hypervisors {
		if hypervisorItem.HypervisorType == hypervisorType {
			hypervisor = &hypervisorItem
			break
		}
	}

	if hypervisor == nil {
		return nil, fmt.Errorf("Hypervisor not found for type: %s", hypervisorType)
	}

	for _, template := range hypervisor.Templates {
		if template.Description == name {
			return &template, nil
		}
	}
	return nil, fmt.Errorf("Template not found for %s", name)
}

type HypervisorType int

const (
	Microsoft_Hyper_V HypervisorType = 1 + iota
	VMWare_Cloud_Pro
	Microsoft_Hyper_V_Low_Cost
	VMWare_Cloud_Smart
)

var hypervisors = [...]string{
	"Microsoft Hyper-V (Cloud Pro)",
	"VMWare (Cloud Pro)",
	"Microsoft Hyper-V Low Cost (Cloud Pro)",
	"VMWare (Cloud Smart)",
}

// String returns the name of the Hypervisor.
func (m HypervisorType) String() string {
	return hypervisors[m-1]
}

type Hypervisor struct {
	HypervisorServerType int
	HypervisorType       HypervisorType
	Templates            []OSTemplate `json:"Templates"`
}

type OSTemplate struct {
	ApplianceType                   int
	ArchitectureType                int
	CompanyID                       int
	CompatiblePreConfiguredPackages []interface{}
	Description                     string
	Enabled                         bool
	ExportEnabled                   bool
	FeatureTypes                    []interface{}
	Icon                            interface{}
	Id                              int
	IdentificationCode              string
	Ipv6Compatible                  bool
	Name                            string
	OSFamily                        int
	OSVersion                       string
	OwnerUserId                     string
	ParentTemplateID                int
	ProductId                       int
	ResourceBounds                  []ResourceBounds
	Revision                        string
	SshKeyInitializationSupported   bool
	TemplateExtendedDescription     string
	TemplateOwnershipType           int
	TemplatePassword                string
	TemplateSellingStatus           int
	TemplateStatus                  interface{}
	TemplateType                    int
	TemplateUsername                string
	ToolsAvailable                  bool
}

type ResourceBounds struct {
	ResourceType int
	Default      int
	Min          int
	Max          int
}

type OSTemplateDetails struct {
	Id           int
	Name         string
	Description  string
	ResourceType int
	ResourceId   int
	CompanyId    int
	ProductId    int
	UserId       int
}
