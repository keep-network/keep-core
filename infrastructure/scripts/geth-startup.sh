#!/bin/sh
kubectl create -f infrastructure/kube/dashboard.yaml
kubectl create -f infrastructure/kube/miner-nodes.yaml
kubectl create -f infrastructure/kube/tx-nodes.yaml
