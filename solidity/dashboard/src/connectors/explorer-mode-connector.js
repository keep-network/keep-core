import { AbstractConnector } from "./abstract-connector"
import { WALLETS } from "../constants/constants"
import { getWsUrl } from "./utils"
import Web3ProviderEngine from "web3-provider-engine"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket"
import ExplorerModeSubprovider from "./explorerModeSubprovider"
import { EventEmitter } from "events"

export class ExplorerModeConnector extends AbstractConnector {
  /** @type {string} To store selected account by user. */
  selectedAccount = ""

  constructor() {
    super(WALLETS.EXPLORER_MODE.name)
    this.websocketSubprovider = new WebsocketSubprovider({ rpcUrl: getWsUrl() })
    this.eventEmitter = new EventEmitter()
  }

  setSelectedAccount = (address) => {
    this.selectedAccount = address
  }

  enable = async () => {
    if (!this.selectedAccount) {
      throw new Error("Account is not selected")
    }

    this.provider = new Web3ProviderEngine()
    this.provider.addProvider(new ExplorerModeSubprovider(this.eventEmitter))
    this.provider.addProvider(this.websocketSubprovider)

    this.provider.start()

    return [this.selectedAccount]
  }

  disconnect = async () => {
    this.provider.stop()
  }
}
