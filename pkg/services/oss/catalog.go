package oss

import (
	"fmt"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
)

type ServicePlanMetadata struct {
	Engine           string `json:"engine"`
	EngineVersion    string `json:"engine_version"`
	Class            string `json:"class"`
	CPU              string `json:"cpu"`
	Memory           string `json:"memory"`
	Storage          string `json:"storage"`
	HighAvailability string `json:"high_availability"`
	Type             string `json:"type"`
	VpcID            string `json:"vpcID"`
	VSwitchID        string `json:"vswitchID"`
}

const (
	SelfDefineServicePlan = "oss-edc2badc-d93b-4d9c-9d8e-da2f1c8c2222"
)

func getCatalog() *brokerapi.Catalog {

	instanceProperty := &brokerapi.ServiceInstanceSchema{
		Create: &brokerapi.InputParameters{
			Parameters: brokerapi.ParameterMapSchemas{
				Schema: "http://json-schema.org/draft-04/schema#",
				Type:   "object",
				Properties: map[string]brokerapi.ParameterProperty{
					"AccessKeyId": {
						Description: "AccessKeyId: ram accessKeyid used to create oss instance and temporary ram user to access it.",
						Type:        "string",
					},
					"AccessKeySecret": {
						Description: "AccessKeySecret: ram accessKeySecret used to create oss instance and temporary ram user to access it.",
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
					"AccessKeyId": {
						Description: "AccessKeyId: ram accessKeyid used to create oss instance and temporary ram user to access it.",
						Type:        "string",
					},
					"AccessKeySecret": {
						Description: "AccessKeySecret: ram accessKeyid used to create oss instance and temporary ram user to access it.",
						Type:        "string",
					},
					"Region": {
						Description: "Region: optional, the region of the oss instance, default is cluster's region",
						Type:        "string",
					},
					"Acl": {
						Description: "Acl: optional, ACL type: Private, Public Read, Public Read/Write, default is Private",
						Type:        "string",
					},
					"StorageClass": {
						Description: "StorageClass: optional, StorageClass type: Standard Storage, Infrequent Access Storage, Archive Storage default is  Standard Storage",
						Type:        "string",
					},
				},
			},
		},
	}

	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				ID:             "oss-997b8372-8dac-40ac-ae65-758b4a502222",
				Name:           "alibaba-cloud-oss",
				Description:    "Alibaba Cloud OSS",
				Bindable:       true,
				PlanUpdateable: true,
				Tags:           []string{"AlibabaCloud", "OSS"},
				Plans: []brokerapi.ServicePlan{
					{
						ID:          "oss-edc2badc-d93b-4d9c-9d8e-da2f1c8c2221",
						Name:        "default",
						Description: "ACLPrivate,StorageStandard",
						Free:        false,
						Schemas: &brokerapi.Schemas{
							ServiceInstances: instanceProperty,
						},
					},

					{
						ID:          SelfDefineServicePlan,
						Name:        "oss-self-define",
						Description: "oss-self-define",
						Free:        false,
						Schemas: &brokerapi.Schemas{
							ServiceInstances: optionalInstanceProperty,
						},
					},
				},
			},
		},
	}
}

func dealSelfDefinedServicePlanMetadata(parameter map[string]interface{}) (*ServicePlanMetadata, error) {
	return nil, nil
}

func findServicePlanMetadata(serviceID, planID string, parameter map[string]interface{}) (*ServicePlanMetadata, error) {

	catalog := getCatalog()

	for _, service := range catalog.Services {
		if service.ID == serviceID {
			for _, plan := range service.Plans {
				if plan.ID == planID {
					if plan.ID == SelfDefineServicePlan {
						return dealSelfDefinedServicePlanMetadata(parameter)
					}
					return nil, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("Unsupported Service Plan for service id: %s, plan id: %s", serviceID, planID)
}

func GetAccessInfoFromInstanceParameters(parameterIn map[string]interface{}) (
	accessKeyId, accessKeySecret, region, acl, storageClass string, err error) {

	vAccessKeyId, ok := parameterIn["AccessKeyId"]
	if ok {
		if value, ok := vAccessKeyId.(string); ok {
			accessKeyId = value
		} else {
			err = fmt.Errorf("Service instance AccessKeyId  type is wrong suppose to be string.")
			return
		}
	}

	vAccessKeySecret, ok := parameterIn["AccessKeySecret"]
	if ok {
		if value, ok := vAccessKeySecret.(string); ok {
			accessKeySecret = value
		} else {
			err = fmt.Errorf("Service instance vAccessKeySecret  type is wrong suppose to be string.")
			return
		}
	}

	vRegion, ok := parameterIn["Region"]
	if ok {
		if value, ok := vRegion.(string); ok {
			region = value
		} else {
			err = fmt.Errorf("Service instance region  type is wrong suppose to be string.")
			return
		}
	}

	vACL, ok := parameterIn["ACL"]
	if ok {
		if value, ok := vACL.(string); ok {
			acl = value
		} else {
			err = fmt.Errorf("Service instance acl type is wrong suppose to be string.")
			return
		}
	}

	vStorageClass, ok := parameterIn["StorageClass"]
	if ok {
		if value, ok := vStorageClass.(string); ok {
			storageClass = value
		} else {
			err = fmt.Errorf("Service instance storageClass type is wrong suppose to be string.")
			return
		}
	}
	return
}
