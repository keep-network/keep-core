#!/bin/sh
kubectl create -f infrastructure/kube/lcl/dashboard.yaml
kubectl create -f infrastructure/kube/lcl/miner-nodes.yaml
kubectl create -f infrastructure/kube/lcl/tx-nodes.yaml
