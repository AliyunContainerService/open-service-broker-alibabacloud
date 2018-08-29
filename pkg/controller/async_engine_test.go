package controller

import (
	"testing"
)

func TestNewAsyncEngine(t *testing.T) {
	asyncEngine := NewAsyncEngine()
	if asyncEngine == nil || asyncEngine.instancesRunMapMutex == nil || asyncEngine.instancesRunMap == nil {
		t.Fatalf("NewAsyncEngine create async engine failed.")
	}

	asyncEngine.storageProvider = nil
	instanceRunInfo := &InstanceRunInfo{InstanceId: "0123456789"}
	err := asyncEngine.AddAsyncInstance(instanceRunInfo, true)
	if err != nil || len(asyncEngine.instancesRunMap) != 1 {
		t.Fatalf("AddAsyncInstance add instance information failed.")
	}

	err = asyncEngine.AddAsyncInstance(instanceRunInfo, true)
	if err == nil {
		t.Fatalf("AddAsyncInstance repeat add instance information failed.")
	}

	instanceRunInfoNew := &InstanceRunInfo{InstanceId: "0123456789", ServiceID: "0123456789"}
	err = asyncEngine.UpdateAsyncInstance(instanceRunInfoNew)
	if err != nil || len(asyncEngine.instancesRunMap) != 1 {
		t.Fatalf("UpdateAsyncInstance update instance information failed.")
	}

	instanceRunInfoNew.ServiceID = "9876543210"
	err = asyncEngine.UpdateAsyncInstance(instanceRunInfoNew)
	if err != nil {
		t.Fatalf("UpdateAsyncInstance update non instance information failed.")
	}

	instanceInfo := asyncEngine.GetAsyncInstance("0123456789")
	if instanceInfo == nil || instanceInfo.InstanceId != "0123456789" || instanceInfo.ServiceID != "9876543210" {
		t.Fatalf("GetAsyncInstance get instance information failed.")
	}

	instanceInfo = asyncEngine.GetAsyncInstance("9876543210")
	if instanceInfo != nil {
		t.Fatalf("GetAsyncInstance get non exist instance information failed.")
	}

	err = asyncEngine.DeleteAsyncInstance("0123456789")
	if err != nil || len(asyncEngine.instancesRunMap) != 0 {
		t.Fatalf("DeleteAsyncInstance add instance information failed.")
	}

	err = asyncEngine.DeleteAsyncInstance("0123456789")
	if err == nil {
		t.Fatalf("DeleteAsyncInstance repeat add instance information failed.")
	}
}
