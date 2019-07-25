import React, { Component } from 'react';
import Emitter from './Emitter'
import Requester from './Requester'

class App extends Component {

  state = { loading: true }

  componentDidMount() {
    const { drizzle } = this.props

    this.unsubscribe = drizzle.store.subscribe(() => {
      const drizzleState = drizzle.store.getState()

      if (drizzleState.drizzleStatus.initialized) {
        this.setState({ loading: false, drizzleState })
      }
    })
  }

  compomentWillUnmount() {
    this.unsubscribe();
  }

  render() {
    const { drizzle } = this.props
    const { drizzleState } = this.state

    if (this.state.loading) return "Loading Drizzle..."

    return (
      <div className="app">
        <header className="header">
          <div className="logo">
            {/* TODO: Keep "K" logo */}
            randomBeacon
          </div>
        </header>
        <div className="body">
          <div className="nav">
            <div className="nav-item selected">
              {/* TODO: "Emitter tower" logo */}
              Emitter
            </div>
          </div>
          <div className="main">
            <Emitter drizzleState={drizzleState} drizzle={drizzle} />
            <Requester drizzleState={drizzleState} drizzle={drizzle} />
          </div>
        </div>
      </div>
    )
  }
}

export default App
