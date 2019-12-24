#!/bin/bash

export GO111MODULE=on
go mod vendor
operator-sdk build toversus/aws-ssm-operator:v0.2.1
docker push toversus/aws-ssm-operator:v0.2.1
