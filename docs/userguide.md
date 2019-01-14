## Open Service Broker for Alibaba Cloud
This is an implementation of a Open Service Broker API v2.13 to provision and bind service instances of Alibaba Cloud services.
## Prerequisites
* Running Kubernetes cluster >= 1.9
* [Helm](https://github.com/kubernetes/helm) >= v2.7.0
* [Kubernetes Service catalog](https://kubernetes.io/docs/tasks/service-catalog/install-service-catalog-using-helm/) installed >= v0.1.9

## Installing the Broker
The Alibaba Cloud service broker can be installed using the Helm chart in this repository.

```
$ git clone https://github.com/AliyunContainerService/open-service-broker-alibabacloud.git
$ cd open-service-broker-alibabacloud
$ helm install --name alibabacloud-servicebroker --namespace catalog charts/alibabacloud-servicebroker
```
After that, Alibaba Cloud service broker has been installed on Kubernetes cluster, and also registered in service catalog by default.

In case you prefer to register the broker separately later. Just update `registerBroker` to false in charts/alibabacloud-servicebroker/values.yaml before installing the chart.
To register the broker, use below template. Do replace value of 'spec.url' in examples/broker.yaml with the create service path of alibabacloud-servicebroker.
```
$ kubectl create -f examples/broker.yaml
```
If the broker was successfully registered, you can query clusterservicebroker, clusterserviceclass and clusterserviceplan in service catalog now.

```
$ kubectl get clusterservicebroker
NAME         AGE
alibabacloud-servicebroker   1d
$kubectl get clusterserviceclass
NAME                                   AGE
997b8372-8dac-40ac-ae65-758b4a502222   1d
997b8372-8dac-40ac-ae65-758b4a5075a5   1d
$ kubectl get clusterserviceplan
NAME                                   AGE
427559f1-bf2a-45d3-8844-32374a3e58aa   1d
edc2badc-d93b-4d9c-9d8e-da2f1c8c2221   1d
edc2badc-d93b-4d9c-9d8e-da2f1c8c2222   1d
edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1c   1d
edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1d   1d
edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1e   1d
```
So far, alibabacloud-servicebroker is ready to use.

## Usage 

### Create service instance

```
$ kubectl create -f examples/rds-mysql-instance.yaml
```
This will result in a service instance of Alibaba Cloud RDS service with MySQL engine:
```
$ kubectl get serviceinstance
NAME            AGE
rds-mysql-1     5s
```
### Create service binding
```
kubectl create -f examples/rds-binding.yaml
```
Notice: Please set username and password for your binding in examples/rds-binding.yaml. They are used to access the previous RDS instance.
This will result in a service binding for the newly created RDS service instance.
```
$ kubectl get servicebindings
NAME           AGE
rds-binding    1m
```
A secret object is created to store above credential in the cluster.

```
$ kubectl get secret
NAME                        TYPE                                  DATA      AGE
default-token-5wwkc         kubernetes.io/service-account-token   3         36d
rds-instance-credentials    Opaque                                6         5m

$ kubectl describe  secret rds-instance-credentials
Name:         rds-instance-credentials
Namespace:    default
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
username:           44 bytes
password:           16 bytes
```
