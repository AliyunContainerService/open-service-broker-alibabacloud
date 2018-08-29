package rds

import (
	"fmt"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"strconv"
)

type ServiceInstanceMetadata struct {
	Engine        string `json:"engine"`
	EngineVersion string `json:"engine_version"`
	Class         string `json:"class"`
	Storage       string `json:"storage"`
	ZoneId        string `json:"zoneId"`
	NetworkType   string `json:"networktype"`
	VpcID         string `json:"vpcID"`
	VSwitchID     string `json:"vswitchID"`
}

type ServiceInstanceBasicMetadata struct {
	Engine        string `json:"engine"`
	EngineVersion string `json:"engine_version"`
	Class         string `json:"class"`
	Storage       string `json:"storage"`
}

const (
	SelfDefineServicePlan = "rds-edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1e"
)

func getCatalog() *brokerapi.Catalog {

	bindingProperty := &brokerapi.ServiceBindingSchema{
		Create: &brokerapi.InputParameters{
			Parameters: brokerapi.ParameterMapSchemas{
				Schema: "http://json-schema.org/draft-04/schema#",
				Type:   "object",
				Properties: map[string]brokerapi.ParameterProperty{
					"Username": {
						Description: "Username: the username of rds, 2 to 16 characters including lower-case letters, digits, or underscores. It must begin with a letter and end with a letter or a digit.",
						Type:        "string",
					},
					"Password": {
						Description: "Password: the password of rds, 8 to 32 characters including at least three of the following:Capital letters, Lower-case letters,Digits,Special characters ( !@#$%^&*()_-+=)",
						Type:        "string",
					},
				},
			},
		},
	}

	instanceProperty := &brokerapi.ServiceInstanceSchema{
		Create: &brokerapi.InputParameters{
			Parameters: brokerapi.ParameterMapSchemas{
				Schema: "http://json-schema.org/draft-04/schema#",
				Type:   "object",
				Properties: map[string]brokerapi.ParameterProperty{
					"Engine": {
						Description: "Engine:MySQL",
						Type:        "string",
					},
					"EngineVersion": {
						Description: "EngineVersion:5.5,5.6,5.7",
						Type:        "string",
					},
					"Class": {
						Description: "Class: mysql-n2-medium-1, and so on",
						Type:        "string",
					},
					"Storage": {
						Description: "Storage:20,25,30",
						Type:        "string",
					},
					"NetworkType": {
						Description: "NetworkType: optional, vpc, classic, default is classic",
						Type:        "string",
					},
					"VSwitchID": {
						Description: "VSwitchID: optional, the vSwitchID of rds, default is cluster's vSwitchID",
						Type:        "string",
					},
					"VpcID": {
						Description: "VpcID: optional, the VpcID of rds, default is cluster's VpcID",
						Type:        "string",
					},
					"ZoneId": {
						Description: "ZoneId: optional, the ZoneId of rds, default is cluster's ZoneId",
						Type:        "string",
					},
				},
			},
		},
	}

	optionalInstanceProperty := &brokerapi.ServiceInstanceSchema{
		Create: &brokerapi.InputParameters{
			Parameters: brokerapi.ParameterMapSchemas{
				Schema: "http://json-schema.org/draft-04/schema#",
				Type:   "object",
				Properties: map[string]brokerapi.ParameterProperty{
					"NetworkType": {
						Description: "NetworkType: optional, vpc, classic, default is classic",
						Type:        "string",
					},
					"VSwitchID": {
						Description: "VSwitchID: optional, the vSwitchID of rds, default is cluster's vSwitchID",
						Type:        "string",
					},
					"VpcID": {
						Description: "VpcID: optional, the VpcID of rds, default is cluster's VpcID",
						Type:        "string",
					},
					"ZoneId": {
						Description: "ZoneId: optional, the ZoneId of rds, default is cluster's ZoneId",
						Type:        "string",
					},
				},
			},
		},
	}

	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				ID:             "rds-997b8372-8dac-40ac-ae65-758b4a5075a5",
				Name:           "alibaba-cloud-rds-mysqldb",
				Description:    "Alibaba Cloud RDS Database for MySQL (Experimental)",
				Bindable:       true,
				PlanUpdateable: true,
				Tags:           []string{"AlibabaCloud", "MySQL", "Database", "RDS"},
				Plans: []brokerapi.ServicePlan{
					{
						ID:          "rds-427559f1-bf2a-45d3-8844-32374a3e58aa",
						Name:        "rds-mysql-t1-small",
						Description: "Basic Tier, 1 Core 1G 20G",
						Free:        false,
						Metadata: ServiceInstanceBasicMetadata{
							Engine:        "MySQL",
							EngineVersion: "5.6",
							Class:         "rds.mysql.t1.small",
							Storage:       "20",
						},
						Schemas: &brokerapi.Schemas{
							ServiceInstances: optionalInstanceProperty,
							ServiceBindings:  bindingProperty,
						},
					},
					{
						ID:          "rds-edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1c",
						Name:        "rds-mysql-s1-small",
						Description: "Standard Tier, 1 Core 2G 20G",
						Free:        false,
						Metadata: ServiceInstanceBasicMetadata{
							Engine:        "MySQL",
							EngineVersion: "5.6",
							Class:         "rds.mysql.s1.small",
							Storage:       "20",
						},
						Schemas: &brokerapi.Schemas{
							ServiceInstances: optionalInstanceProperty,
							ServiceBindings:  bindingProperty,
						},
					},
					{
						ID:          "rds-edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1d",
						Name:        "rds-mysql-s2-large",
						Description: "Advanced Tier, 2 Core 4G 20G",
						Free:        false,
						Metadata: ServiceInstanceBasicMetadata{
							Engine:        "MySQL",
							EngineVersion: "5.6",
							Class:         "rds.mysql.s2.large",
							Storage:       "20",
						},
						Schemas: &brokerapi.Schemas{
							ServiceInstances: optionalInstanceProperty,
							ServiceBindings:  bindingProperty,
						},
					},
					{
						ID:          SelfDefineServicePlan,
						Name:        "rds-self-define",
						Description: "rds-self-define service plan is a generic way to create RDS instance with user provided parameters",
						Free:        false,
						Schemas: &brokerapi.Schemas{
							ServiceInstances: instanceProperty,
							ServiceBindings:  bindingProperty,
						},
					},
				},
			},
		},
	}
}

func dealServiceInstanceMetadata(parameter map[string]interface{}) (*ServiceInstanceMetadata, error) {

	metadata := &ServiceInstanceMetadata{}
	engine, ok := parameter["Engine"]
	if !ok {
		metadata.Engine = "MySQL"
	} else {
		if value, ok := engine.(string); ok {
			metadata.Engine = value
		} else {
			return nil, fmt.Errorf("Service plan's Engine type is wrong suppose to be string.")
		}
	}

	engineVersion, ok := parameter["EngineVersion"]
	if !ok {
		metadata.EngineVersion = "5.6"
	} else {
		if value, ok := engineVersion.(float64); ok {
			metadata.EngineVersion = strconv.FormatFloat(value, 'f', -1, 64)
		} else {
			return nil, fmt.Errorf("Service plan's EngineVersion type is wrong suppose to be float64.")
		}
	}

	class, ok := parameter["Class"]
	if !ok {
		metadata.Class = "rds.mysql.s2.large"
	} else {
		if value, ok := class.(string); ok {
			metadata.Class = value
		} else {
			return nil, fmt.Errorf("Service plan's Class type is wrong suppose to be string.")
		}
	}

	storage, ok := parameter["Storage"]
	if !ok {
		metadata.Storage = "20"
	} else {
		if value, ok := storage.(float64); ok {
			metadata.Storage = strconv.FormatFloat(value, 'f', -1, 64)
		} else {
			return nil, fmt.Errorf("Service plan's Storage type is wrong suppose to be float64.")
		}
	}

	err := dealOptionalServicePlanMetadata(metadata, parameter)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func dealOptionalServicePlanMetadata(metadata *ServiceInstanceMetadata, parameter map[string]interface{}) error {

	zoneId, ok := parameter["ZoneId"]
	if !ok {
		metadata.ZoneId = ""
	} else {
		if value, ok := zoneId.(string); ok {
			metadata.ZoneId = value
		} else {
			return fmt.Errorf("Service plan's HighAvailability type is wrong suppose to be string.")
		}
	}

	networkType, ok := parameter["NetworkType"]
	if !ok {
		metadata.NetworkType = ""
	} else {
		if value, ok := networkType.(string); ok {
			metadata.NetworkType = value
		} else {
			return fmt.Errorf("Service plan's Type type is wrong suppose to be string.")
		}
	}

	vpcID, ok := parameter["VpcID"]
	if !ok {
		metadata.VpcID = ""
	} else {
		if value, ok := vpcID.(string); ok {
			metadata.VpcID = value
		} else {
			return fmt.Errorf("Service plan's VpcID type is wrong suppose to be string.")
		}
	}

	vSwitchID, ok := parameter["VSwitchID"]
	if !ok {
		metadata.VSwitchID = ""
	} else {
		if value, ok := vSwitchID.(string); ok {
			metadata.VSwitchID = value
		} else {
			return fmt.Errorf("Service plan's VSwitchID type is wrong suppose to be string.")
		}
	}

	return nil
}

func findServicePlanMetadata(serviceID, planID string, parameter map[string]interface{}) (*ServiceInstanceMetadata, error) {

	catalog := getCatalog()

	for _, service := range catalog.Services {
		if service.ID == serviceID {
			for _, plan := range service.Plans {
				if plan.ID == planID {
					if plan.ID == SelfDefineServicePlan {
						metadata, err := dealServiceInstanceMetadata(parameter)
						if err != nil {
							return nil, err
						}
						return metadata, nil
					} else {
						metadata := &ServiceInstanceMetadata{}
						basicMetadata, _ := plan.Metadata.(ServiceInstanceBasicMetadata)
						err := dealOptionalServicePlanMetadata(metadata, parameter)
						if err != nil {
							return nil, err
						}
						metadata.Engine = basicMetadata.Engine
						metadata.EngineVersion = basicMetadata.EngineVersion
						metadata.Class = basicMetadata.Class
						metadata.Storage = basicMetadata.Storage
						return metadata, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("Unsupported Service Plan for service id: %s, plan id: %s", serviceID, planID)
}

func GetAccountInfoFromBindingParameters(parameter map[string]interface{}) (string, string) {

	username := defaultUsername
	vUsername, ok := parameter["Username"]
	if ok {
		if value, ok := vUsername.(string); ok {
			username = value
		}
	}

	password := defaultPassword
	vPassword, ok := parameter["Password"]
	if ok {
		if value, ok := vPassword.(string); ok {
			password = value
		}
	}
	return username, password
}
