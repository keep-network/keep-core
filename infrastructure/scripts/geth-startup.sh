#!/bin/bash
DIRNAME=`dirname $0`
kubectl create -f $DIRNAME/../kube/lcl/dashboard.yaml
kubectl create -f $DIRNAME/../kube/lcl/miner-nodes.yaml
kubectl create -f $DIRNAME/../kube/lcl/tx-nodes.yaml
