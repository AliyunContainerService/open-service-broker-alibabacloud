package controller

import (
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"testing"
)

type fakeBroker struct {
	info string
}

const (
	TestInstanceID = "997b8372-8dac-40ac-ae65-758b4a000000"
	TestServiceID  = "997b8372-8dac-40ac-ae65-758b4a111111"
	TestPlanID     = "997b8372-8dac-40ac-ae65-758b4a222222"
	TestBindingId  = "997b8372-8dac-40ac-ae65-758b4a333333"
)

func (b *fakeBroker) Catalog() (*brokerapi.Catalog, error) {
	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				ID:             TestServiceID,
				Name:           "fakeBroker",
				Description:    "a fake broker for test",
				Bindable:       true,
				PlanUpdateable: true,
				Plans: []brokerapi.ServicePlan{
					{
						ID:          TestPlanID,
						Name:        "fake-broker-plan",
						Description: "a fake broker service plan",
						Free:        false,
					},
				},
			},
		},
	}, nil
}

func (b *fakeBroker) GetInstanceStatus(instanceID, serviceID, planID string, parameterIn map[string]interface{}) (bool, error) {
	return true, nil
}

func (b *fakeBroker) Provision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (b *fakeBroker) Deprovision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) error {
	return nil
}

func (b *fakeBroker) GetBindingStatus(instanceID, serviceID, planID, bindingID string, parameterInstance, parameterIn map[string]interface{}) (bool, error) {
	return true, nil
}

func (b *fakeBroker) Bind(instanceID, serviceID, planID, bindingID string, parameterInstance,
	parameterIn map[string]interface{}) (map[string]interface{}, brokerapi.Credential, error) {
	return nil, brokerapi.Credential{}, nil
}

func (b *fakeBroker) UnBind(instanceID, serviceID, planID, bindingID string, parameterInstance, parameterIn map[string]interface{}) error {
	return nil
}

func TestBaseController(t *testing.T) {
	baseController := NewBaseController()
	if baseController == nil {
		t.Fatalf("NewBaseController failed.")
	}

	RegisterBroker("fakeBroker", &fakeBroker{info: "fakeBroker"})

	baseController.asyncEngine.storageProvider = nil

	fakeBroker := baseController.getBrokerByName("fakeBroker")
	if fakeBroker == nil {
		t.Fatalf("getBrokerByName failed.")
	}

	_, err := baseController.GetCatalog()
	if err != nil {
		t.Fatalf("GetCatalog failed.")
	}

	gotBroker := baseController.getBrokerByName("fakeBroker")
	if gotBroker == nil {
		t.Fatalf("getBrokerByName failed.")
	}

	err = baseController.CreateServiceInstance(TestInstanceID, TestServiceID, TestPlanID, nil)
	if err != nil {
		t.Fatalf("CreatebaseControllerInstance failed.")
	}

	_, err = baseController.GetServiceInstanceStatus(TestInstanceID, TestServiceID, TestPlanID, "")
	if err != nil {
		t.Fatalf("GetServiceInstanceStatus failed.")
	}

	_, err = baseController.CreateServiceBinding(TestInstanceID,
		TestServiceID, TestPlanID, TestBindingId, nil)
	if err != nil {
		t.Fatalf("CreateServiceBinding failed.")
	}

	_, err = baseController.GetServiceBindingStatus(TestInstanceID,
		TestServiceID, TestPlanID, TestBindingId, "")
	if err != nil {
		t.Fatalf("GetServiceBindingStatus failed.")
	}

	err = baseController.DeleteServiceBinding(TestInstanceID,
		TestServiceID, TestPlanID, TestBindingId, nil)
	if err != nil {
		t.Fatalf("DeleteServiceBinding failed.")
	}

	err = baseController.DeleteServiceInstance(TestInstanceID, TestServiceID, TestPlanID, nil)
	if err != nil {
		t.Fatalf("DeleteServiceInstance failed.")
	}
}
