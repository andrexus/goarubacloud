package goarubacloud

const purchasedIpsListPath = "GetPurchasedIpAddresses"
const purchaseIpPath = "SetPurchaseIpAddress"
const removeIpPath = "SetRemoveIpAddress"

type PurchasedIPsService interface {
	List() ([]PurchasedIP, *Response, error)
	Purchase() (*PurchasedIP, *Response, error)
	Delete(int) (*Response, error)
}

// PurchasedIPsServiceOp handles communication with the purchased IPs related methods of the
// Arubacloud API.
type PurchasedIPsServiceOp struct {
	client *Client
}

var _ PurchasedIPsService = &PurchasedIPsServiceOp{}

// PurchasedIP represents an Arubacloud purchased IP.
type PurchasedIP struct {
	Value            string
	SubNetMask       string
	Gateway          string
	GatewayIPv6      string
	PrefixIPv6       int
	SubnetPrefixIPv6 string
	StartRangeIPv6   string
	EndRangeIPv6     string
	ServerId         int
	UserId           int
	CompanyId        int
	ProductId        int
	ResourceId       int
	ResourceType     int
	LoadBalancerID   interface{}
}

type purchasedIpRoot struct {
	PurchasedIP *PurchasedIP `json:"Value"`
}

type purchasedIpsRoot struct {
	PurchasedIps []PurchasedIP `json:"Value"`
}

type PurchasedIpRemoveRequest struct {
	IpAddressResourceId int
}

// List all purchased IPs.
func (s *PurchasedIPsServiceOp) List() ([]PurchasedIP, *Response, error) {
	req, err := s.client.NewRequest(purchasedIpsListPath, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(purchasedIpsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.PurchasedIps, resp, nil
}

// Purchase a new IP.
func (s *PurchasedIPsServiceOp) Purchase() (*PurchasedIP, *Response, error) {
	req, err := s.client.NewRequest(purchaseIpPath, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(purchasedIpRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.PurchasedIP, resp, nil
}

// Delete purchased IP.
func (s *PurchasedIPsServiceOp) Delete(IpAddressResourceId int) (*Response, error) {
	req, err := s.client.NewRequest(removeIpPath,
		PurchasedIpRemoveRequest{IpAddressResourceId:IpAddressResourceId})

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
