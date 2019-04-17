#!/bin/bash

kubectl create -f ../kube/keep-dev/keep-client/keep-client-bootstrap-peer-0.yaml
kubectl create -f ../kube/keep-dev/keep-client/keep-client-peer-0.yaml
kubectl create -f ../kube/keep-dev/keep-client/keep-client-peer-1.yaml
kubectl create -f ../kube/keep-dev/keep-client/keep-client-peer-2.yaml
kubectl create -f ../kube/keep-dev/keep-client/keep-client-peer-3.yaml
