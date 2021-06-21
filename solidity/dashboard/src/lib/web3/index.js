import {
  ContractFactory,
  BaseContract,
  Web3jsContractWrapper,
} from "./contract"

class Web3LibWrapper {
  constructor(_lib) {
    this.lib = _lib
  }

  getTransaction = async (hash) => {
    return await this._getTransaction(hash)
  }

  createContractInstance = (artifact) => {
    return this._createContractInstance(artifact)
  }

  get defaultAccount() {
    return this._defaultAccount
  }

  set defaultAccount(defaultAccount) {
    this._defaultAccount = defaultAccount
  }
}

class Web3jsWrapper extends Web3LibWrapper {
  _getTransaction = async (hash) => {
    return await this.lib.eth.getTransaction(hash)
  }

  _createContractInstance = (artifact) => {
    const { networks, abi } = artifact
    const networkId = Object.keys(networks)[0]

    const address = networks[networkId].address
    const deploymentTxnHash = networks[networkId].transactionHash

    const contract = new this.lib.eth.Contract(abi, address)
    contract.options.defaultAccount = this.defaultAccount
    contract.options.handleRevert = true

    return ContractFactory.createWeb3jsContract(
      contract,
      deploymentTxnHash,
      this
    )
  }

  set _defaultAccount(defaultAccount) {
    this.lib.eth.defaultAccount = defaultAccount
  }

  get _defaultAccount() {
    return this.lib.eth.defaultAccount
  }
}

export {
  Web3jsWrapper,
  Web3LibWrapper,
  ContractFactory,
  BaseContract,
  Web3jsContractWrapper,
}
