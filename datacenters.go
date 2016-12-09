package goarubacloud

const virtualDatacenterPath = "GetVirtualDatacenter"
const activeJobsPath = "GetJobs"

type DataCentersService interface {
	GetVirtualDatacenter() (*VirtualDatacenter, *Response, error)
	GetJobs() ([]ActiveJob, *Response, error)
}

type DataCentersServiceOp struct {
	client *Client
}

var _ DataCentersService = &DataCentersServiceOp{}

type VirtualDatacenter struct {
	CustomProductEntities []interface{}
	DatacenterRegion      DataCenterRegion `json:"DatacenterId"`
	FTP                   interface{}
	IpAddresses           []interface{}
	LoadBalancers         []interface{}
	PleskLicenses         []interface{}
	PrivateCloudEntities  []interface{}
	PublicIpAddresses     []interface{}
	Servers               []CloudServerDetails
	SharedStorages        []interface{}
	VLans                 []interface{}
}

type virtualDataCenterRoot struct {
	VirtualDatacenter *VirtualDatacenter `json:"Value"`
}

type activeJobsRoot struct {
	ActiveJobs []ActiveJob `json:"Value"`
}

type DataCenterRegion int

const (
	Italy_1 DataCenterRegion = 1 + iota
	Italy_2
	Czech_Republic
	France
	Germany
	UK
)

var datacenter_regions = [...]string{
	"Italy 1",
	"Italy 2",
	"Czech Republic",
	"France",
	"Germany",
	"UK",
}

// String returns the name of the Datacenter.
func (m DataCenterRegion) String() string { return datacenter_regions[m-1] }

// Get info about used services in the datacenter
func (s *DataCentersServiceOp) GetVirtualDatacenter() (*VirtualDatacenter, *Response, error) {
	req, err := s.client.NewRequest(virtualDatacenterPath, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(virtualDataCenterRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.VirtualDatacenter, resp, err
}

// Get active jobs
func (s *DataCentersServiceOp) GetJobs() ([]ActiveJob, *Response, error) {
	req, err := s.client.NewRequest(activeJobsPath, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(activeJobsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ActiveJobs, resp, err
}
