package goarubacloud

const snapshotPath = "SetEnqueueServerSnapshot"

// SnapshotsService is an interface for interfacing with the cloud server actions
// endpoints of the Arubacloud API
type SnapshotsService interface {
	Create(int) (*Response, error)
	Restore(int) (*Response, error)
	Delete(int) (*Response, error)
}

// SnapshotsServiceOp handles communication with the cloud server action related
// methods of the Arubacloud API.
type SnapshotsServiceOp struct {
	client *Client
}

type snapshotRequest struct {
	ServerId               int
	SnapshotOperationTypes string
}

var _ SnapshotsService = &SnapshotsServiceOp{}

func (s *SnapshotsServiceOp) Create(serverId int) (*Response, error) {
	action := &snapshotRequest{ServerId: serverId, SnapshotOperationTypes: "Create"}
	return s.doAction(action)
}

func (s *SnapshotsServiceOp) Restore(serverId int) (*Response, error) {
	action := &snapshotRequest{ServerId: serverId, SnapshotOperationTypes: "Restore"}
	return s.doAction(action)
}

func (s *SnapshotsServiceOp) Delete(serverId int) (*Response, error) {
	action := &snapshotRequest{ServerId: serverId, SnapshotOperationTypes: "Delete"}
	return s.doAction(action)
}

func (s *SnapshotsServiceOp) doAction(snapshotRequest *snapshotRequest) (*Response, error) {
	data := struct {
		Snapshot interface{} `json:"Snapshot"`
	}{snapshotRequest}

	req, err := s.client.NewRequest(snapshotPath, data)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
