package controller

import (
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"github.com/golang/glog"
)

type BindingPayload struct {
	InstanceID string
	BindingID  string
	Request    brokerapi.BindingRequest
	Result     chan brokerapi.WorkerResponse
}

func (s BindingPayload) Deal() error {
	glog.Infof("BindingPayload in.%s\n", s)
	response := brokerapi.WorkerResponse{}
	controller := GetBaseController()
	credential, err := controller.CreateServiceBinding(s.InstanceID,
		s.Request.ServiceID, s.Request.PlanID, s.BindingID, s.Request.Parameters)
	if err != nil {
		response.Message = nil
		response.Err = err
	} else {
		response.Message = credential
		response.Err = nil
	}
	glog.Infof("BindingPayload get response.%s\n", s)
	s.Result <- response
	return err
}

type UnBindingPayload struct {
	InstanceID string
	BindingID  string
	PlanID     string
	ServiceID  string
	Result     chan brokerapi.WorkerResponse
}

func (s UnBindingPayload) Deal() error {
	response := brokerapi.WorkerResponse{}

	controller := GetBaseController()
	err := controller.DeleteServiceBinding(s.InstanceID, s.ServiceID, s.PlanID, s.BindingID, nil)
	if err != nil {
		response.Message = nil
		response.Err = err
	} else {
		response.Message = "Unbind success."
		response.Err = nil
	}
	s.Result <- response
	return err
}

type BindingLastOperationPayload struct {
	InstanceID string
	BindingID  string
	Request    brokerapi.BindingLastOperationRequest
	Result     chan brokerapi.WorkerResponse
}

func (s BindingLastOperationPayload) Deal() error {
	controller := GetBaseController()
	status, err := controller.GetServiceBindingStatus(s.InstanceID, s.Request.ServiceID,
		s.Request.PlanID, s.BindingID, s.Request.OperationKey)
	response := brokerapi.WorkerResponse{Message: status, Err: err}
	s.Result <- response
	return err
}
