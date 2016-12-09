package goarubacloud

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
