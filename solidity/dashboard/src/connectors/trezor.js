import TrezorConnect from 'trezor-connect'
import { TrezorSubprovider } from '@0x/subproviders/lib/src/subproviders/trezor'
import { AbstractHardwareWalletConnector } from './abstract'

TrezorConnect.init({
  lazyLoad: true, // this param will prevent iframe injection until TrezorConnect.method will be called
  manifest: {
    email: 'work@keep.network',
    appUrl: 'keep.network',
  },
  popup: true,
  //  debug: true,
})

export class TrezorProvider extends AbstractHardwareWalletConnector {
  constructor() {
    super(new TrezorSubprovider({
      trezorConnectClientApi: TrezorConnect,
      networkId: 1,
    }))
  }
}
