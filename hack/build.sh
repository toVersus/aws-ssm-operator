#!/bin/bash

export GO111MODULE=on
go mod vendor
operator-sdk build toversus/aws-ssm-operator:v0.0.1
docker push toversus/aws-ssm-operator:v0.0.1
