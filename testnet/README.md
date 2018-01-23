
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


## Preparing deployment files
Rather than creating YAML files with hardcoded values we're gonna use [Helm](https://helm.sh/) package manager that can generate them based on the provided variables in **values.yaml **


### Bootnode configuration

| Parameter                  | Description                        | Default                                                    |
| -----------------------    | ---------------------------------- | ---------------------------------------------------------- |
| `replicaCount`    | Number of replicas  | 1
| `image.repository` | `geth` image   | ethereum/client-go
| `image.tag` | `geth` image tag | alltools-v1.7.3 (Geth + Tools)                                       

These are used in `bootnode.deployment.yaml` and `bootnode.service.yaml`. A quick note on bootnode setup:
>Each ethereum node, including a bootnode is identified by an enode identifier. These identifiers are derived from a key. Therefore you will need to give the bootnode such key.

We will instruct bootnode to generate a key and use Kubernetes `initContainers` feature to do this in advance.


### Miner configuration and genesis block configuration

| Parameter                  | Description                        | Default                                                    |
| -----------------------    | ---------------------------------- | ---------------------------------------------------------- |
| `replicaCount`    | Number of replicas  | 1
| `image.repository` | `geth` image   | ethereum/client-go
| `image.tag` | `geth` image tag | v1.7.3
| `networkId`    | Network ID  | 1101
| `initialAccounts` | Initial accounts |   

These are used in `miner.configmap.yaml`, `miner.deployment.yaml` and `miner.service.yaml`. 
The following actions are performed by `initContainers` step before the miner is started:

* init our network with genesis.json
* Get bootnode enode identifier
 


