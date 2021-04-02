import React from "react"
import Web3 from "web3"
import { Web3Context } from "./WithWeb3Context"
import { MessagesContext, useShowMessage } from "./Message"
import { getContracts, resolveWeb3Deferred } from "../contracts"
import { connect } from "react-redux"
import { WALLETS } from "../constants/constants"
import { getNetworkName } from "../utils/ethereum.utils"
import { isSameEthAddress } from "../utils/general.utils"

class Web3ContextProvider extends React.Component {
  static contextType = MessagesContext

  constructor(props) {
    super(props)
    this.state = {
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
      isConnected: false,
      connector: null,
    }
  }

  componentWillUnmount() {
    this.disconnect(false)
  }

  connectAppWithWallet = async (connector, payload = null) => {
    this.setState({ isFetching: true })
    let web3
    let yourAddress
    let contracts
    let networkId
    let chainId
    try {
      const accounts = await connector.enable()
      networkId = await connector.getNetworkId()
      chainId = await connector.getChainId()
      console.log(
        `Connected to the network; chainId: ${chainId.toString()}, networkId: ${networkId}`
      )
      yourAddress = accounts[0]

      if (
        this.state.connector?.name === WALLETS.EXPLORER_MODE.name &&
        connector.name !== WALLETS.EXPLORER_MODE.name
      ) {
        if (
          this.state.yourAddress &&
          !isSameEthAddress(this.state.yourAddress, yourAddress)
        ) {
          throw new Error(
            "Connected address is different from the one used in Explorer Mode."
          )
        }
      }

      web3 = new Web3(connector.getProvider())
      web3.eth.defaultAccount = yourAddress

      await resolveWeb3Deferred(web3)

      if (payload) {
        await web3.eth.currentProvider.sendAsync(payload)
      }
    } catch (error) {
      this.setState({ providerError: error.message, isFetching: false })
      throw error
    }

    try {
      contracts = await getContracts(web3, networkId)
    } catch (error) {
      this.setState({
        isFetching: false,
        error: error.message,
      })
      throw error
    }

    this.props.fetchKeepTokenBalance()

    this.setState({
      web3,
      yourAddress,
      networkType: getNetworkName(chainId),
      chainId,
      networkId,
      ...contracts,
      utils: web3.utils,
      eth: web3.eth,
      isFetching: false,
      connector,
      isConnected: true,
    })
  }

  abortWalletConnection = () => {
    this.setState({
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
      isConnected: false,
    })
  }

  connectAppWithAccount = async () => {
    const { connector } = this.state
    await this.connectAppWithWallet(connector)
  }

  refreshProvider = async ([yourAddress]) => {
    // if (!yourAddress) {
    //   this.setState({
    //     isFetching: false,
    //     yourAddress: "",
    //     token: { options: { address: "" } },
    //     stakingContract: { options: { address: "" } },
    //     grantContract: { options: { address: "" } },
    //   })
    //   return
    // }
    // const { connector, provider } = this.state
    // await this.connectAppWithWallet(connector, provider)

    // This is a temporary solution to prevent a situation when a user changed
    // an account but data has not been updated. After migrate to redux the dapp
    // fetches data only once and updates data based on emitted events. This
    // solution doesn't support a case where a user changed an account. We are
    // going to address it in a follow up work.
    window.location.reload()
  }

  disconnect = async (shouldSetState = true) => {
    const { connector } = this.state
    if (!connector) {
      return
    }

    await connector.disconnect()
    if (shouldSetState) {
      this.setState({
        web3: null,
        isFetching: false,
        yourAddress: "",
        networkType: "",
        utils: {},
        eth: {},
        error: "",
        isConnected: false,
        connector: null,
      })
    }
  }

  render() {
    return (
      <Web3Context.Provider
        value={{
          ...this.state,
          connectAppWithAccount: this.connectAppWithAccount,
          connectAppWithWallet: this.connectAppWithWallet,
          abortWalletConnection: this.abortWalletConnection,
          disconnect: this.disconnect,
        }}
      >
        {this.props.children}
      </Web3Context.Provider>
    )
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    fetchKeepTokenBalance: () =>
      dispatch({ type: "keep-token/balance_request" }),
  }
}

export default connect(null, mapDispatchToProps)(Web3ContextProvider)
