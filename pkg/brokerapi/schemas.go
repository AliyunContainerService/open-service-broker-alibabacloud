package brokerapi

// Schemas represents a plan's schemas for service instance and binding create
// and update.
type Schemas struct {
	ServiceInstances *ServiceInstanceSchema `json:"service_instance,omitempty"`
	ServiceBindings  *ServiceBindingSchema  `json:"service_binding,omitempty"`
}

// ServiceInstanceSchema represents a plan's schemas for a create and update
// of a service instance.
type ServiceInstanceSchema struct {
	Create *InputParameters `json:"create,omitempty"`
	Update *InputParameters `json:"update,omitempty"`
}

// ServiceBindingSchema represents a plan's schemas for the parameters
// accepted for binding creation.
type ServiceBindingSchema struct {
	Create *InputParameters `json:"create,omitempty"`
}

// InputParameters represents a schema for input parameters for creation or
// update of an API resource.
type InputParameters struct {
	Parameters interface{} `json:"parameters,omitempty"`
}

// ParameterProperty describe a parameter's information and type
type ParameterProperty struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}

// ParameterMapSchemas describe the parameters information
type ParameterMapSchemas struct {
	Schema     string                       `json:"$schema"`
	Type       string                       `json:"type"`
	Properties map[string]ParameterProperty `json:"properties"`
}
