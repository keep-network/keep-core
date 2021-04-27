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

  chooseWalletAndSendTransactionHandler = (payload) => {
    this.emit("chooseWalletAndSendTransaction", payload)
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
      this.chooseWalletAndSendTransactionHandler
    )

    this.provider.start()

    return [this.selectedAccount]
  }

  disconnect = async () => {
    this.explorerModeSubprovider.removeListener(
      "chooseWalletAndSendTransaction",
      this.chooseWalletAndSendTransactionHandler
    )
    this.provider.stop()
    this.emit("disconnect")
  }
}
