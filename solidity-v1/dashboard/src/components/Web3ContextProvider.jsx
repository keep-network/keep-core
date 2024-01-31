import React from "react"
import Web3 from "web3"
import { Web3Context } from "./WithWeb3Context"
import { MessagesContext } from "./Message"
import {
  getContracts,
  resolveWeb3Deferred,
  Web3Loaded,
  ContractsLoaded,
  Keep,
  KeepExplorerMode,
} from "../contracts"
import { getNetworkName } from "../utils/ethereum.utils"
import { isSameEthAddress } from "../utils/general.utils"
import { WALLETS } from "../constants/constants"
import { getWsUrl } from "../connectors/utils"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"

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

      this.checkIfConnectionToWalletIsPossible(
        this.state.connector,
        this.state.yourAddress,
        connector,
        yourAddress
      )

      web3 = new Web3(connector.getProvider())
      web3.eth.defaultAccount = yourAddress

      await resolveWeb3Deferred(web3)
      Keep.setProvider(connector.getProvider())
      Keep.defaultAccount = yourAddress

      const explorerModeConnector = new ExplorerModeConnector()
      explorerModeConnector.setSelectedAccount(yourAddress)
      explorerModeConnector.enable()
      KeepExplorerMode.setProvider(explorerModeConnector.getProvider())
      KeepExplorerMode.defaultAccount = yourAddress
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
    Keep.defaultAccount = yourAddress
    KeepExplorerMode.defaultAccount = yourAddress
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
    // Set provider to the default one to fetch data w/o connected wallet.
    Keep.setProvider(new Web3.providers.WebsocketProvider(getWsUrl()))
    KeepExplorerMode.setProvider(
      new Web3.providers.WebsocketProvider(getWsUrl())
    )
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

  checkIfConnectionToWalletIsPossible = (
    prevConnector,
    prevAddress,
    nextConnector,
    nextAddress
  ) => {
    // Checks if an address on the wallet that the user is trying to connect to
    // is the same that was used in Explorer Mode
    if (
      prevConnector?.name === WALLETS.EXPLORER_MODE.name &&
      nextConnector?.name !== WALLETS.EXPLORER_MODE.name &&
      prevAddress &&
      !isSameEthAddress(prevAddress, nextAddress)
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
