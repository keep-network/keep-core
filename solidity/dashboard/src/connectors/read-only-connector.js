import { AbstractConnector } from "./abstract-connector"
import { WALLETS } from "../constants/constants"
import { getWsUrl } from "./utils"
import Web3ProviderEngine from "web3-provider-engine"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket"

export class ReadOnlyConnector extends AbstractConnector {
  /** @type {string} To store selected account by user. */
  selectedAccount = ""

  constructor() {
    super(WALLETS.READ_ONLY_ADDRESS.name)
    this.websocketSubprovider = new WebsocketSubprovider({ rpcUrl: getWsUrl() })
  }

  setSelectedAccount = (address) => {
    this.selectedAccount = address
  }

  enable = async () => {
    if (!this.selectedAccount) {
      throw new Error("Account is not selected")
    }

    this.provider = new Web3ProviderEngine()
    this.provider.addProvider(
      new WebsocketSubprovider({
        rpcUrl: getWsUrl(),
      })
    )

    this.provider.start()

    return [this.selectedAccount]
  }

  disconnect = async () => {
    this.provider.stop()
  }
}
