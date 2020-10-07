# KEEP Core Subgraph

## Run KEEP subgraph locally
### 1. Run Ganache

Note that host set to `0.0.0.0` is necessary for Ganache to be accessible from within Docker and from other machines. By default, Ganache only binds to 127.0.0.1, which can only be accessed from the host machine that Ganache runs on. If you use a ganache-cli run the node with `ganache-cli -h 0.0.0.0`. If you use a Ganche in GUI version go to `settings -> server ` and set `HOSTNAME` to `0.0.0.0 - All Interfaces`.

### 2. Run a local Graph Node
`docker-compose up`

### 3. Deploy contracts to Ganache

### 4. Install dependencies Run script that creates the subgraph manifest
Run `yarn && yarn create-manifest`

This script creates the subgraph manifest and adds contracts abi in the `abis` dir based on the `@keep-network/*` packages.

Note: for local development link the local version of  `@keep-network/*` pacakges.

### 5. Deploy the subgraph to the local Graph Node

#### 5.1 Run code generation
`yarn codegen`

#### 5.2 Allocate the subgraph name in the local Graph Node
`yarn create-local`

Note: use it only if your subgraph is not created in the local Graph node.

#### 5.3 Deploy the subgraph to your local Graph Node.
`yarn deploy:local`