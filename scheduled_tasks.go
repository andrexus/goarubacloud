package goarubacloud

const getScheduledOperationsPath = "GetScheduledOperations"
const addScheduledOperationPath = "SetAddServerScheduledOperation"
const updateScheduledOperationPath = "SetUpdateServerScheduledOperation"
const removeScheduledOperationPath = "SetRemoveServerScheduledOperation"

// ScheduledTasksService is an interface for interfacing with the cloud server actions
// endpoints of the Arubacloud API
type ScheduledTasksService interface {
	List(*ScheduledTasksRequest) ([]ScheduledTask, *Response, error)
	Add() (*Response, error)
	Update(int) (*Response, error)
	Delete(int) (*Response, error)
}

// ScheduledTasksServiceOp handles communication with the cloud server action related
// methods of the Arubacloud API.
type ScheduledTasksServiceOp struct {
	client *Client
}

var _ ScheduledTasksService = &ScheduledTasksServiceOp{}

type ScheduledTask struct {
	ServerId             int
	ServerName           string
	OperationType        ScheduledTaskType
	OperationParameter   []interface{}
	ScheduledOperationID int
	ScheduledPlan        ScheduledPlan
}

type ScheduledPlan struct {
	FirstExecutionTime        string
	LastExecutionTime         string
	ScheduleDaysOfMonth       interface{}
	ScheduleEndDateTime       interface{}
	ScheduleFrequency         interface{}
	ScheduleFrequencyType     int
	ScheduleOperationLabel    string
	ScheduleStartDateTime     string
	ScheduleWeekDays          []interface{}
	ScheduledMontlyRecurrence interface{}
	ScheduledOwnerType        int
	ScheduledPlanId           int
	ScheduledPlanStatus       int
}

type ScheduledTaskType int

const (
	SWITCH_ON        ScheduledTaskType = 2
	FORCE_SHUTDOWN   ScheduledTaskType = 3
	SWITCH_OFF       ScheduledTaskType = 25
	CREATE_SNAPSHOT  ScheduledTaskType = 27
	RESTORE_SNAPSHOT ScheduledTaskType = 28
	DELETE_SNAPSHOT  ScheduledTaskType = 29
)

// String returns the name of the Datacenter.
func (m ScheduledTaskType) String() string {
	scheduled_task_types := map[ScheduledTaskType]string{
		2:  "Switch On",
		3:  "Force Shutdown",
		25: "Switch Off",
		27: "Create snapshot",
		28: "Restore Snapshot",
		29: "Delete snapshot",
	}

	return scheduled_task_types[m]
}

type ScheduledTasksRequest struct {
	Period *Period `json:"GetScheduledOperations"`
}

type Period struct {
	StartDate string
	EndDate   string
}

type ScheduledTaskCreateRequest struct {
	ScheduledOperationTypes   string
	ScheduleOperationLabel    string
	ServerID                  int
	ScheduleStartDateTime     string
	ScheduleEndDateTime       string
	ScheduleFrequencyType     string
	ScheduledPlanStatus       string
	ScheduledMontlyRecurrence string
}

type ScheduledTaskUpdateRequest struct {
	SetUpdateServerScheduledOperation ScheduledTaskIdCreate `json:"SetUpdateServerScheduledOperation"`
}

type ScheduledTaskIdCreate struct {
	ScheduledOperationId int
}

func (s ScheduledTasksServiceOp) List(request *ScheduledTasksRequest) ([]ScheduledTask, *Response, error) {
	// Not implemented yet
	return nil, nil, nil
}

func (s ScheduledTasksServiceOp) Add() (*Response, error) {
	// Not implemented yet
	return nil, nil
}

func (s ScheduledTasksServiceOp) Update(scheduledOperationId int) (*Response, error) {
	// Not implemented yet
	return nil, nil
}

func (s ScheduledTasksServiceOp) Delete(scheduledOperationId int) (*Response, error) {
	// Not implemented yet
	return nil, nil
}

func (s *ScheduledTasksServiceOp) doAction(scheduledOperationPath string, scheduledOperationRequest interface{}) (*Response, error) {
	req, err := s.client.NewRequest(scheduledOperationPath, scheduledOperationRequest)

	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
