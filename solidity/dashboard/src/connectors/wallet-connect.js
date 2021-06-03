import Web3ProviderEngine from "web3-provider-engine"
import WalletConnectSubprovider from "@walletconnect/web3-subprovider"
import CacheSubprovider from "web3-provider-engine/subproviders/cache"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket"
import { AbstractConnector } from "./abstract-connector"
import config from "../config/config.json"
import { WALLETS } from "../constants/constants"
import { overrideCacheMiddleware, getWsUrl, getChainId } from "./utils"
import { UserRejectedConnectionRequestError } from "."

export class WalletConnectConnector extends AbstractConnector {
  constructor(
    options = {
      rpc: {
        // Mainnet
        1: config.networks["1"].rpcURL,
        // Ropsten
        3: config.networks["3"].rpcURL,
        // Internal keep-dev network
        1101: config.networks["1101"].rpcURL,
        // Development- private network, change if you use a different one.
        1337: "http://localhost:8545",
      },
      chainId: getChainId(),
    }
  ) {
    super(WALLETS.WALLET_CONNECT.name)
    this._options = options
  }

  _onDisconnect = (error) => {
    if (error) {
      this.emit("error", error)
      return
    }
    this.emit("disconnect")
  }

  _onSessionUpdate = (error, payload) => {
    if (error) {
      this.emit("error", error)
      return
    }
    // Handle session update
    this._updateState.apply(this.walletConnectSubprovider, [payload.params[0]])
    this._handleAccountChanged(payload.params[0])
  }

  enable = async () => {
    try {
      this.provider = new Web3ProviderEngine()

      this.walletConnectSubprovider = new WalletConnectSubprovider(
        this._options
      )
      this.provider.addProvider(this.walletConnectSubprovider)

      const cacheSubprovider = new CacheSubprovider()
      this.provider.addProvider(cacheSubprovider)
      overrideCacheMiddleware(cacheSubprovider)

      this.provider.addProvider(
        new WebsocketSubprovider({ rpcUrl: getWsUrl() })
      )

      // There is a bug in the WallecConnect subprovider. The bug is that the
      // WalletConnect subprovider doesn't implement `updateState` and` emit`
      // functions so there is a need to override them to avoid throwing errors.
      this.walletConnectSubprovider.updateState = this._updateState
      this.walletConnectSubprovider.emit = (type, ...args) => {
        this.emit(type, ...args)
      }
      const wc = await this.walletConnectSubprovider.getWalletConnector()

      this.provider.start()
      wc.on("disconnect", this._onDisconnect)
      wc.on("session_update", this._onSessionUpdate)

      // Set the correct chainId and networkid for the WallectConnect
      // subprovider. Also override the `getChainId` and `getNetworkId` to avoid
      // additional requests to the node- we can return values from the
      // WalleConnect subprovider object.
      this.walletConnectSubprovider.chainId = await this.getChainId()
      this.walletConnectSubprovider.networkId = await this.getNetworkId()
      this.getChainId = async () => this.walletConnectSubprovider.chainId
      this.getNetworkId = async () => this.walletConnectSubprovider.networkId

      const accounts = this.walletConnectSubprovider.accounts
      // Save accounts to detect account change.
      this.accounts = accounts
      return accounts
    } catch (error) {
      if (error.message === "User closed modal") {
        throw new UserRejectedConnectionRequestError()
      }

      throw error
    }
  }

  disconnect = async () => {
    const wc = await this.walletConnectSubprovider.getWalletConnector({
      disableSessionCreation: true,
    })
    await wc.killSession()
    this.provider.stop()
    this.emit("disconnect")
  }

  async _handleAccountChanged(sessionParams) {
    const { accounts } = sessionParams
    if (!this.accounts || (accounts && this.accounts !== accounts)) {
      this.accounts = accounts
      this.emit("accountsChanged", accounts[0])
    }
  }

  async _updateState(sessionParams) {
    const { chainId, networkId } = sessionParams
    // Check if chainId changed and trigger event
    if (
      !this.chainId ||
      (chainId && Number(this.chainId) !== Number(chainId))
    ) {
      this.chainId = chainId
      this.emit("chainChanged", chainId)
    }
    // Check if networkId changed and trigger event
    if (!this.networkId || (networkId && this.networkId !== networkId)) {
      this.networkId = networkId
      this.emit("networkChanged", networkId)
    }
  }
}
