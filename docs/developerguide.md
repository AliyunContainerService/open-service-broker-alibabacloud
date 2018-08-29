##Overview
This repository is organized as similarly to Kubernetes. Below is a summary of the repository's layout:
```
├── charts                  # Helm charts
│   └── alibabacloud-servicebroker          # Helm chart for deploying the Alibaba Cloud service broker
├── pkg                     # Contains all non-"main" Go packages
│   └── server              # The basic web server of Alibaba cloud service broker
│   └── brokerapi           # The interface defined for Open Service Broker API
│   └── controller          # Controller implements the generic logic of Alibaba Cloud service broker.
│   │                       # All concrete brokers just handel service specific logic, and are registered to the controller.
│   └── services            # The cloud services supported by the Alibaba Cloud service broker, e.g. RDS and OSS
│   └── util                # The common functions
├── docs                    # Documentation
├── vendor                  # Dependencies
├── main.go                 # Main
├── Dockerfile              # Dockerfile for building Alibaba Cloud service broker Docker image
│...

```
##Prerequisites
At a minimum, to develop service broker you will need:

* Docker installed locally
* git
 
These will allow you to build and test service broker within a Docker container.

If you want to deploy service broker to Kubernetes cluster manually, you will also need:

* A working Kubernetes cluster and kubectl installed in your local PATH, properly configured to access that cluster. The version of Kubernetes and kubectl must be >= 1.9.
* Helm (Tiller) installed in your Kubernetes cluster and the helm binary located in your local PATH
* To be pre-authenticated to a Docker registry (if using a remote cluster)

##Build Alibaba Cloud Service Broker
###Download source code

```
$ git clone https://github.com/AliyunContainerService/open-service-broker-alibabacloud.git
```
###Build Docker image

```
$ cd open-service-broker-alibabacloud
$ docker build --tag alibabacloud-servicebroker:xxx .
```
After building Docker image successfully, push it to your public accessiable Docker registry.

##Deploy Alibaba Cloud Service Broker

* Modify values.yaml in charts/alibabacloud-servicebroker folder, to use your own Docker images location.
* Use helm to install Alibaba Cloud service broker into Kubernetes cluster.

```
$ cd open-service-broker-alibabacloud
$ helm install
```

##Extend Alibaba Cloud Service Broker to support new service

To add new service implementation:
* create a new folder in pkg/services to locate the broker implementation of new service
* implement ServiceBroker interface defined in pkg/brokerapi/service_broker.go
* register new service broker in function registerBrokers() in pkg/server/server.go, like this

```
	// register broker for specific cloud services
	brokers["rds-broker"] = rds.CreateBroker()
```
##Test
