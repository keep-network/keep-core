import TrezorConnect from 'trezor-connect'
import Web3ProviderEngine from 'web3-provider-engine'
import { TrezorSubprovider } from '@0x/subproviders/lib/src/subproviders/trezor'
import { RPCSubprovider } from '@0x/subproviders/lib/src/subproviders/rpc_subprovider'

export class TrezorProvider extends Web3ProviderEngine {
  constructor(
    manifestEmail,
    manifestAppUrl,
  ) {
    super()
    TrezorConnect.manifest({
      email: manifestEmail,
      appUrl: manifestAppUrl,
    })
    this.addProvider(
      new TrezorSubprovider({
        trezorConnectClientApi: TrezorConnect,
        networkId: 1101,

      })
    )

    this.addProvider(new RPCSubprovider('localhost:8545'))

    this.start(() => {
      console.log('callback')
    })
  }
}
