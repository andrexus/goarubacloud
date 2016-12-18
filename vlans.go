package goarubacloud

const vLANsListPath = "GetPurchasedVLans"
const vLANPurchasePath = "SetPurchaseVLan"
const vLANRemovePath = "SetRemoveVLan"
const vLANAttachPath = "SetEnqueueAssociateVLan"
const vLANDetachPath = "SetEnqueueDeassociateVLan"

// VLANsService is an interface for interfacing with the purchased VLANs
// endpoint of the Arubacloud API
type VLANsService interface {
	List() ([]PurchasedVLAN, *Response, error)
	Purchase(name string) (*PurchasedVLAN, *Response, error)
	Delete(vlan_resource_id int) (*Response, error)
	Attach(attachRequest *purchasedVLANAttachRequest) (*PurchasedVLAN, *Response, error)
	Detach(network_adapter_id int, vlan_resource_id int) (*Response, error)
}

// VLANsServiceOp handles communication with the purchased VLANs related methods of the
// Arubacloud API.
type VLANsServiceOp struct {
	client *Client
}

var _ VLANsService = &VLANsServiceOp{}

type PurchasedVLAN struct {
	Name         string
	VlanCode     string
	ServerIds    []int
	ResourceId   int
	ResourceType int
	ProductId    int
	CompanyId    int
	UserId       int
}

type purchasedVLANAttachRequest struct {
	NetworkAdapterId int
	VLanResourceId   int
	Gateway          string
	IP               string
	SubnetMask       string
}

type privateIP struct {
	Gateway    string `json:"GateWay"`
	IP         string
	SubNetMask string
}

type purchasedVLANRoot struct {
	PurchasedVLAN *PurchasedVLAN `json:"Value"`
}

type purchasedVLANsRoot struct {
	PurchasedVLANs []PurchasedVLAN `json:"Value"`
}

type vLANRequestRoot struct {
	NetworkAdapterId    int
	VLanResourceId      int
	SetOnVirtualMachine bool
	PrivateIps          []privateIP `json:"PrivateIps,omitempty"`
}

func NewPurchasedVLanAttachRequest(network_adapter_id int, vlan_resource_id int) *purchasedVLANAttachRequest {
	return &purchasedVLANAttachRequest{
		NetworkAdapterId: network_adapter_id,
		VLanResourceId:   vlan_resource_id,
	}
}

// List all purchased VLANs.
func (s *VLANsServiceOp) List() ([]PurchasedVLAN, *Response, error) {
	req, err := s.client.NewRequest(vLANsListPath, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(purchasedVLANsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.PurchasedVLANs, resp, nil
}

func (s *VLANsServiceOp) Purchase(name string) (*PurchasedVLAN, *Response, error) {
	body := struct{ VLanName string }{VLanName: name}
	req, err := s.client.NewRequest(vLANPurchasePath, body)

	if err != nil {
		return nil, nil, err
	}

	root := new(purchasedVLANRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.PurchasedVLAN, resp, nil
}

func (s *VLANsServiceOp) Delete(vlan_resource_id int) (*Response, error) {
	body := struct{ VLanResourceId int }{VLanResourceId: vlan_resource_id}
	req, err := s.client.NewRequest(vLANRemovePath, body)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *VLANsServiceOp) Attach(attachRequest *purchasedVLANAttachRequest) (*PurchasedVLAN, *Response, error) {
	if attachRequest == nil {
		return nil, nil, NewArgError("attachRequest", "cannot be nil")
	}

	vLANRequestRoot := vLANRequestRoot{
		NetworkAdapterId:    attachRequest.NetworkAdapterId,
		VLanResourceId:      attachRequest.VLanResourceId,
		SetOnVirtualMachine: false,
		PrivateIps:          []privateIP{},
	}
	if attachRequest.Gateway != "" {
		vLANRequestRoot.SetOnVirtualMachine = true
		vLANRequestRoot.PrivateIps = append(vLANRequestRoot.PrivateIps, privateIP{
			Gateway:    attachRequest.Gateway,
			IP:         attachRequest.IP,
			SubNetMask: attachRequest.SubnetMask,
		})
	}

	body := struct {
		VLanRequest interface{}
	}{vLANRequestRoot}

	req, err := s.client.NewRequest(vLANAttachPath, body)

	if err != nil {
		return nil, nil, err
	}

	root := new(PurchasedVLAN)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *VLANsServiceOp) Detach(network_adapter_id int, vlan_resource_id int) (*Response, error) {
	if network_adapter_id == 0 {
		return nil, NewArgError("network_adapter_id", "cannot be nil")
	}

	if vlan_resource_id == 0 {
		return nil, NewArgError("vlan_resource_id", "cannot be nil")
	}

	vLANRequestRoot := vLANRequestRoot{
		NetworkAdapterId:    network_adapter_id,
		VLanResourceId:      vlan_resource_id,
		SetOnVirtualMachine: false,
	}

	body := struct {
		VLanRequest interface{}
	}{vLANRequestRoot}

	req, err := s.client.NewRequest(vLANDetachPath, body)

	if err != nil {
		return nil, err
	}

	root := new(PurchasedVLAN)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return resp, err
	}

	return resp, err
}
