## Open Service Broker for Alibaba Cloud
This is an implementation of a Open Service Broker API v2.13 to provision and bind service instances of Alibaba Cloud services.
## Prerequisites
* Running Kubernetes cluster >= 1.9
* [Helm](https://github.com/kubernetes/helm) >= v2.7.0
* [Kubernetes Service catalog](https://kubernetes.io/docs/tasks/service-catalog/install-service-catalog-using-helm/) installed >= v0.1.9

## Install
Alibaba Cloud service broker can be installed using the Helm chart in this repository.

```
$ git clone https://github.com/AliyunContainerService/open-service-broker-alibabacloud.git
$ cd open-service-broker-alibabacloud
$ helm install --name alibabacloud-servicebroker --namespace catalog charts/alibabacloud-servicebroker
```
Service broker runs as pod in your Kubernetes cluster. The above Helm chart uses `registry.cn-hangzhou.aliyuncs.com/service-catalog/alibabacloud-servicebroker` as default Docker image repository. 
User can build his own image using the Dockerfile provided in this repo. In that case, user need update `charts/alibabacloud-servicebroker/values.yaml` with his own image info. 

After that, Alibaba Cloud service broker has been installed on Kubernetes cluster, and also registered in service catalog by default.

In case you prefer to register the broker separately later. Just update `registerBroker` to false in `charts/alibabacloud-servicebroker/values.yaml` before installing the chart.
To register the broker, use below template. Do replace value of `spec.url` in `examples/broker.yaml` with the create service path of alibabacloud-servicebroker.
```
$ kubectl create -f examples/broker.yaml
```
If the broker was successfully registered, you can query clusterservicebroker, clusterserviceclass and clusterserviceplan in service catalog now.

```
$ kubectl get clusterservicebroker
NAME         AGE
alibabacloud-servicebroker   1d

$kubectl get clusterserviceclass
NAME                                       AGE
oss-997b8372-8dac-40ac-ae65-758b4a502222   7d
rds-997b8372-8dac-40ac-ae65-758b4a5075a5   7d

$ kubectl get clusterserviceplan
NAME                                       AGE
oss-edc2badc-d93b-4d9c-9d8e-da2f1c8c2221   7d
oss-edc2badc-d93b-4d9c-9d8e-da2f1c8c2222   7d
rds-427559f1-bf2a-45d3-8844-32374a3e58aa   7d
rds-edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1c   7d
rds-edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1d   7d
rds-edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1e   7d
```
So far, alibabacloud-servicebroker is ready to use.

## Usage 

### Prerequisite: setup permission for operating Alibaba Cloud services
One behave of user, service broker creates/deletes/binds/unbinds instances of Alibaba Cloud service. User need assign appropriate permissions to alibabacloud-servicebroker first. 
Alibaba Cloud services leverage [RAM](https://www.aliyun.com/product/ram) policy to enable access control. 
User can config RAM policies for alibabacloud-servicebroker in web console of [Alibaba Cloud Container Service for Kubernetes](https://cs.console.aliyun.com) as below steps.
1.  In web console, select the Kubernetes cluster where alibabacloud-servicebroker is installed. Click on `manage`, go to cluster details page
2.  In `Basic Information` tab, click on 'ROS' link (naming with prefix 'k8s-for-cs-') on `Cluster Resource` panel, go to ROS stack overview page
3.  In `Resource` tab, find resource name `KubernetesWorkerRole`, click on its resource ID link (naming with prefix 'KubernetesWorkerRole-'), go to RAM Role details page
4.  In `Role Authorization Policies` tab, find authorization policy name with prefix 'k8sWorkerRolePolicy-', click `View Permissions`, go to Authorization Policy Details page
5.  On right of `Policy Details` panel, click on `Modify Authorization Policy`
6.  In opened window, edit `Policy Content`. Add RDS and OSS service operation permissions.
E.g. As shown below, user can authorize the 'full' permission to service broker to perform any action to any resource type of RDS and OSS service.

```
    {
      "Action": [
        "rds:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "oss:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    }
```
WARN: The above example policy might be over authorized. It is used only for demonstration. 
      Users are suggested to config detailed policies according to exact requirement. For more details, please refer to authorization policies of [RDS](https://help.aliyun.com/knowledge_detail/58932.html) and [OSS](https://help.aliyun.com/knowledge_detail/58905.html) .

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
Notice: Please set username and password for your binding in `examples/rds-binding.yaml`. They are used to access the previous RDS instance.
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
Then user's application pod can connect to RDS instance `rds-mysql-1` by using this secret.