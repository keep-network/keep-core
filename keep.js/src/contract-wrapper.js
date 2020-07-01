class ContractWrapper {
  constructor(instance, deployedAtBlock) {
    this.deployedAtBlock = deployedAtBlock
    this.instance = instance
  }

  async makeCall(contractMethodName, ...args) {
    return await this.instance.methods[contractMethodName](...args).call()
  }

  sendTransaction(contractMethodName, ...args) {
    return this.instance.methods[contractMethodName](...args).send
  }

  async getPastEvents(eventName, filter, fromBlock = this.deployedAtBlock) {
    const searchFilter = { fromBlock, filter }
    return await this.instance.getPastEvents(eventName, searchFilter)
  }

  get address() {
    return this.instance.options.address
  }

  get methods() {
    return this.instance.methods
  }
}

class ContractFactory {
  static async createContractInstance(artifact, config) {
    const { web3, networkId } = config

    const lookupArtifactAddress = () => {
      const networkInfo = artifact.networks[networkId]
      if (!networkInfo) {
        throw new Error(
          `No contract ${artifact.contractName} for a given network ID ${networkId}.`
        )
      }
      return networkInfo.address
    }

    const contractDeployedAtBlock = async () => {
      const deployTransactionHash = artifact.networks[networkId].transactionHash
      const transaction = await web3.eth.getTransaction(deployTransactionHash)

      return transaction.blockNumber.toString()
    }

    const address = lookupArtifactAddress(artifact)
    const deployedAtBlock = await contractDeployedAtBlock(web3, artifact)
    const instance = new web3.eth.Contract(artifact.abi, address)
    instance.options.from = web3.eth.defaultAccount

    return new ContractWrapper(instance, deployedAtBlock)
  }
}

export default ContractFactory
