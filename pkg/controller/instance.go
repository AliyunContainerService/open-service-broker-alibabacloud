package controller

import (
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"github.com/golang/glog"
)

type CreateServiceInstancePayload struct {
	InstanceID string
	Request    brokerapi.CreateServiceInstanceRequest
	Result     chan brokerapi.WorkerResponse
}

func (s CreateServiceInstancePayload) Deal() error {
	glog.Infof("Deal CreateServiceInstance with payload :%v.", s)
	response := brokerapi.WorkerResponse{}
	controller := GetBaseController()

	err := controller.CreateServiceInstance(s.InstanceID, s.Request.ServiceID, s.Request.PlanID, s.Request.Parameters)
	if err != nil {
		response.Message = ""
		response.Err = err

	} else {
		response.Message = brokerapi.StateInProgress
		response.Err = nil
	}
	s.Result <- response
	return err
}

type DeleteServiceInstancePayload struct {
	InstanceID string
	Request    brokerapi.DeleteServiceInstanceRequest
	Result     chan brokerapi.WorkerResponse
}

func (s DeleteServiceInstancePayload) Deal() error {
	response := brokerapi.WorkerResponse{}

	controller := GetBaseController()
	err := controller.DeleteServiceInstance(s.InstanceID, s.Request.ServiceID, s.Request.PlanID, nil)
	if err != nil {
		response.Message = ""
		response.Err = err
	} else {
		response.Message = "Delete instance success."
		response.Err = nil
	}
	s.Result <- response
	return err
}

type InstanceLastOperationPayload struct {
	InstanceID string
	Request    brokerapi.LastOperationRequest
	Result     chan brokerapi.WorkerResponse
}

func (s InstanceLastOperationPayload) Deal() error {
	glog.Infof("Deal with InstanceLastOperationPayload:%v.", s)
	controller := GetBaseController()
	status, err := controller.GetServiceInstanceStatus(s.InstanceID, s.Request.ServiceID, s.Request.PlanID, "")
	response := brokerapi.WorkerResponse{}
	if err != nil {
		response.Message = ""
		response.Err = err
	} else {
		state := ""
		if status == StateProvisionInstanceInProgress {
			state = brokerapi.StateInProgress
		} else if status == StateProvisionInstanceSucceeded {
			state = brokerapi.StateSucceeded
		} else if status == StateProvisionInstanceFailed {
			state = brokerapi.StateFailed
		}
		response.Message = state
		response.Err = nil
	}
	glog.Infof("InstanceLastOperationPayload with response:%v.", response)
	s.Result <- response
	return err
}
