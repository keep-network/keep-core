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

  createContractInstance = (
    abi,
    address,
    deploymentTxnHash,
    deploymentBlock
  ) => {
    return this._createContractInstance(
      abi,
      address,
      deploymentTxnHash,
      deploymentBlock
    )
  }

  get defaultAccount() {
    return this._defaultAccount
  }

  set defaultAccount(defaultAccount) {
    this._defaultAccount = defaultAccount
  }

  setProvider = (provider) => {
    this._setProvider(provider)
  }
}

class Web3jsWrapper extends Web3LibWrapper {
  _getTransaction = async (hash) => {
    return await this.lib.eth.getTransaction(hash)
  }

  _createContractInstance = (
    abi,
    address,
    deploymentTxnHash = null,
    deploymentBlock = null
  ) => {
    const contract = new this.lib.eth.Contract(abi, address)
    contract.options.from = this.defaultAccount
    contract.options.handleRevert = true

    return ContractFactory.createWeb3jsContract(
      contract,
      deploymentTxnHash,
      this,
      deploymentBlock
    )
  }

  set _defaultAccount(defaultAccount) {
    this.lib.eth.defaultAccount = defaultAccount
  }

  get _defaultAccount() {
    return this.lib.eth.defaultAccount
  }

  _setProvider = (provider) => {
    this.lib.setProvider(provider)
  }
}

export {
  Web3jsWrapper,
  Web3LibWrapper,
  ContractFactory,
  BaseContract,
  Web3jsContractWrapper,
}
