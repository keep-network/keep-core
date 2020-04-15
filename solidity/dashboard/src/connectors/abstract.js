import Web3ProviderEngine from 'web3-provider-engine'
import WebsocketSubprovider from 'web3-provider-engine/subproviders/websocket'
import CacheSubprovider from 'web3-provider-engine/subproviders/cache'
import { getRpcURL } from './utils'

export class AbstractHardwareWalletConnector extends Web3ProviderEngine {
    provider

    constructor(provider) {
      super()
      this.provider = provider
      this.addProvider(this.provider)
      this.addProvider(new CacheSubprovider())
      this.addProvider(new WebsocketSubprovider({ rpcUrl: getRpcURL(), debug: false }))
    }

    enable = async () => {
      this.start()

      return await this.getAccounts()
    }

    getAccounts = async () => {
      return await this.provider.getAccountsAsync()
    }

    setProvider = (provider) => {
      this.provider = provider
    }

    getProvider = () => {
      return this.provider
    }
}
