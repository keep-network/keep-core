import React, { Component } from 'react';

class Emitter extends Component {

  state = { previousEntryKey: null }

  componentDidMount() {
    const { drizzle } = this.props
    const randomBeaconService = drizzle.contracts.KeepRandomBeaconServiceImplV1
    if (randomBeaconService) {
      const previousEntryKey = randomBeaconService.methods['previousEntry'].cacheCall()
      this.setState({ previousEntryKey });
    }
  }

  render() {
    const { previousEntryKey } = this.state
    const { drizzleState } = this.props

    const randomBeaconServiceState = drizzleState.contracts.KeepRandomBeaconServiceImplV1
    const previousEntry = randomBeaconServiceState && randomBeaconServiceState['previousEntry'][previousEntryKey]

    return (
      <div className="emitter">
        <div className="label">Latest Entry Emitted</div>
        <div className="entry">
          {(previousEntry && previousEntry.value) || 'none'}
        </div>
      </div>
    )
  }
}

export default Emitter