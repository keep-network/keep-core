import TrezorConnect from 'trezor-connect'
import Web3ProviderEngine from 'web3-provider-engine'
import WebsocketSubprovider from 'web3-provider-engine/subproviders/websocket'
import CacheSubprovider from 'web3-provider-engine/subproviders/cache'
import { TrezorSubprovider } from '@0x/subproviders/lib/src/subproviders/trezor'
import { RPCSubprovider } from '@0x/subproviders/lib/src/subproviders/rpc_subprovider'

export class TrezorProvider extends Web3ProviderEngine {
  constructor(
    manifestEmail,
    manifestAppUrl,
  ) {
    super()
    TrezorConnect.init({
      lazyLoad: true, // this param will prevent iframe injection until TrezorConnect.method will be called
      manifest: {
        email: manifestEmail,
        appUrl: manifestAppUrl,
      },
      popup: true,
      debug: true,
    })
    this.addProvider(
      new TrezorSubprovider({
        trezorConnectClientApi: TrezorConnect,
        networkId: 1,
      })
    )
    this.addProvider(new WebsocketSubprovider({ rpcUrl: 'ws://localhost:8545', debug: true }))
    this.addProvider(new RPCSubprovider('http://localhost:8545'))
    this.addProvider(new CacheSubprovider())

    this.start()
  }
}
