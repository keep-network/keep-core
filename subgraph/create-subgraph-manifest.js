const fs = require('fs')
const { ethers } = require('ethers')
const Mustache = require('mustache')
const TokenStakingJSON = require('@keep-network/keep-core/artifacts/TokenStaking.json')
const TokenGrantJSON = require('@keep-network/keep-core/artifacts/TokenGrant.json') 
const ManagedGrantFactoryJSON = require('@keep-network/keep-core/artifacts/ManagedGrantFactory.json') 

const contracts = [
    { artifact: TokenStakingJSON, name: "TokenStaking" },
    { artifact: TokenGrantJSON, name: "TokenGrant" },
    { artifact: ManagedGrantFactoryJSON, name: "ManagedGrant" }

]

// The artifacts from `@keep-network/*` for a given build only support a single network id
function getFirstNetworkIdFromArtifact() {
    return Object.keys(TokenStakingJSON.networks)[0]
}

const provider = new ethers.providers.JsonRpcProvider()

const updateSubgraphManifest = async () => {
    const networkId = getFirstNetworkIdFromArtifact()
    let network

    if (networkId === "3"){
        network = "ropsten"
    } else if (networkId === "1") {
        network= "mainnet"
    } else {
        network = "local"
    }

    const mustacheView = {}
    for(const contract of contracts) {
        const address = contract.artifact.networks[networkId].address
        const deployTransactionHash = contract.artifact.networks[networkId].transactionHash
        const startBlock = (await provider.getTransaction(deployTransactionHash)).blockNumber
    
        mustacheView[contract.name] = { address, startBlock }
        fs.writeFileSync(`abis/${contract.name}.json`, JSON.stringify(contract.artifact.abi))
    }

    mustacheView.network = network
    const subgraphTemplate = fs.readFileSync("subgraph.template.yaml", "utf8")

    const rendered = Mustache.render(subgraphTemplate, mustacheView)
    fs.writeFileSync(`subgraph.yaml`, rendered)
}


updateSubgraphManifest()
