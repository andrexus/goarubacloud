package goarubacloud

const cloudServerPowerOffPath = "SetEnqueueServerPowerOff"
const cloudServerPowerOnPath = "SetEnqueueServerStart"
const cloudServerArchivePath = "ArchiveVirtualServer"
const cloudServerRestorePath = "SetEnqueueServerRestore"
const cloudServerReinitializePath = "SetEnqueueReinitializeServer"

// CloudServerActionsService is an interface for interfacing with the Cloud Server actions
// endpoints of the Arubacloud API
type CloudServerActionsService interface {
	PowerOff(int) (*Response, error)
	PowerOn(int) (*Response, error)
	PowerCycle(int) (*Response, error)
	Archive(int) (*Response, error)
	Restore(serverId int, CPUQuantity int, RAMQuantity int) (*Response, error)
	Reinitialize(*ServerReinitializeRequest) (*Response, error)
}

// CloudServerActionsServiceOp handles communication with the Cloud Server action related
// methods of the Arubacloud API.
type CloudServerActionsServiceOp struct {
	client *Client
}

var _ CloudServerActionsService = &CloudServerActionsServiceOp{}

type ServerIdCreate struct {
	ServerId    int
	CPUQuantity int `json:"CPUQuantity,omitempty"`
	RAMQuantity int `json:"RAMQuantity,omitempty"`
}

type ServerReinitializeRequest struct {
	ServerId              int
	AdministratorPassword string `json:"AdministratorPassword,omitempty"`
	OSTemplateID          int    `json:"OSTemplateID,omitempty"`
	ConfigureIPv6         bool   `json:"ConfigureIPv6,omitempty"`
}

type cloudServerActionRequest struct {
	ActionPath            string
	ServerIdCreateRequest *ServerIdCreate
}

// PowerOff a Cloud Server
func (s *CloudServerActionsServiceOp) PowerOff(serverId int) (*Response, error) {
	action := &cloudServerActionRequest{ActionPath: cloudServerPowerOffPath,
		ServerIdCreateRequest: &ServerIdCreate{ServerId: serverId}}
	return s.doAction(action)
}

// PowerOn a Cloud Server
func (s *CloudServerActionsServiceOp) PowerOn(serverId int) (*Response, error) {
	action := &cloudServerActionRequest{ActionPath: cloudServerPowerOnPath,
		ServerIdCreateRequest: &ServerIdCreate{ServerId: serverId}}
	return s.doAction(action)
}

// PowerCycle a Cloud Server
func (s *CloudServerActionsServiceOp) PowerCycle(serverId int) (*Response, error) {
	serverDetails, resp, err := s.client.CloudServers.Get(serverId)
	if err != nil {
		return resp, err
	}

	if serverDetails.ServerStatus == OFF {
		resp, err = s.PowerOn(serverId)
		if err != nil {
			return nil, err
		}
	}

	resp, err = s.PowerOff(serverId)
	if err != nil {
		return resp, err
	}
	err = WaitForServerStatus(s.client, serverId, OFF)
	if err != nil {
		return nil, err
	}

	return s.PowerOn(serverId)
}

// Archive Cloud Server
func (s *CloudServerActionsServiceOp) Archive(serverId int) (*Response, error) {
	data := struct {
		ArchiveVirtualServer interface{} `json:"ArchiveVirtualServer"`
	}{ServerIdCreate{ServerId: serverId}}

	req, err := s.client.NewRequest(cloudServerArchivePath, data)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Restore Cloud Server
func (s *CloudServerActionsServiceOp) Restore(serverId int, qpu_quantity int, ram_quantity int) (*Response, error) {
	data := struct {
		SetEnqueueServerRestore interface{} `json:"SetEnqueueServerRestore"`
	}{ServerIdCreate{
		ServerId:    serverId,
		CPUQuantity: qpu_quantity,
		RAMQuantity: ram_quantity}}

	req, err := s.client.NewRequest(cloudServerRestorePath, data)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Restore Cloud Server
func (s *CloudServerActionsServiceOp) Reinitialize(serverReinitializeRequest *ServerReinitializeRequest) (*Response, error) {
	serverId := serverReinitializeRequest.ServerId
	serverDetails, resp, err := s.client.CloudServers.Get(serverId)
	if err != nil {
		return resp, err
	}

	if serverDetails.ServerStatus == ON {
		resp, err = s.PowerOff(serverId)
		if err != nil {
			return nil, err
		}
	}

	err = WaitForServerStatus(s.client, serverId, OFF)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(cloudServerReinitializePath, serverReinitializeRequest)

	if err != nil {
		return nil, err
	}

	resp, err = s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (s *CloudServerActionsServiceOp) doAction(actionRequest *cloudServerActionRequest) (*Response, error) {
	req, err := s.client.NewRequest(actionRequest.ActionPath, actionRequest.ServerIdCreateRequest)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
