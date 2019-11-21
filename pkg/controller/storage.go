package controller

import (
	"context"
	"time"

	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"sort"
)

type Storage interface {
	Init()
	WriteObject(instanceInfo InstanceRunInfo)
	GetObject(instanceID string) (InstanceRunInfo, error)
	DeleteObject(instanceID string)

	GetAllObject() ([]*InstanceRunInfo, error)
}

type StorageProvider struct {
	client *client.Client
}

func getEtcdClient() (*client.Client, error) {
	cfg := client.Config{
		Endpoints:               []string{"http://localhost:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	client, err := client.New(cfg)
	return &client, err
}

func NewStorageProvider() *StorageProvider {
	sp := StorageProvider{}
	c, err := getEtcdClient()
	if err != nil {
		glog.Infof("NewStorageProvider Get etcd client failed when write instance info.%v\n", err)
		return nil
	}
	sp.client = c
	return &sp
}

func (sp *StorageProvider) Init() error {
	kApi := client.NewKeysAPI(*(sp.client))

	o := client.SetOptions{Dir: true}
	_, err := kApi.Set(context.Background(), "/alibroker", "", &o)
	if err != nil {
		return fmt.Errorf("NewStorageProvider Create etcd alibroker directory failed. %v\n", err)
	}
	return nil
}

//Persistent storage(etcd) api
func (sp *StorageProvider) WriteObject(instanceInfo InstanceRunInfo) {
	key := instanceInfo.InstanceId
	value, err := json.Marshal(instanceInfo)
	if err != nil {
		glog.Infof("Instance Info can not marshal when write object.")
		return
	}

	kApi := client.NewKeysAPI(*(sp.client))

	resp, err := kApi.Set(context.Background(), "/alibroker/"+key, string(value), nil)
	if err != nil {
		glog.Fatal(err)
	} else {
		glog.Infof("Write instance %s info ok.\n", resp)
	}

	return
}

func (sp *StorageProvider) GetObject(instanceID string) (InstanceRunInfo, error) {
	var instanceInfo InstanceRunInfo

	kApi := client.NewKeysAPI(*(sp.client))

	resp, err := kApi.Get(context.Background(), "/alibroker/"+instanceID, nil)
	if err != nil {
		glog.Infof("Get instance %s info failed. error:%v\n", err)
		return instanceInfo, err
	} else {
		glog.Infof("Get instance %s info ok.%v\n", instanceID, resp)
		glog.Infof("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
		err := json.Unmarshal([]byte(resp.Node.Value), &instanceInfo)
		if err != nil {
			glog.Infof("Json unmarshal instance %s failed,error:%v\n.", err)
			return instanceInfo, err
		}
	}
	return instanceInfo, nil
}

func (sp *StorageProvider) DeleteObject(instanceID string) {

	kApi := client.NewKeysAPI(*(sp.client))

	resp, err := kApi.Delete(context.Background(), "/alibroker/"+instanceID, nil)
	if err != nil {
		glog.Infof("Delete instance %s info failed error:%v.\n", err)
	} else {
		glog.Infof("Delete instance %s info ok.\n", resp)
	}
	return
}

func getAllNode(rootNode *client.Node, nodesInfo map[string]string) {
	if !rootNode.Dir {
		nodesInfo[rootNode.Key] = rootNode.Value
		return
	}
	for node := range rootNode.Nodes {
		getAllNode(rootNode.Nodes[node], nodesInfo)
	}
}

func (sp *StorageProvider) GetAllObject() ([]*InstanceRunInfo, error) {
	var instanceInfoSlice []*InstanceRunInfo

	time.Sleep(10 * time.Second)

	kApi := client.NewKeysAPI(*(sp.client))

	resp, err := kApi.Get(context.Background(), "/alibroker", nil)
	if err != nil {
		glog.Infof("Failed to get instance info in stoarge. %v", err.Error())
		return nil, err
	}

	glog.Infof("Get all existed instance info from storage.")

	infoStrMap := make(map[string]string)

	sort.Sort(resp.Node.Nodes)
	for _, n := range resp.Node.Nodes {
		glog.Infof("Key: %q, Value: %q\n", n.Key, n.Value)
		infoStrMap[n.Key] = n.Value
	}

	for _, infoStr := range infoStrMap {
		instanceInfo := new(InstanceRunInfo)
		err := json.Unmarshal([]byte(infoStr), &instanceInfo)
		if err != nil {
			glog.Infof("Failed to JSON unmarshal instance info %v. %v\n.", instanceInfo, err)
			continue
		}
		instanceInfoSlice = append(instanceInfoSlice, instanceInfo)
	}

	return instanceInfoSlice, nil
}
