##Overview
The open-service-broker-alibabacloud is an implementation of Open Service Broker (OSB) API v2.13 for Alibaba Cloud services.
So far Alibaba Cloud RDS (Relational Database Service) and OSS (Object Store Service) are included. More services will be supported soon.
The open-service-broker-alibabacloud is designed in an extensible way.
In general, there is a base controller, and several service brokers registered to the controller.

The base controller implements most of generic works including
* a service broker server to serving broker API routes
* a job queue mechanism to handle requests asynchronously (because cloud service provider may create instance slowly, it's better to handle requests asynchronously)
* a simple storage to persistent all service instance and binding information, which can recover server back to correct status
* register all available brokers for specific cloud services, and route particular request to specific broker
* update status of a request flow according to the progress of backend service provider

With such architecture, it has following advantages:
* new service can easily been added
* developer needn't pay too much attention to interaction process of OSB server, but just focus on the service specific implementation align with OSB API spec, like provision and bind service instance.

The broker has 3 internal modules:
* Dispatcher
* Async Engine
* Persistent Storage

## Dispatcher

The dispatcher is designed to accept highly concurrent requests and ensure stable and controllable resource consumption.
A job queue to cache the request jobs.
A worker pool to provide workers to deal with request.
The dispatcher watches the job channel and dispatch the new coming job to some worker's channel.

## Async Engine

Async Engine is designed to interact with the cloud service to get the latest status of service instance and bindings.
The async engine scan progress of backend service provider processing the request of instance and binding.
If they are still working in process, engine will update the latest status of specific service instance and binding.
The base controller response to the query of "last_operation" according to those status.

## Persistent Storage

Persistent storage is designed to save service instance and binding information. The base controller reads historical information from storage when starting.
Service instance information is kept in a map, and every update will be persistent into storage.
Now the storage is implemented with an embedded Etcd.
