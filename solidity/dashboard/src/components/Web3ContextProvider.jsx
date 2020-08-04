import React from "react"
import Web3 from "web3"
import { Web3Context } from "./WithWeb3Context"
import { MessagesContext, messageType } from "./Message"
import { getContracts, Web3Deferred } from "../contracts"

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

  connectAppWithWallet = async (connector, providerName) => {
    this.setState({ isFetching: true })
    let web3
    let yourAddress
    let contracts
    try {
      const accounts = await connector.enable()
      yourAddress = accounts[0]

      web3 = new Web3(connector)
      web3.eth.defaultAccount = yourAddress

      Web3Deferred.resolve(web3)
    } catch (error) {
      this.setState({ providerError: error.message, isFetching: false })
      this.context.showMessage({
        type: messageType.ERROR,
        title: error.message,
      })
      return
    }

    try {
      contracts = await getContracts(web3)
    } catch (error) {
      this.setState({
        isFetching: false,
        error: "Please select correct network",
      })
      return
    }

    web3.eth.currentProvider.on("accountsChanged", this.onAccountChanged)

    this.setState({
      web3,
      provider: providerName,
      yourAddress,
      networkType: await web3.eth.net.getNetworkType(),
      ...contracts,
      utils: web3.utils,
      eth: web3.eth,
      isFetching: false,
      connector,
    })
  }

  abortWalletConnection = () => {
    this.setState({
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
    })
  }

  connectAppWithAccount = async () => {
    const { connector, provider } = this.state
    await this.connectAppWithWallet(connector, provider)
  }

  onAccountChanged = async ([yourAddress]) => {
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

    const { connector, provider } = this.state
    await this.connectAppWithWallet(connector, provider)
  }

  render() {
    return (
      <Web3Context.Provider
        value={{
          ...this.state,
          connectAppWithAccount: this.connectAppWithAccount,
          connectAppWithWallet: this.connectAppWithWallet,
          abortWalletConnection: this.abortWalletConnection,
        }}
      >
        {this.props.children}
      </Web3Context.Provider>
    )
  }
}
