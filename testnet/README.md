
#Private Ethereum on Kubernetes cluster

## Introduction
According to the [official docs](https://github.com/ethereum/go-ethereum/wiki/Private-network) to set up a Private Ethereum Network we need to configure the following:

#### Bootnode
You need to start a bootstrap node that others can use to find each other in your network. `go-ethereum` offers a bootnode implementation that can be configured and run in your private network.

#### Miner
A single CPU miner instance is more than enough for practical purposes as it can produce a stable stream of blocks at the correct intervals without needing heavy resources (consider running on a single thread, no need for multiple ones either)

#### Non default Network ID
The main Ethereum network has id 1 (the default). So if you supply your own custom network ID which is different than the main network your nodes will not connect to other nodes and form a private network.

#### Genesis Block
JSON file where you setup Network ID (chainID), configure network settings and allocate pre-funds accounts.

## Kubernetes checklist
Make sure you authenticated to use your cluster and the right context is set

```
kubectl config get-contexts

... list of contexts

kubectl config use-context name_of_your_cluster

```
To swich namespace, use the following

```
kubectl config set-context $(kubectl config current-context) --namespace=<namespace_name>

```


## Intro to HELM

Rather than creating YAML files with hardcoded values we're gonna use [Helm](https://helm.sh/) package manager that can generate them based on the provided variables.

### Installation
```
brew install kubernetes-helm
```
### Basic usage
```
cd /project_folder

# deploy project
helm install --name project .

# check status
helm status project

# update project
helm upgrade project .

# delete project
helm del --purge project
```

### Variables
The variables are set in **values.yaml**



## Deployment: Bootnode

To [setup a bootnode](https://github.com/ethereum/go-ethereum/wiki/Private-network#network-connectivity) we need to:

* Run bootnode once to generate a key `bootnode --genkey=boot.key`
* Start bootnode with this key `bootnode --nodekey=boot.key`
* Grab the returned **enode URL** and make sure other nodes use that URL

###bootnode.deployment.yaml 
######Initial variables (values.yaml)
| Parameter                  | Description                        | Default                                                    |
| -----------------------    | ---------------------------------- | ---------------------------------------------------------- |
| `replicaCount`    | Number of replicas  | 1
| `image.repository` | `geth` image   | ethereum/client-go
| `image.tag` | `geth` image tag | alltools-v1.7.3 (Geth + Tools)                                       

######Containers

* [init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) to generate bootnode key in advance
* Bootnode main container
* Contaner to broadcast enode URL via netcast. So it's available accross all the pods in the cluster.

>Note on broadcasting enode URL. There are other possible ways of doing this i.e. using daemonSet or StatefulSet. Please share your knowledge if you think of a better solution.

###bootnode.service.yaml

Making bootnode available at port 30301 and netcat broadcasting bootnode enode URL at port 80



## Deployment: Miner
For the [initial setup](https://github.com/ethereum/go-ethereum/wiki/Private-network#creating-the-genesis-block) of the network we need to provide:

* Genesis block
* Path to a data folder

For a subsequent [miner setup](https://github.com/ethereum/go-ethereum/wiki/Private-network#running-a-private-miner) we need to provide:

* Bootnode enode URL
* Network ID
* Path to a data folder

###genesis.json
[Genesis block](https://github.com/ethereum/go-ethereum/wiki/Private-network#creating-the-genesis-block) contains your network settings, initial accounts and allocated funds.

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

###miner.deployment.yaml 
######Initial variables (values.yaml)

| Parameter                  | Description                        | Default                                                    |
| -----------------------    | ---------------------------------- | ---------------------------------------------------------- |
| `replicaCount`    | Number of replicas  | 1
| `image.repository` | `geth` image   | ethereum/client-go
| `image.tag` | `geth` image tag | v1.7.3
| `networkId`    | Network ID  | 1101
| `initialAccounts` | Initial accounts |   

######Containers

* [init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) to setup genesis block.
* [init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) to fetch bootnode enode URL.
* Miner main container

###miner.service.yaml 
Making jsonrpc an ipc ports available.

## Deployment: Monitor

###monitor.deployment.yaml 
This is a straightforward setup using [docker image](https://hub.docker.com/r/ethereumex/eth-stats-dashboard/) of the popular Ethereum Network Stats dashboard https://github.com/cubedro/eth-netstats

###monitor.service.yaml
Making dashboard available at port 80


## Deployment: VPN  (optional)
There is a great [HELM package](https://github.com/kubernetes/charts/tree/master/stable/openvpn
)  that will install openVPN on your cluster. 

>The chart will automatically configure dns to use kube-dns and route all network traffic to kubernetes pods and services through the vpn. By connecting to this vpn a host is effectively inside a cluster's network.

Installation

```
helm repo add stable http://storage.googleapis.com/kubernetes-charts
helm install stable/openvpn
```


## Credits
Big thanks to the following github contributors

* [https://github.com/kuberstack/charts/tree/master/incubator/geth](https://github.com/kuberstack/charts/tree/master/incubator/geth)
* [https://github.com/jpoon/kubernetes-ethereum-chart](https://github.com/jpoon/kubernetes-ethereum-chart)
* [https://github.com/MaximilianMeister/kuberneteth](https://github.com/MaximilianMeister/kuberneteth)
* [https://github.com/wbuchwalter/ethereum-kubernetes-cluster](https://github.com/wbuchwalter/ethereum-kubernetes-cluster)
* [https://github.com/rehive/geth-helm-chart/tree/master/geth-chart](https://github.com/rehive/geth-helm-chart/tree/master/geth-chart)
* [https://github.com/skozin/ethereum-dev-net](https://github.com/skozin/ethereum-dev-net)