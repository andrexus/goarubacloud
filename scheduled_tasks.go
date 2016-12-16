package goarubacloud

import (
	"fmt"
	"time"
)

const getScheduledOperationsPath = "GetScheduledOperations"
const addScheduledOperationPath = "SetAddServerScheduledOperation"
const updateScheduledOperationPath = "SetUpdateServerScheduledOperation"
const removeScheduledOperationPath = "SetRemoveServerScheduledOperation"

// ScheduledTasksService is an interface for interfacing with the cloud server actions
// endpoints of the Arubacloud API
type ScheduledTasksService interface {
	List(*Interval) ([]ScheduledTask, *Response, error)
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

type scheduledTasksRoot struct {
	ScheduledTasks []ScheduledTask `json:"Value"`
}

type Interval struct {
	StartDate time.Time
	EndDate   time.Time
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

func (s ScheduledTasksServiceOp) List(interval *Interval) ([]ScheduledTask, *Response, error) {
	data := struct {
		StartDate string
		EndDate   string
	}{
		StartDate: fmt.Sprintf("/Date(%d)/", interval.StartDate.Unix()),
		EndDate:   fmt.Sprintf("/Date(%d)/", interval.EndDate.Unix()),
	}

	req, err := s.client.NewRequest(getScheduledOperationsPath, data)

	if err != nil {
		return nil, nil, err
	}

	root := new(scheduledTasksRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ScheduledTasks, resp, err
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
