class ContractWrapper {
  constructor(instance, deployedAtBlock) {
    this.address = instance.options.address
    this.deployedAtBlock = deployedAtBlock
    this.instance = instance
  }

  instance
  deployedAtBlock
  address

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

  get instance() {
    return this.instance
  }

  get deployedAtBlock() {
    return this.deployedAtBlock
  }

  get address() {
    return this.address
  }

  get methods() {
    this.instance.methods
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
      const deployTransactionHash = contract.networks[networkId].transactionHash
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
