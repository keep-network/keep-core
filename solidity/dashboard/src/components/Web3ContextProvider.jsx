import React from "react"
import Web3 from "web3"
import { Web3Context } from "./WithWeb3Context"
import { MessagesContext } from "./Message"
import {
  getContracts,
  resolveWeb3Deferred,
  Web3Loaded,
  ContractsLoaded,
} from "../contracts"
import { getNetworkName } from "../utils/ethereum.utils"
import { isSameEthAddress } from "../utils/general.utils"
import { WALLETS } from "../constants/constants"

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

  connectAppWithWallet = async (connector, shouldSetError = true) => {
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

      this.explorerModeToWalletConnectionCheck(
        this.state.connector,
        connector,
        this.state.yourAddress,
        yourAddress
      )

      web3 = new Web3(connector.getProvider())
      web3.eth.defaultAccount = yourAddress

      await resolveWeb3Deferred(web3)
    } catch (error) {
      this.setState({ providerError: error.message, isFetching: false })
      throw error
    }

    try {
      contracts = await getContracts(web3, networkId)
    } catch (error) {
      this.setState({
        isFetching: false,
        error: shouldSetError ? error.message : null,
      })
      throw error
    }

    connector.on("accountsChanged", this.onAccountsChanged)
    connector.on("chainChanged", () => window.location.reload())
    connector.once("disconnect", this.disconnect)

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

  onAccountsChanged = async (yourAddress) => {
    if (!yourAddress) {
      await this.disconnect()
      return
    }

    const web3 = await Web3Loaded
    web3.eth.defaultAccount = yourAddress
    const contracts = await ContractsLoaded
    for (const contractInstance of Object.values(contracts)) {
      contractInstance.options.from = web3.eth.defaultAccount
    }

    this.setState({
      web3,
      yourAddress,
      ...contracts,
      utils: web3.utils,
      eth: web3.eth,
      isConnected: true,
    })
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

  explorerModeToWalletConnectionCheck = (
    explorerModeConnector,
    walletConnector,
    explorerModeAddress,
    walletAddress
  ) => {
    if (
      explorerModeConnector?.name === WALLETS.EXPLORER_MODE.name &&
      walletConnector?.name !== WALLETS.EXPLORER_MODE.name &&
      explorerModeAddress &&
      !isSameEthAddress(explorerModeAddress, walletAddress)
    ) {
      throw new Error(
        "Connected address is different from the one used in Explorer Mode."
      )
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

export default Web3ContextProvider
