package controller

import (
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"testing"
	"time"
)

func TestBindingCreateDeal(t *testing.T) {
	aliBroker := NewBaseController()
	if aliBroker == nil {
		t.Fatalf("BaseControllerInit failed.")
	}

	RegisterBroker("fakeBroker", &fakeBroker{info: "fakeBroker"})

	aliBroker.asyncEngine.storageProvider = nil

	t.Logf("aliBroker servicemap:%v.", aliBroker.serviceMap)

	createServiceInstancePayload := CreateServiceInstancePayload{
		InstanceID: TestInstanceID,
		Request:    brokerapi.CreateServiceInstanceRequest{PlanID: TestPlanID, ServiceID: TestServiceID},
		Result:     make(chan brokerapi.WorkerResponse, 2),
	}

	t.Logf("createServiceInstancePayload begin.")

	go func(Result chan brokerapi.WorkerResponse) {
		<-Result
		return
	}(createServiceInstancePayload.Result)

	err := createServiceInstancePayload.Deal()
	if err != nil {
		t.Fatalf("createServiceInstancePayload failed.")
	}

	t.Logf("createServiceInstancePayload success.sleep 10s")

	time.Sleep((AsyncLoopTime + 1) * time.Second)

	bindingPayload := BindingPayload{
		InstanceID: TestInstanceID,
		BindingID:  TestBindingId,
		Request: brokerapi.BindingRequest{
			PlanID: TestPlanID, ServiceID: TestServiceID,
		},
		Result: make(chan brokerapi.WorkerResponse),
	}

	go func(Result chan brokerapi.WorkerResponse) {
		<-Result
		return
	}(bindingPayload.Result)

	err = bindingPayload.Deal()
	if err != nil {
		t.Fatalf("TestBindingCreateDeal failed.err:%v", err)
	}

	bindingLastOperationPayload := BindingLastOperationPayload{
		InstanceID: TestInstanceID,
		BindingID:  TestBindingId,
		Request: brokerapi.BindingLastOperationRequest{
			PlanID: TestPlanID, ServiceID: TestServiceID,
		},
		Result: make(chan brokerapi.WorkerResponse),
	}

	go func(Result chan brokerapi.WorkerResponse) {
		<-Result
		return
	}(bindingLastOperationPayload.Result)

	err = bindingLastOperationPayload.Deal()
	if err != nil {
		t.Fatalf("TestBindingLastOperationDeal failed.")
	}

	UnBindingPayload := UnBindingPayload{
		InstanceID: TestInstanceID,
		BindingID:  TestBindingId,
		PlanID:     TestPlanID,
		ServiceID:  TestServiceID,
		Result:     make(chan brokerapi.WorkerResponse),
	}

	go func(Result chan brokerapi.WorkerResponse) {
		<-Result
		return
	}(UnBindingPayload.Result)

	err = UnBindingPayload.Deal()
	if err != nil {
		t.Fatalf("TestUnBindingDeal failed.")
	}
}
