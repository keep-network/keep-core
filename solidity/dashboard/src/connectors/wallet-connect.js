import Web3ProviderEngine from "web3-provider-engine"
import WalletConnectSubprovider from "@walletconnect/web3-subprovider"
import CacheSubprovider from "web3-provider-engine/subproviders/cache"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket"
import { AbstractConnector } from "./abstract-connector"
import { WALLETS } from "../constants/constants"
import { overrideCacheMiddleware, getWsUrl } from "./utils"

export class WalletConnectConnector extends AbstractConnector {
  constructor(
    options = {
      // TODO add chains
      rpc: {
        1337: "http://localhost:8545",
      },
    }
  ) {
    super(WALLETS.WALLET_CONNECT.name)
    this.walletConnectSubprovider = new WalletConnectSubprovider(options)
  }

  enable = async () => {
    try {
      this.provider = new Web3ProviderEngine()

      this.provider.addProvider(this.walletConnectSubprovider)

      const cacheSubprovider = new CacheSubprovider()
      this.provider.addProvider(cacheSubprovider)
      overrideCacheMiddleware(cacheSubprovider)

      this.provider.addProvider(
        new WebsocketSubprovider({ rpcUrl: getWsUrl() })
      )

      this.provider.start()

      // Set the correct chainId and networkid for the WallectConnect
      // subprovider. Also override the `getChainId` and `getNetworkId` to avoid
      // additional requests to the node- we can return values from the
      // WalleConnect subprovider object.
      this.walletConnectSubprovider.chainId = await this.getChainId()
      this.walletConnectSubprovider.networkId = await this.getNetworkId()
      this.getChainId = async () => this.walletConnectSubprovider.chainId
      this.getNetworkId = async () => this.walletConnectSubprovider.networkId

      // There is a bug in the WallecConnect subprovider. The bug is that the
      // WalletConnect subprovider doesn't implement `updateState` and` emit`
      // functions so there is a need to override them to avoid throwing errors.
      this.walletConnectSubprovider.updateState = () => {}
      this.walletConnectSubprovider.emit = () => {}

      await this.walletConnectSubprovider.getWalletConnector()

      const accounts = this.walletConnectSubprovider.accounts
      return accounts
    } catch (error) {
      if (error.message === "User closed modal") {
        throw new Error("The user rejected the request")
      }

      throw error
    }
  }

  disconnect = async () => {
    await this.provider.disconnect()
  }
}
