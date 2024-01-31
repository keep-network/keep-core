import EthereumProvider from "@walletconnect/ethereum-provider"
import { getChainId } from "./utils"
import config from "../config/config.json"
import { AbstractConnector } from "./abstract-connector"
import { WALLETS } from "../constants/constants"

const chainId = getChainId()
const rpcUrl = config.networks[chainId.toString()].rpcURL
const walletConnectProjectId = process.env.REACT_APP_WALLET_CONNECT_PROJECT_ID

export class UserRejectedRequestError extends Error {
  constructor() {
    super()
    this.name = this.constructor.name
    this.message = "The user rejected the request."
  }
}

function getSupportedChains({ chains, rpc }) {
  if (chains) {
    return chains
  }

  return rpc ? Object.keys(rpc).map((k) => Number(k)) : []
}

/**
 * Connector for WalletConnect V2
 */
export class WalletConnectV2Connector extends AbstractConnector {
  provider
  config
  rpcMap

  constructor() {
    super(WALLETS.WALLET_CONNECT.name)
    const options = {
      rpc: {
        [Number(chainId)]: rpcUrl,
      },
      chainId: chainId,
    }
    const config = {
      chains: [Number(chainId)],
      rpc: {
        [Number(chainId)]: rpcUrl,
      },
      projectId: walletConnectProjectId,
      showQrModal: true,
    }

    this.config = config
    this.rpcMap = config.rpc
    this._options = options

    this.enable = this.enable.bind(this)
    this.handleOnConnect = this.handleOnConnect.bind(this)
    this.handleOnDisplayUri = this.handleOnDisplayUri.bind(this)

    this.handleChainChanged = this.handleChainChanged.bind(this)
    this.handleAccountsChanged = this.handleAccountsChanged.bind(this)
    this.handleDisconnect = this.handleDisconnect.bind(this)
  }

  handleOnConnect = () => {}

  handleOnDisplayUri = () => {}

  /**
   * Handles chain change
   *
   * @param {string} newChainId
   */
  handleChainChanged = (newChainId) => {
    this.emitUpdate({ chainId: newChainId })
    if (newChainId !== `0x${chainId}`) this.deactivate()
  }

  /**
   * Handles account change
   *
   * @param {string[]} accounts
   */
  handleAccountsChanged = (accounts) => {
    this.emitUpdate({ account: accounts[0] })
  }

  handleDisconnect = () => {
    // We have to do this because of a @walletconnect/web3-provider bug
    if (this.provider) {
      this.provider.removeListener("chainChanged", this.handleChainChanged)
      this.provider.removeListener(
        "accountsChanged",
        this.handleAccountsChanged
      )
      this.provider.removeListener("display_uri", this.handleOnDisplayUri)
      this.provider.removeListener("connect", this.handleOnConnect)
      this.provider = undefined
    }
    this.emit("disconnect")
  }

  enable = async () => {
    // Removes all local storage entries that starts with "wc@2"
    // This is a workaround for "Cannot convert undefined or null to object"
    // error that sometime occur with WalletConnect
    Object.keys(localStorage)
      .filter((x) => x.startsWith("wc@2"))
      .forEach((x) => localStorage.removeItem(x))
    if (!this.provider) {
      const chains = getSupportedChains(this.config)
      if (chains.length === 0) throw new Error("Chains not specified!")
      this.provider = await EthereumProvider.init({
        projectId: this.config.projectId,
        chains: chains,
        rpcMap: this.rpcMap,
        showQrModal: this.config.showQrModal,
      })
    }
    if (chainId !== this.provider.chainId) {
      this.deactivate()
    }

    this.provider.on("connect", this.handleOnConnect)
    this.provider.on("display_uri", this.handleOnDisplayUri)

    const account = await new Promise((resolve, reject) => {
      const userReject = () => {
        // Erase the provider manually
        this.provider = undefined
        reject(new UserRejectedRequestError())
      }

      // Workaround to bubble up the error when user reject the connection
      this.provider.on("disconnect", () => {
        // Check provider has not been enabled to prevent this event callback
        // from being called in the future
        if (!account) {
          userReject()
        }
      })

      this.provider
        .enable()
        .then((accounts) => resolve(accounts[0]))
        .catch((error) => {
          // TODO: ideally this would be a better check
          if (error.message === "User closed modal") {
            userReject()
            return
          }
          reject(error)
        })
    }).catch((err) => {
      this.emitError(err)
      throw err
    })

    this.provider.on("disconnect", this.handleDisconnect)
    this.provider.on("chainChanged", this.handleChainChanged)
    this.provider.on("accountsChanged", this.handleAccountsChanged)
    return [account]
  }

  /**
   * @return {any} provider
   */
  getProvider = () => {
    return this.provider
  }

  /**
   * @param {any} provider
   */
  setProvider = (provider) => {
    this.provider = provider
  }

  /**
   * @return {Promise<number>} chainId
   */
  getChainId = () => {
    return Promise.resolve(this.provider.chainId)
  }

  /**
   * @return {Promise<number>} networkId
   */
  getNetworkId = () => {
    return Promise.resolve(this.provider.chainId)
  }

  /**
   * @return {string[]} accounts
   */
  getAccounts = () => {
    return Promise.resolve(this.provider.accounts)
  }

  disconnect = () => {
    if (this.provider) {
      this.provider.removeListener("disconnect", this.handleDisconnect)
      this.provider.removeListener("chainChanged", this.handleChainChanged)
      this.provider.removeListener(
        "accountsChanged",
        this.handleAccountsChanged
      )
      this.provider.removeListener("display_uri", this.handleOnDisplayUri)
      this.provider.removeListener("connect", this.handleOnConnect)
      this.provider.disconnect()
      this.provider = undefined
      this.emit("disconnect")
    }
  }

  close = () => {
    this.emit("disconnect")
  }
}
