package brokerapi

// ServiceBinding represents a binding to a service instance
type ServiceBinding struct {
	ID                string                 `json:"id"`
	ServiceID         string                 `json:"service_id"`
	AppID             string                 `json:"app_id"`
	ServicePlanID     string                 `json:"service_plan_id"`
	PrivateKey        string                 `json:"private_key"`
	ServiceInstanceID string                 `json:"service_instance_id"`
	BindResource      map[string]interface{} `json:"bind_resource,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
}

// BindingRequest represents a request to bind to a service instance
type BindingRequest struct {
	AppGUID      string                 `json:"app_guid,omitempty"`
	PlanID       string                 `json:"plan_id,omitempty"`
	ServiceID    string                 `json:"service_id,omitempty"`
	BindResource map[string]interface{} `json:"bind_resource,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

// CreateServiceBindingResponse represents a response to a service binding
// request
type CreateServiceBindingResponse struct {
	Credentials Credential `json:"credentials"`
	// Async indicates whether the broker is handling the provision request
	// asynchronously.
	//Async bool `json:"async"`
}

// Credential represents connection details, username, and password that are
// provisioned when a consumer binds to a service instance
type Credential map[string]interface{}

// BindingLastOperationRequest represents a request to a broker to give the
// state of the action on a binding it is completing asynchronously.
type BindingLastOperationRequest struct {
	// InstanceID is the instance of the service to query the last operation
	// for.
	InstanceID string `json:"instance_id"`
	// BindingID is the binding to query the last operation for.
	BindingID string `json:"binding_id"`
	// ServiceID is the ID of the service the instance is provisioned from.
	// Optional, but recommended.
	ServiceID string `json:"service_id,omitempty"`
	// PlanID is the ID of the plan the instance is provisioned from.
	// Optional, but recommended.
	PlanID string `json:"plan_id,omitempty"`
	// OperationKey is the operation key provided by the broker in the response
	// to the initial request. Optional, but must be sent if supplied in the
	// response to the original request.
	OperationKey string `json:"operation,omitempty"`
	// OriginatingIdentity requires a client API version >= 2.13.
	//
	// OriginatingIdentity is the identity on the platform of the user making
	// this request.
	//OriginatingIdentity *OriginatingIdentity `json:"originatingIdentity,omitempty"`
}

// GetBindingRequest represents a request to do a GET on a particular binding.
type GetBindingRequest struct {
	// InstanceID is the ID of the instance the binding is for.
	InstanceID string `json:"instance_id"`
	// BindingID is the ID of the binding to delete.
	BindingID string `json:"binding_id"`
}

// GetBindingResponse is sent as the response to doing a GET on a particular
// binding.
type GetBindingResponse struct {
	// Credentials is a free-form hash of credentials that can be used by
	// applications or users to access the service.
	Credentials Credential `json:"credentials,omitempty"`
	// SyslogDrainURl is a URL to which logs must be streamed. CF-specific. May
	// only be supplied by a service that declares a requirement for the
	// 'syslog_drain' permission.
	SyslogDrainURL *string `json:"syslog_drain_url,omitempty"`
	// RouteServiceURL is a URL to which the platform must proxy requests to the
	// application the binding is for. CF-specific. May only be supplied by a
	// service that declares a requirement for the 'route_service' permission.
	RouteServiceURL *string `json:"route_service_url,omitempty"`
	// VolumeMounts is an array of configuration string for mounting volumes.
	// CF-specific. May only be supplied by a service that declares a
	// requirement for the 'volume_mount' permission.
	VolumeMounts []interface{} `json:"volume_mounts,omitempty"`
	// Parameters is configuration parameters for the binding.
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}
