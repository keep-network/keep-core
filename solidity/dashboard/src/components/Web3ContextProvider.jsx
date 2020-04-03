import React from 'react'
import Web3 from 'web3'
import { TrezorProvider } from '../connectors/trezor'
import { Web3Context } from './WithWeb3Context'
import { MessagesContext, messageType } from './Message'
import { getContracts } from '../contracts'

export default class Web3ContextProvider extends React.Component {
    static contextType = MessagesContext

    constructor(props) {
      super(props)
      this.state = {
        showModal: true,
        provider: null,
        web3: null,
        isFetching: false,
        yourAddress: '',
        networkType: '',
        token: { options: { address: '' } },
        stakingContract: { options: { address: '' } },
        grantContract: { options: { address: '' } },
        utils: {},
        eth: {},
        error: '',
      }
    }

    showModal = () => this.setState({ showModal: true })
    closeModal = () => this.setState({ showModal: false })

    getWeb3 = (providerName) => {
      switch (providerName) {
        case 'TREZOR': {
          return new Web3(new TrezorProvider('test@email.com', 'https://keep.network/'))
        }
        case 'METAMASK': {
          if (window.ethereum || window.web3) {
            return new Web3(window.ethereum || window.web3.currentProvider)
          }
          throw new Error('No browser extention')
        }
        case 'COINBASE': {
          throw new Error('Coinbase wallet is not yet supported')
        }
        case 'LEDGER': {
          throw new Error('Ledger wallet is not yet supported')
        }
        default:
          throw new Error('Unsupported wallet')
      }
    }

    connectAppWithWallet = (providerName) => {
      let web3
      try {
        web3 = this.getWeb3(providerName)
      } catch (error) {
        this.setState({ providerError: error.message })
        this.context.showMessage({ type: messageType.ERROR, title: error.message })
        return
      }
      this.setState({ web3, provider: providerName, showModal: false }, this.setData)
    }

    initialize = async () => {
      const web3 = this.getWeb3()
      if (!web3) {
        return
      }
      this.setState({ web3 }, this.setData)
    }

    setData = async () => {
      const { web3 } = this.state
      this.setState({ isFetching: true })
      const accounts = await web3.eth.getAccounts()
      console.log('accounts', accounts)
      this.connectAppWithAccount(!accounts || accounts.length === 0, accounts)
      this.initializeContracts()
      this.state.web3.eth.currentProvider.on('accountsChanged', this.accountHasBeenChanged)
    }

    connectAppWithAccount = async (withInfoMessage = true, accounts) => {
      const { web3 } = this.state
      this.setState({ isFetching: true })
      withInfoMessage && this.context.showMessage({ type: messageType.INFO, title: 'Please check web3 provider' })

      try {
        const [account] = accounts
        console.log('account', account)
        this.setState({
          yourAddress: account,
          networkType: await web3.eth.net.getNetworkType(),
          isFetching: false,
        })
      } catch (error) {
        console.log('error', error)
        this.context.showMessage({ type: 'error', title: error.message })
        this.setState({ isFetching: false })
      }
    }

    initializeContracts = async () => {
      const { web3 } = this.state
      try {
        const contracts = await getContracts(web3)
        this.setState({
          ...contracts,
          defaultContract: contracts.stakingContract,
          utils: web3.utils,
          eth: web3.eth,
        })
      } catch (error) {
        console.log('error contracts', error)
        this.setState({
          error: 'Please select correct network',
        })
      }
    }

    accountHasBeenChanged = ([yourAddress]) => {
      if (!yourAddress) {
        this.setState({
          isFetching: false,
          yourAddress: '',
          token: { options: { address: '' } },
          stakingContract: { options: { address: '' } },
          grantContract: { options: { address: '' } },
        })
        return
      }
      this.setState({ yourAddress })
    }

    render() {
      const { showModal, ...restState } = this.state

      return (
        <Web3Context.Provider
          value={{
            ...restState,
            connectAppWithAccount: this.connectAppWithAccount,
            connectAppWithWallet: this.connectAppWithWallet,
            showConnectWalletModal: this.showModal,
          }}
        >
          {this.props.children}
        </Web3Context.Provider>

      )
    }
}
