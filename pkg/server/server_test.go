package server

import (
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
)

//
// Test of server /v2/catalog endpoint.
//

// Make sure that Controller stub implements controller.Controller interface

type fakeController struct {
	name string
}

func (a *fakeController) GetCatalog() (*brokerapi.Catalog, error) {
	return nil, nil
}

func (a *fakeController) GetAliBrokerInstanceStatus(instanceId, serviceID, PlanID string, operation string) (string, error) {
	return "", nil
}
func (a *fakeController) CreateAliBrokerInstance(instanceId, serviceID, PlanID string, Parameter map[string]interface{}) error {
	return nil
}
func (a *fakeController) DeleteAliBrokerInstance(instanceId, serviceID, PlanID string, Parameter map[string]string) error {
	return nil
}

func (a *fakeController) GetAliBrokerBindingStatus(instanceId, serviceID, planID, bindingId string, operation string) (string, error) {
	return "", nil
}

func (a *fakeController) CreateAliBrokerBinding(instanceId, serviceID, planID, bindingId string,
	Parameter map[string]interface{}) (brokerapi.Credential, error) {
	return brokerapi.Credential{}, nil
}

func (a *fakeController) DeleteAliBrokerBinding(instanceId, serviceID, planID, bindingId string, Parameter map[string]string) error {
	return nil
}
