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
	Purchase(VLanName string) (*PurchasedVLAN, *Response, error)
	Delete(VLan *PurchasedVLAN) (*Response, error)
	Attach(attachRequest *PurchasedVLanAttachRequest) (*PurchasedVLAN, *Response, error)
	Detach(attachRequest *PurchasedVLanAttachRequest) (*Response, error)
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

type PurchasedVLanCreateRequest struct {
	VLanName string
}
type PurchasedVLanDeleteRequest struct {
	VLanResourceId int
}

type PurchasedVLanAttachRequest struct {
	NetworkAdapterId int
	VLanResourceId   int
	Gateway          string
	IP               string
	SubnetMask       string
}

type PrivateIp struct {
	Gateway    string `json:"GateWay"`
	IP         string
	SubNetMask string
}

type purchasedVLANsRoot struct {
	PurchasedVLans []PurchasedVLAN `json:"Value"`
}

type vLANRequestRoot struct {
	NetworkAdapterId    int
	VLanResourceId      int
	SetOnVirtualMachine bool
	PrivateIps          []PrivateIp `json:"PrivateIps,omitempty"`
}

// List all purchased VLans.
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

	return root.PurchasedVLans, resp, nil
}

func (s *VLANsServiceOp) Purchase(VLanName string) (*PurchasedVLAN, *Response, error) {
	req, err := s.client.NewRequest(vLANPurchasePath, PurchasedVLanCreateRequest{VLanName: VLanName})

	if err != nil {
		return nil, nil, err
	}

	root := new(PurchasedVLAN)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

func (s *VLANsServiceOp) Delete(VLan *PurchasedVLAN) (*Response, error) {
	req, err := s.client.NewRequest(vLANRemovePath, PurchasedVLanDeleteRequest{VLanResourceId: VLan.ResourceId})

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *VLANsServiceOp) Attach(attachRequest *PurchasedVLanAttachRequest) (*PurchasedVLAN, *Response, error) {
	if attachRequest == nil {
		return nil, nil, NewArgError("attachRequest", "cannot be nil")
	}

	vLANRequestRoot := vLANRequestRoot{
		NetworkAdapterId:    attachRequest.NetworkAdapterId,
		VLanResourceId:      attachRequest.VLanResourceId,
		SetOnVirtualMachine: false,
		PrivateIps:          []PrivateIp{},
	}
	if attachRequest.Gateway != "" {
		vLANRequestRoot.SetOnVirtualMachine = true
		vLANRequestRoot.PrivateIps = append(vLANRequestRoot.PrivateIps, PrivateIp{
			Gateway:    attachRequest.Gateway,
			IP:         attachRequest.IP,
			SubNetMask: attachRequest.SubnetMask,
		})
	}

	data := struct {
		VLanRequest interface{}
	}{vLANRequestRoot}

	req, err := s.client.NewRequest(vLANAttachPath, data)

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

func (s *VLANsServiceOp) Detach(attachRequest *PurchasedVLanAttachRequest) (*Response, error) {
	if attachRequest == nil {
		return nil, NewArgError("attachRequest", "cannot be nil")
	}

	vLANRequestRoot := vLANRequestRoot{
		NetworkAdapterId:    attachRequest.NetworkAdapterId,
		VLanResourceId:      attachRequest.VLanResourceId,
		SetOnVirtualMachine: false,
	}

	data := struct {
		VLanRequest interface{}
	}{vLANRequestRoot}

	req, err := s.client.NewRequest(vLANDetachPath, data)

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
