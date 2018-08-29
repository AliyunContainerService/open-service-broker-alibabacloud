package ots

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
	SelfDefineServicePlan = "edc2badc-d93b-4d9c-9d8e-da2f1c8c3333"
)

func getCatalog() *brokerapi.Catalog {

	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				ID:             "997b8372-8dac-40ac-ae65-758b4a503330",
				Name:           "alibaba-cloud-ots",
				Description:    "Alibaba Cloud OTS",
				Bindable:       true,
				PlanUpdateable: true,
				Tags:           []string{"AlibabaCloud", "OTS"},
				Plans: []brokerapi.ServicePlan{
					{
						ID:          "edc2badc-d93b-4d9c-9d8e-da2f1c8c3331",
						Name:        "ots-default",
						Description: "ACLPrivate,StorageStandard",
						Free:        false,
					},
					{
						ID:          SelfDefineServicePlan,
						Name:        "ots-self-define",
						Description: "ots-self-define",
						Free:        false,
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
