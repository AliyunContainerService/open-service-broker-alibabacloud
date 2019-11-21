package server

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/brokerapi"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/controller"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/services/oss"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/services/polardb"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/services/rds"
	//"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/services/ots"
	"github.com/AliyunContainerService/open-service-broker-alibabacloud/pkg/util"
	"github.com/golang/glog"

	"fmt"
	"github.com/gorilla/mux"
)

var Options struct {
	Port int
	RunServices string
}

type server struct {
	controller Controller
}

// CreateHandler creates Broker HTTP handler based on an implementation
// of a controller.Controller interface.
func createHandler(c Controller) http.Handler {
	s := server{
		controller: c,
	}
	var router = mux.NewRouter()
	router.HandleFunc(
		"/v2/catalog", s.catalog,
	).Methods("GET")
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/last_operation", s.getServiceInstanceLastOperation,
	).Methods("GET")
	router.HandleFunc(
		"/v2/service_instances/{instance_id}", s.createServiceInstance,
	).Methods("PUT")
	router.HandleFunc(
		"/v2/service_instances/{instance_id}", s.removeServiceInstance,
	).Methods("DELETE")
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.bind,
	).Methods("PUT")
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}", s.unBind,
	).Methods("DELETE")
	router.HandleFunc(
		"/v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation", s.pollBindingLastOperation,
	).Methods("GET")
	return router
}

// Run creates the HTTP handler based on an implementation of a
// controller.Controller interface, and begins to listen on the specified address.
func Run(ctx context.Context, addr string, c Controller) error {
	listenAndServe := func(srv *http.Server) error {
		return srv.ListenAndServe()
	}
	return run(ctx, addr, listenAndServe, c)
}

// RunTLS creates the HTTPS handler based on an implementation of a
// controller.Controller interface, and begins to listen on the specified address.
func RunTLS(ctx context.Context, addr string, cert string, key string, c Controller) error {
	var decodedCert, decodedKey []byte
	var tlsCert tls.Certificate
	var err error
	decodedCert, err = base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return err
	}
	decodedKey, err = base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}
	tlsCert, err = tls.X509KeyPair(decodedCert, decodedKey)
	if err != nil {
		return err
	}
	listenAndServe := func(srv *http.Server) error {
		srv.TLSConfig = new(tls.Config)
		srv.TLSConfig.Certificates = []tls.Certificate{tlsCert}
		return srv.ListenAndServeTLS("", "")
	}
	return run(ctx, addr, listenAndServe, c)
}

type BrokerFactoryFunc func() brokerapi.ServiceBroker

type BrokerFactory struct {
	all map[string]BrokerFactoryFunc
}

func BrokerFactoryInit() *BrokerFactory {
	return &BrokerFactory{
		all: map[string]BrokerFactoryFunc{
			"oss-broker":     oss.CreateBroker,
			"rds-broker":     rds.CreateBroker,
			"polardb-broker": polardb.CreateBroker,
		},
	}
}

func (b *BrokerFactory) createByName(svrName string) (brokerapi.ServiceBroker, error) {
	brokerFunc, ok := b.all[svrName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("serviceName: %v not exists", svrName))
	} else {
		glog.Infof("service broker : %v created", svrName)
		return brokerFunc(), nil
	}
}

func (b *BrokerFactory) createAll() map[string]brokerapi.ServiceBroker {
	brokers := map[string]brokerapi.ServiceBroker{}
	for svrName, brokerFunc := range b.all {
		brokers[svrName] = brokerFunc()
	}
	return brokers
}

func registerBrokers() error {
	var registered bool = false
	glog.Info("Start to register service brokers")

	brokers := make(map[string]brokerapi.ServiceBroker)
	if Options.RunServices != "" && Options.RunServices != "all" {
		runServices := strings.Split(Options.RunServices, ",")
		for _, runSvr := range runServices {
			broker, err := BrokerFactoryInit().createByName(runSvr)
			if err != nil {
				return err
			}
			brokers[runSvr] = broker
		}
	} else {
		brokers = BrokerFactoryInit().createAll()
	}

	registered = controller.RegisterBrokers(brokers)
	if !registered {
		return fmt.Errorf("None broker can be registered for any cloud service.")
	}
	glog.Info("Completed to register service brokers.\n")
	return nil
}

func run(ctx context.Context, addr string, listenAndServe func(srv *http.Server) error, c Controller) error {
	glog.Infof("Starting service broker server on %s\n", addr)
	NewJobQueue()
	dispatcher := NewDispatcher()
	dispatcher.Run()

	err := registerBrokers()
	if err != nil {
		glog.Error(err.Error())
		return fmt.Errorf("No service broker available.")
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: createHandler(c),
	}
	go func() {
		<-ctx.Done()
		c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if srv.Shutdown(c) != nil {
			srv.Close()
		}
	}()
	return listenAndServe(srv)
}

func (s *server) catalog(w http.ResponseWriter, r *http.Request) {
	//glog.Infof("Get Service Broker Catalog...")

	if result, err := s.controller.GetCatalog(); err == nil {
		util.WriteResponse(w, http.StatusOK, result)
	} else {
		util.WriteErrorResponse(w, http.StatusBadRequest, err)
	}
}

func (s *server) getServiceInstanceLastOperation(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]
	q := r.URL.Query()
	serviceID := q.Get("service_id")
	planID := q.Get("plan_id")
	operation := q.Get("operation")
	glog.Infof("getServiceInstanceLastOperation ... %s\n", instanceID)

	var req brokerapi.LastOperationRequest

	req.Operation = operation
	req.PlanID = planID
	req.ServiceID = serviceID

	sChan := make(chan brokerapi.WorkerResponse)
	instanceLastOperationPayload := controller.InstanceLastOperationPayload{
		InstanceID: instanceID,
		Request:    req,
		Result:     sChan}

	EnqueueJob(Job{Payload: instanceLastOperationPayload})
	response := <-sChan
	glog.Infof("Service broker server getServiceInstanceLastOperation get response:%v.\n", response)
	if response.Err == nil {
		if state, ok := response.Message.(string); ok {
			result := &brokerapi.LastOperationResponse{State: state}
			util.WriteResponse(w, http.StatusOK, result)
			glog.Infof("getServiceInstanceLastOperation get result:%v.\n", result)
			return
		}
	}

	util.WriteErrorResponse(w, http.StatusBadRequest, response.Err)
	return
}

func (s *server) createServiceInstance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["instance_id"]
	var req brokerapi.CreateServiceInstanceRequest
	if err := util.BodyToObject(r, &req); err != nil {
		glog.Errorf("Failed to unmarshal CreateServiceInstanceRequest: %v", err)
		util.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	// TODO: Check if parameters are required, if not, this thing below is ok to leave in,
	// if they are ,they should be checked. Because if no parameters are passed in, this will
	// fail because req.Parameters is nil.
	if req.Parameters == nil {
		req.Parameters = make(map[string]interface{})
	}
	glog.Infof("CreateServiceInstance %s, with request:\n%v", id, r)
	sChan := make(chan brokerapi.WorkerResponse)
	createServiceInstancePayload := controller.CreateServiceInstancePayload{
		InstanceID: id,
		Request:    req,
		Result:     sChan}

	EnqueueJob(Job{Payload: createServiceInstancePayload})
	glog.Infof("CreateServiceInstance wait for response.\n")
	response := <-sChan
	glog.Infof("Service broker server createServiceInstance get response:%v.\n", response)
	if response.Err == nil {
		result := &brokerapi.CreateServiceInstanceResponse{Operation: response.Message.(string)}
		util.WriteResponse(w, http.StatusAccepted, result)
	} else {
		util.WriteErrorResponse(w, http.StatusBadRequest, response.Err)
	}
}

func (s *server) removeServiceInstance(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]
	q := r.URL.Query()
	serviceID := q.Get("service_id")
	planID := q.Get("plan_id")
	acceptsIncomplete := q.Get("accepts_incomplete") == "true"
	glog.Infof("RemoveServiceInstance %s...\n", instanceID)

	sChan := make(chan brokerapi.WorkerResponse)
	deleteServiceInstancePayload := controller.DeleteServiceInstancePayload{
		InstanceID: instanceID,
		Request:    brokerapi.DeleteServiceInstanceRequest{ServiceID: serviceID, PlanID: planID, AcceptsIncomplete: acceptsIncomplete},
		Result:     sChan}

	EnqueueJob(Job{Payload: deleteServiceInstancePayload})
	response := <-sChan
	glog.Infof("Service broker server removeServiceInstance get response:%v.\n", response)
	if response.Err == nil {
		util.WriteResponse(w, http.StatusOK, &brokerapi.DeleteServiceInstanceResponse{})
	} else {
		util.WriteErrorResponse(w, http.StatusBadRequest, response.Err)
	}
}

func (s *server) bind(w http.ResponseWriter, r *http.Request) {
	bindingID := mux.Vars(r)["binding_id"]
	instanceID := mux.Vars(r)["instance_id"]
	q := r.URL.Query()

	acceptsIncomplete := q.Get("accepts_incomplete") == "true"

	glog.Infof("Bind binding_id=%s, instance_id=%s acceptsIncomplete:%v request:%v\n", bindingID, instanceID, acceptsIncomplete, r)

	var req brokerapi.BindingRequest

	if err := util.BodyToObject(r, &req); err != nil {
		glog.Errorf("Failed to unmarshall request: %v", err)
		util.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	// TODO: Check if parameters are required, if not, this thing below is ok to leave in,
	// if they are ,they should be checked. Because if no parameters are passed in, this will
	// fail because req.Parameters is nil.
	if req.Parameters == nil {
		req.Parameters = make(map[string]interface{})
	}

	// Pass in the instanceId to the template.
	req.Parameters["instanceId"] = instanceID

	sChan := make(chan brokerapi.WorkerResponse)
	bindingPayload := controller.BindingPayload{
		InstanceID: instanceID,
		BindingID:  bindingID,
		Request:    req,
		Result:     sChan}

	EnqueueJob(Job{Payload: bindingPayload})

	glog.Infof("Bind binding_id=%s for instance_id=%s wait response\n", bindingID, instanceID)
	response := <-sChan
	glog.Infof("Service broker server bind binding_id=%s for instance_id=%s get response:%v\n", bindingID, instanceID, response)
	if response.Err == nil {
		credential, ok := (response.Message).(brokerapi.Credential)
		if !ok {
			return
		}
		result := &brokerapi.CreateServiceBindingResponse{Credentials: credential}
		util.WriteResponse(w, http.StatusOK, result)
	} else {
		util.WriteErrorResponse(w, http.StatusBadRequest, response.Err)
	}
}

func (s *server) unBind(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]
	bindingID := mux.Vars(r)["binding_id"]
	q := r.URL.Query()
	serviceID := q.Get("service_id")
	planID := q.Get("plan_id")
	glog.Infof("UnBind: Service instance guid: %s:%s", bindingID, instanceID)

	sChan := make(chan brokerapi.WorkerResponse)
	unBindingPayload := controller.UnBindingPayload{
		InstanceID: instanceID,
		BindingID:  bindingID,
		ServiceID:  serviceID,
		PlanID:     planID,
		Result:     sChan}

	EnqueueJob(Job{Payload: unBindingPayload})
	response := <-sChan
	glog.Infof("Service broker server unbind binding_id=%s for instance_id=%s get response:%v\n", bindingID, instanceID, response)
	if response.Err == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{}") //id)
	} else {
		util.WriteErrorResponse(w, http.StatusBadRequest, response.Err)
	}
}

func (s *server) pollBindingLastOperation(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]
	bindingID := mux.Vars(r)["binding_id"]
	q := r.URL.Query()
	serviceID := q.Get("service_id")
	planID := q.Get("plan_id")
	operation := q.Get("operation")

	var req brokerapi.BindingLastOperationRequest
	if err := util.BodyToObject(r, &req); err != nil {
		glog.Errorf("error unmarshalling: %v", err)
		util.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	glog.Infof("pollBindingLastOperation: Service instance guid: bindingID:%s instanceID:%s serviceID:%s planID:%s :%s",
		bindingID, instanceID, serviceID, planID, operation)

	sChan := make(chan brokerapi.WorkerResponse)
	bindingLastOperationPayload := controller.BindingLastOperationPayload{
		InstanceID: instanceID,
		BindingID:  bindingID,
		Request:    req,
		Result:     sChan}

	EnqueueJob(Job{Payload: bindingLastOperationPayload})
	response := <-sChan
	if response.Err == nil {
		util.WriteResponse(w, http.StatusOK, response.Message)
	} else {
		util.WriteErrorResponse(w, http.StatusBadRequest, response.Err)
	}
	return
}
