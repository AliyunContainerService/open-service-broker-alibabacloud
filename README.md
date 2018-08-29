[![Build Status](https://travis-ci.org/AliyunContainerService/open-service-broker-alibabacloud.svg?branch=master)](https://travis-ci.org/AliyunContainerService/open-service-broker-alibabacloud)
[![CircleCI](https://circleci.com/gh/AliyunContainerService/open-service-broker-alibabacloud.svg?style=svg)](https://circleci.com/gh/AliyunContainerService/open-service-broker-alibabacloud)
[![Go Report Card](https://goreportcard.com/badge/github.com/AliyunContainerService/open-service-broker-alibabacloud)](https://goreportcard.com/report/github.com/AliyunContainerService/open-service-broker-alibabacloud)

# Open Service Broker for Alibaba Cloud

The open-service-broker-alibabacloud is an implementation of [Open Service Broker (OSB) API](https://github.com/openservicebrokerapi/servicebroker/blob/v2.13/spec.md) for Alibaba Cloud services.
So far Alibaba Cloud RDS (Relational Database Service) and OSS (Object Store Service) are included. More services will be supported soon.

open-service-broker-alibabacloud is now able to run on Kubernetes. It works together with [Kubernetes Service Catalog](https://github.com/kubernetes-incubator/service-catalog)

[Learn more about the Kubernetes Service Catalog](https://svc-cat.io/) and its great [Walkthrough](https://svc-cat.io/docs/walkthrough/)

## Prerequisites

1. Kubernetes cluster
2. [Helm](https://github.com/kubernetes/helm)
3. [Service Catalog](https://github.com/kubernetes-incubator/service-catalog) - follow the [walkthrough](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/walkthrough.md)

## How to use

Please refer to specific documentation for details.

Read [design document](docs/design.md) to understand how it's designed and how it works;
Read [user guide document](docs/userguide.md) to learn how to use it in your Kubernetes;
Read [developer guide document](docs/developerguide.md) to learn how to extend open-service-broker-alibabacloud by adding a new implementation for specific cloud service
