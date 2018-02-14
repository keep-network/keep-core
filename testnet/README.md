# Private Ethereum on Kubernetes cluster
## Introduction
According to the official docs ([go-ethereum/wiki](https://github.com/ethereum/go-ethereum/wiki/Private-network), [ethdocs.org](http://ethdocs.org/en/latest/network/test-networks.html)) in order to set up a Private Ethereum Network we need to configure the following:

#### Bootnode
You need to start a bootstrap node that others can use to find each other in
your network. `go-ethereum` offers a bootnode implementation that can be
configured and run in your private network.

#### Miner
A single CPU miner instance is more than enough for practical purposes as it can
produce a stable stream of blocks at the correct intervals without needing heavy
resources (consider running on a single thread, no need for multiple ones
either)

#### Non default Network ID
The main Ethereum network has id 1 (the default). So if you supply your own
custom network ID which is different than the main network your nodes will not
connect to other nodes and form a private network.

#### Custom genesis block
JSON file where you setup Network ID (chainID), configure network settings and
allocate pre-funds accounts. You need custom genesis block because if someone
accidentally connects to your testnet using the real chain, your local copy will
be considered a stale fork and updated to the "real" one.

#### Custom Data Directory
Choose a location that is separate from your public Ethereum chain folder,
otherwise, in order to successfully mine a block, you would need to mine against
the difficulty of the last block present in your local copy of the blockchain -
which may take several hours.

## Kubernetes checklist
Make sure you authenticated to use your cluster and the right context is set:

```
kubectl config get-contexts

... list of contexts

kubectl config use-context name_of_your_cluster
```

To swich namespace, use the following:
```
kubectl config set-context $(kubectl config current-context) --namespace=<namespace_name>
```

## Intro to HELM

Rather than creating YAML files with hardcoded values we are using [Helm](https://helm.sh/)
package manager that can generate them based on the provided variables.

### Installation
```
brew install kubernetes-helm
```
### Basic usage
```
cd /project_folder

# deploy project release
helm install --name <release_name> .

# list all releases
helm list

# check status
helm status <release_name>

# update release
helm upgrade <release_name> .

# delete release
helm del --purge <release_name>
```

### Variables
The variables are set in **values.yaml**


## Deployment
### TL;DR

```
helm install --name my-testnet-v1 .
```
Install private ethereum network with one miner node.

### Notice
It can take some time (20-30min) for the miner node to generate DAG file before
it starts mining. You can check the progress with the following command:

```
kubectl logs <pod_name_with_miner> | grep "DAG"

> INFO [...] Generating DAG in progress epoch=0 percentage=28 elapsed=7m46.574s
```

## Deployment step-by-step guide
### 1. Bootnode

To [setup a bootnode](https://github.com/ethereum/go-ethereum/wiki/Private-network#network-connectivity) containers in this deployment will make sure to:

* Run bootnode once to generate a key `bootnode --genkey=boot.key`
* Start bootnode with this key `bootnode --nodekey=boot.key`
* Grab the returned **enode URL** and make sure other nodes use that URL

#### bootnode.deployment.yaml 

###### Initial variables (values.yaml)

| Parameter                  | Description                        | Default                                                    |
| -----------------------    | ---------------------------------- | ---------------------------------------------------------- |
| `replicaCount`    | Number of replicas  | 1
| `image.repository` | `geth` image   | ethereum/client-go
| `image.tag` | `geth` image tag | alltools-v1.7.3 (Geth + Tools)


###### List of containers created by this deployment:

* [init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) to generate bootnode key in advance
* Bootnode main container
* Contaner to broadcast enode URL via netcast. So it's available accross all the
  pods in the cluster.

>Note on broadcasting enode URL. There are other possible ways of doing this
>i.e. using daemonSet or StatefulSet. Please share your knowledge if you think
>of a better solution.

#### bootnode.service.yaml

Making bootnode available at port 30301 and netcat broadcasting bootnode enode
URL at port 80

### 2. Node 
For the [initial setup](https://github.com/ethereum/go-ethereum/wiki/Private-network#creating-the-genesis-block) of the network we need to provide:

* Genesis block
* Path to a data folder

For a subsequent [miner setup](https://github.com/ethereum/go-ethereum/wiki/Private-network#running-a-private-miner) we need to provide:

* Bootnode enode URL
* Network ID
* Path to a data folder

#### genesis.json
[Genesis block](https://github.com/ethereum/go-ethereum/wiki/Private-network#creating-the-genesis-block) contains your network settings, 
initial accounts and allocated funds.

```javascript
{
  "config": {
    "chainId": 98052,
    "homesteadBlock": 0,
    "eip155Block": 0,
    "eip158Block": 0
  },
  "difficulty" : "0x20000",
  "gasLimit"   : "0x493E00",
  "alloc": {
    "0x2932b7a2355d6fecc4b5c0b6bd44cc31df247a2e": {
      "balance": "1000000000000000000000"
    }
  }
}
```

#### miner.deployment.yaml 
###### Initial variables (values.yaml)

| Parameter                  | Description                        | Default                                                    |
| -----------------------    | ---------------------------------- | ---------------------------------------------------------- |
| `replicaCount`    | Number of replicas  | 1
| `image.repository` | `geth` image   | ethereum/client-go
| `image.tag` | `geth` image tag | v1.7.3
| `networkId`    | Network ID  | 1101
| `initialAccounts` | Initial accounts |   

###### List of containers created by this deployment:

* [init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) to setup genesis block.
* [init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) to fetch bootnode enode URL.
* Miner main container

#### miner.service.yaml 
Making jsonrpc an ipc ports available.

### 3. Monitor

#### monitor.deployment.yaml 
This is a straightforward setup using [docker image](https://hub.docker.com/r/ethereumex/eth-stats-dashboard/) of the popular 
Ethereum Network Stats dashboard https://github.com/cubedro/eth-netstats

#### monitor.service.yaml
Making dashboard available at port 80


### 4. VPN  (optional)
There is a great [HELM package](https://github.com/kubernetes/charts/tree/master/stable/openvpn
)  that will install openVPN on your cluster. 

>The chart will automatically configure dns to use kube-dns and route all
>network traffic to kubernetes pods and services through the vpn. By connecting
>to this vpn a host is effectively inside a cluster's network.

#### Installation

```
helm repo add stable http://storage.googleapis.com/kubernetes-charts
helm install stable/openvpn
```

#### Create client

As seen on [https://github.com/kubernetes/charts/tree/master/stable/openvpn](https://github.com/kubernetes/charts/tree/master/stable/openvpn)

```bash
#!/bin/bash

if [ $# -ne 1 ]
then
  echo "Usage: $0 <CLIENT_KEY_NAME>"
  exit
fi

KEY_NAME=$1
NAMESPACE=$(kubectl get pods --all-namespaces -l type=openvpn -o jsonpath='{.items[0].metadata.namespace}')
POD_NAME=$(kubectl get pods -n $NAMESPACE -l type=openvpn -o jsonpath='{.items[0].metadata.name}')
SERVICE_NAME=$(kubectl get svc -n $NAMESPACE -l type=openvpn  -o jsonpath='{.items[0].metadata.name}')
SERVICE_IP=$(kubectl get svc -n $NAMESPACE $SERVICE_NAME -o go-template='{{range $k, $v := (index .status.loadBalancer.ingress 0)}}{{$v}}{{end}}')
kubectl -n $NAMESPACE exec -it $POD_NAME /etc/openvpn/setup/newClientCert.sh $KEY_NAME $SERVICE_IP
kubectl -n $NAMESPACE exec -it $POD_NAME cat /etc/openvpn/certs/pki/$KEY_NAME.ovpn > $KEY_NAME.ovpn


```
Example usage: ``` the_script_above.sh <CLIENT_KEY_NAME>```

#### Usage
Use the generated config file with VPN client of your choice, i.e. [Tunnelblick](https://tunnelblick.net/downloads.html). 
Once connected, you should be able to access your cluster services via internal IPs.

### 5. Nginx Ingress controller (Optional)
This will allow you to have useful features such as virtual hosts, TLS,
whitelisted source IP range 

#### Install ingress-nginx
Install kubernetes [ingress-nginx](https://github.com/kubernetes/ingress-nginx) via [this HELM package](https://github.com/kubernetes/charts/tree/master/stable/nginx-ingress) with the following command:
	
```console
helm install stable/nginx-ingress --name <RELEASE_NAME>
```

Example with limiting connection only from specific IP range:
	
```console
helm install stable/nginx-ingress --name <RELEASE_NAME> \ 
--set controller.service.loadBalancerSourceRanges={<IP_ADDRESS>/32}
```

#### Create self-signed TLS certificate

Notice: make sure to provide the name of your host when asked for *Common Name*.
You can use wildcard for common name as well ***.example.com**


```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls-key.key -out tls-testnet.crt

kubectl create secret tls keep-tls-testnet --key tls-key.key --cert tls-testnet.crt
```

#### Create ingress resource

Example of multi host ingress usinf TLS certificate: 

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/tls-acme: "true"
    kubernetes.io/ingress.class: nginx
    kubernetes.io/ssl-redirect: "true"
  name: testnet-ingress
spec:
  tls:
  - hosts:
    - app1.example.com
    - app2.example.com
    secretName: tls-testnet
  rules:
    - host: app1.example.com
      http:
        paths:
        - path: /
          backend:
            serviceName: app1-svc
            servicePort: 8545
    - host: app2.example.com
      http:
        paths:
        - path: /
          backend:
            serviceName: app2-svc
            servicePort: 3001

```

## Credits
Big thanks to the following github contributors for tips and inspiration:

* [https://github.com/kuberstack/charts/tree/master/incubator/geth](https://github.com/kuberstack/charts/tree/master/incubator/geth)
* [https://github.com/jpoon/kubernetes-ethereum-chart](https://github.com/jpoon/kubernetes-ethereum-chart)
* [https://github.com/MaximilianMeister/kuberneteth](https://github.com/MaximilianMeister/kuberneteth)
* [https://github.com/wbuchwalter/ethereum-kubernetes-cluster](https://github.com/wbuchwalter/ethereum-kubernetes-cluster)
* [https://github.com/rehive/geth-helm-chart/tree/master/geth-chart](https://github.com/rehive/geth-helm-chart/tree/master/geth-chart)
* [https://github.com/skozin/ethereum-dev-net](https://github.com/skozin/ethereum-dev-net)