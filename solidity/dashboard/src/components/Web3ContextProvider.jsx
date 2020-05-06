import React from "react"
import Web3 from "web3"
import { TrezorProvider } from "../connectors/trezor"
import { LedgerProvider, LEDGER_DERIVATION_PATHS } from "../connectors/ledger"
import { Web3Context } from "./WithWeb3Context"
import { MessagesContext, messageType } from "./Message"
import { getContracts } from "../contracts"

export default class Web3ContextProvider extends React.Component {
  static contextType = MessagesContext

  constructor(props) {
    super(props)
    this.state = {
      provider: null,
      web3: null,
      isFetching: false,
      yourAddress: "",
      networkType: "",
      token: { options: { address: "" } },
      stakingContract: { options: { address: "" } },
      grantContract: { options: { address: "" } },
      utils: {},
      eth: {},
      error: "",
    }
  }

  getWeb3 = (providerName) => {
    switch (providerName) {
      case "TREZOR": {
        return new Web3(new TrezorProvider())
      }
      case "METAMASK": {
        if (window.ethereum || window.web3) {
          return new Web3(window.ethereum || window.web3.currentProvider)
        }
        throw new Error("No browser extention")
      }
      case "COINBASE": {
        throw new Error("Coinbase wallet is not yet supported")
      }
      case "LEDGER_LIVE":
      case "LEDGER_LEGACY": {
        return new Web3(
          new LedgerProvider(LEDGER_DERIVATION_PATHS[providerName])
        )
      }
      default:
        throw new Error("Unsupported wallet")
    }
  }

  connectAppWithWallet = async (
    providerName,
    firstAccountAsSelected = false
  ) => {
    let web3
    let accounts
    this.setState({ isFetching: true })
    try {
      web3 = this.getWeb3(providerName)
      accounts = await web3.currentProvider.enable()
    } catch (error) {
      this.setState({ providerError: error.message, isFetching: false })
      this.context.showMessage({
        type: messageType.ERROR,
        title: error.message,
      })
      return
    }
    this.setState(
      {
        web3,
        provider: providerName,
        yourAddress: firstAccountAsSelected ? accounts[0] : null,
        accounts,
        networkType: await web3.eth.net.getNetworkType(),
      },
      this.setData
    )
  }

  setData = async () => {
    this.initializeContracts()
    this.state.web3.eth.currentProvider.on("accountsChanged", this.setAccount)
  }

  connectAppWithAccount = async () => {
    const { web3 } = this.state
    this.setState({ isFetching: true })
    try {
      const [yourAddress] = await web3.currentProvider.enable()
      this.setState({ yourAddress, isFetching: false })
    } catch (error) {
      this.setState({ providerError: error.message, isFetching: false })
      this.context.showMessage({
        type: messageType.ERROR,
        title: error.message,
      })
    }
  }

  initializeContracts = async () => {
    const { web3 } = this.state
    try {
      const contracts = await getContracts(web3)
      this.setState({
        ...contracts,
        utils: web3.utils,
        eth: web3.eth,
        isFetching: false,
      })
    } catch (error) {
      this.setState({
        isFetching: false,
        error: "Please select correct network",
      })
    }
  }

  setAccount = ([yourAddress]) => {
    if (!yourAddress) {
      this.setState({
        isFetching: false,
        yourAddress: "",
        token: { options: { address: "" } },
        stakingContract: { options: { address: "" } },
        grantContract: { options: { address: "" } },
      })
      return
    }
    this.setState({ yourAddress })
  }

  render() {
    return (
      <Web3Context.Provider
        value={{
          ...this.state,
          connectAppWithAccount: this.connectAppWithAccount,
          connectAppWithWallet: this.connectAppWithWallet,
          setAccount: this.setAccount,
        }}
      >
        {this.props.children}
      </Web3Context.Provider>
    )
  }
}
