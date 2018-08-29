package brokerapi

// ServiceBroker the standard interface of open service broker
// which each specific service provider must implement.
type ServiceBroker interface {
	Catalog() (*Catalog, error)

	GetInstanceStatus(instanceID, serviceID, planID string, parameterIn map[string]interface{}) (bool, error)
	Provision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) (map[string]interface{}, error)
	Deprovision(instanceID, serviceID, planID string, parameterIn map[string]interface{}) error

	GetBindingStatus(instanceID, serviceID, planID, bindingID string, parameterInstance, parameterIn map[string]interface{}) (bool, error)
	Bind(instanceID, serviceID, planID, bindingID string, parameterInstance, parameterIn map[string]interface{}) (map[string]interface{}, Credential, error)
	UnBind(instanceID, serviceID, planID, bindingID string, parameterInstance, parameterIn map[string]interface{}) error
}
