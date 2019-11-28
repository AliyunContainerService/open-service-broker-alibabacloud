package polardb

import (
	"errors"
	"fmt"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"reflect"
	"strconv"
)

type ServiceInstanceMetadata struct {
	Engine        string `json:"engine"`
	EngineVersion string `json:"engine_version"`
	Class         string `json:"class"`
	ZoneId        string `json:"zoneId"`
	VpcID         string `json:"vpcID"`
	VSwitchID     string `json:"vswitchID"`

	AccessKeyId     string `json:"accessId"`
	AccessKeySecret string `json:"accessSecret"`

	RegionId        string `json:"regionId"`
	SecurityIps     string `json:"securityIps"`
}

type ServiceInstanceBasicMetadata struct {
	Engine        string `json:"engine"`
	EngineVersion string `json:"engine_version"`
	Class         string `json:"class"`
	Storage       string `json:"storage"`
}

const (
	SelfDefineServicePlan = "polardb-edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1e"
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
						Description: "Engine:MySQL|PostgreSQL|Oracle",
						Type:        "string",
					},
					"EngineVersion": {
						Description: "EngineVersion:MySQL：5.6/8.0, PostgreSQL：11, Oracle：11",
						Type:        "string",
					},
					"Class": {
						Description: "Class: polar.mysql.x2.medium, and so on, see https://help.aliyun.com/document_detail/68498.html",
						Type:        "string",
					},
					"VpcID": {
						Description: "VpcID, the VpcID of rds",
						Type:        "string",
					},
					"VSwitchID": {
						Description: "VSwitchID, the vSwitchID of rds",
						Type:        "string",
					},
					"RegionId": {
						Description: "RegionId, the RegionId of rds",
						Type:        "string",
					},
					"ZoneId": {
						Description: "ZoneId, the ZoneId of rds",
						Type:        "string",
					},
					"AccessKeyId": {
						Description: "AccessKeyId: ram accessKeyid used to create oss instance and temporary ram user to access it.",
						Type:        "string",
					},
					"AccessKeySecret": {
						Description: "AccessKeySecret: ram accessKeyid used to create oss instance and temporary ram user to access it.",
						Type:        "string",
					},
					"SecurityIps": {
						Description: "SecurityIps: The IP addresses to be added to the IP address whitelist group, example: " +
							"https://github.com/AlibabaCloudDocs/polardb/blob/master/intl.en-US/API%20Reference/Whitelist%20management/ModifyDBClusterAccessWhiteList.md",
						Type:        "string",
					},
				},
			},
		},
	}

	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				ID:             "polardb-997b8372-8dac-40ac-ae65-758b4a5075a5",
				Name:           "alibaba-cloud-polardb",
				Description:    "Alibaba Cloud PolarDB Database (Experimental)",
				Bindable:       true,
				PlanUpdateable: true,
				Tags:           []string{"AlibabaCloud", "PolarDB", "Database"},
				Plans: []brokerapi.ServicePlan{
					{
						ID:          SelfDefineServicePlan,
						Name:        "polardb-self-define",
						Description: "polardb-self-define service plan is a generic way to create PolarDB instance with user provided parameters",
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

func dealServiceInstanceMetadataValueString(parameter map[string]interface{}, paramKey string, metadata *ServiceInstanceMetadata) error {
	paramValue, ok := parameter[paramKey]
	if !ok {
		return fmt.Errorf("please specify %v", paramKey)
	} else {
		if value, ok := paramValue.(string); ok {
			reflect.ValueOf(metadata).Elem().FieldByName(paramKey).SetString(value)
		} else {
			return fmt.Errorf("Service plan's %v type is wrong suppose to be string.", paramKey)
		}
	}
	return nil
}

func dealServiceInstanceMetadata(parameter map[string]interface{}) (*ServiceInstanceMetadata, error) {

	metadata := &ServiceInstanceMetadata{}
	if err := dealServiceInstanceMetadataValueString(parameter, "Engine", metadata); err != nil {
		return nil, err
	}

	engineVersion, ok := parameter["EngineVersion"]
	if !ok {
		return nil, fmt.Errorf("please specify EngineVersion")
	} else {
		if value, ok := engineVersion.(float64); ok {
			metadata.EngineVersion = strconv.FormatFloat(value, 'f', -1, 64)
		} else {
			return nil, fmt.Errorf("Service plan's EngineVersion type is wrong suppose to be float64.")
		}
	}

	if err := dealServiceInstanceMetadataValueString(parameter, "Class", metadata); err != nil {
		return nil, err
	}
	if err := dealServiceInstanceMetadataValueString(parameter, "ZoneId", metadata); err != nil {
		return nil, err
	}
	if err := dealServiceInstanceMetadataValueString(parameter, "VpcID", metadata); err != nil {
		return nil, err
	}
	if err := dealServiceInstanceMetadataValueString(parameter, "VSwitchID", metadata); err != nil {
		return nil, err
	}
	if err := dealServiceInstanceMetadataValueString(parameter, "ZoneId", metadata); err != nil {
		return nil, err
	}
	if err := dealServiceInstanceMetadataValueString(parameter, "RegionId", metadata); err != nil {
		return nil, err
	}
	if err := dealServiceInstanceMetadataValueString(parameter, "AccessKeyId", metadata); err != nil {
		return nil, err
	}
	if err := dealServiceInstanceMetadataValueString(parameter, "AccessKeySecret", metadata); err != nil {
		return nil, err
	}
	if err := dealServiceInstanceMetadataValueString(parameter, "SecurityIps", metadata); err != nil {
		return nil, err
	}

	return metadata, nil
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
						return nil, errors.New("can not be here!!")
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
