import TrezorConnect from 'trezor-connect'
import { TrezorSubprovider } from '@0x/subproviders/lib/src/subproviders/trezor'
import { AbstractHardwareWalletConnector } from './abstract'

export class TrezorProvider extends AbstractHardwareWalletConnector {
  manifestEmail
  manifestAppUrl

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
      // debug: true,
    })
    this.setProvider(
      new TrezorSubprovider({
        trezorConnectClientApi: TrezorConnect,
        networkId: 1,
      }))
  }
}
