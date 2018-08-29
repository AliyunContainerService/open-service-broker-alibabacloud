package controller

import (
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"testing"
	"time"
)

func TestInstanceDeal(t *testing.T) {
	baseController := NewBaseController()
	if baseController == nil {
		t.Fatalf("AliBrokerControllerInit failed.")
	}

	RegisterBroker("fakeBroker", &fakeBroker{info: "fakeBroker"})

	baseController.asyncEngine.storageProvider = nil

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

	instanceLastOperationPayload := InstanceLastOperationPayload{
		InstanceID: TestInstanceID,
		Request:    brokerapi.LastOperationRequest{PlanID: TestPlanID, ServiceID: TestServiceID},
		Result:     make(chan brokerapi.WorkerResponse),
	}

	go func(Result chan brokerapi.WorkerResponse) {
		<-Result
		return
	}(instanceLastOperationPayload.Result)

	err = instanceLastOperationPayload.Deal()
	if err != nil {
		t.Fatalf("TestInstanceLastOperationDeal failed.")
	}

	deleteServiceInstancePayload := DeleteServiceInstancePayload{
		InstanceID: TestInstanceID,
		Request:    brokerapi.DeleteServiceInstanceRequest{PlanID: TestPlanID, ServiceID: TestServiceID},
		Result:     make(chan brokerapi.WorkerResponse),
	}

	go func(Result chan brokerapi.WorkerResponse) {
		<-Result
		return
	}(deleteServiceInstancePayload.Result)

	err = deleteServiceInstancePayload.Deal()
	if err != nil {
		t.Fatalf("TestInstanceDeleteDeal failed.")
	}
}
