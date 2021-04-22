import { AbstractConnector } from "./abstract-connector"
import { WALLETS } from "../constants/constants"
import { getWsUrl } from "./utils"
import Web3ProviderEngine from "web3-provider-engine"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket"
import ExplorerModeSubprovider from "./explorerModeSubprovider"

export class ExplorerModeConnector extends AbstractConnector {
  /** @type {string} To store selected account by user. */
  selectedAccount = ""

  constructor() {
    super(WALLETS.EXPLORER_MODE.name)
    this.explorerModeSubprovider = new ExplorerModeSubprovider()
  }

  setSelectedAccount = (address) => {
    this.selectedAccount = address
    this.emit("accountsChanged", address)
  }

  enable = async () => {
    if (!this.selectedAccount) {
      throw new Error("Account is not selected")
    }

    this.provider = new Web3ProviderEngine()
    this.provider.addProvider(this.explorerModeSubprovider)
    this.provider.addProvider(new WebsocketSubprovider({ rpcUrl: getWsUrl() }))

    this.explorerModeSubprovider.on(
      "chooseWalletAndSendTransaction",
      (payload) => {
        this.emit("chooseWalletAndSendTransaction", payload)
      }
    )

    this.provider.start()

    return [this.selectedAccount]
  }

  disconnect = async () => {
    this.explorerModeSubprovider.removeListener(
      "chooseWalletAndSendTransaction",
      () => {}
    )
    this.provider.stop()
    this.emit("disconnect")
  }
}
