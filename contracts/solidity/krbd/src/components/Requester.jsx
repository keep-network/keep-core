import React, { Component } from 'react';

class Requester extends Component {

  state = { requestEntryStackId: null }

  handleRequestEntry = () => {
    const { drizzle } = this.props

    // TODO: Generate random seed
    const seed = drizzle.web3.utils.toBN('27182818284590452353602874713526624977572470936999595749669676277240766303535')
    const account = drizzle.web3.eth.accounts.givenProvider.selectedAddress

    const randomBeaconService = drizzle.contracts.KeepRandomBeaconServiceImplV1
    const requestRelayEntry = randomBeaconService.methods['requestRelayEntry']

    // TODO: Error handling...wow drizzle's error handling is bad...
    var requestEntryStackId = requestRelayEntry.cacheSend(seed, { value: 10, from: account })

    this.setState({ requestEntryStackId })
  }

  getTransactionStatus = () => {
    const { requestEntryStackId } = this.state
    const { transactions, transactionStack } = this.props.drizzleState

    const transactionHash = transactionStack[requestEntryStackId]

    if (!transactionHash || !transactions[transactionHash]) {
      return 'unknown'
    }

    return transactions[transactionHash].status
  }

  render() {
    const status = this.getTransactionStatus()

    return (
      <div className="requester">
        <button onClick={this.handleRequestEntry} disabled={status === 'pending'}>
          {status === 'pending' ? 'Requesting...' : 'Request New Relay Entry'}
        </button>
      </div>
    )
  }
}

export default Requester