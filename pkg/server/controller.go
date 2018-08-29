package server

import (
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
)

// Controller defines the APIs that the base function of Alibaba Cloud service broker
// is expected to support. Implementation should be concurrency-safe
type Controller interface {
	GetCatalog() (*brokerapi.Catalog, error)

	GetServiceInstanceStatus(instanceId, serviceID, PlanID string, operation string) (string, error)
	CreateServiceInstance(instanceId, serviceID, PlanID string, Parameter map[string]interface{}) error
	DeleteServiceInstance(instanceId, serviceID, PlanID string, Parameter map[string]string) error

	GetServiceBindingStatus(instanceId, serviceID, planID, bindingId string, operation string) (string, error)
	CreateServiceBinding(instanceId, serviceID, planID, bindingId string, Parameter map[string]interface{}) (brokerapi.Credential, error)
	DeleteServiceBinding(instanceId, serviceID, planID, bindingId string, Parameter map[string]string) error
}
