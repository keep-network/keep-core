import React, { useEffect } from "react"
import { Link } from "react-router-dom"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Routing from "./components/Routing"
import { Messages } from "./components/Message"
import { SideMenu } from "./components/SideMenu"
import { BrowserRouter as Router } from "react-router-dom"
import { Provider } from "react-redux"
import store from "./store"
import { ModalContextProvider } from "./components/Modal"
import * as Icons from "./components/Icons"
import Footer from "./components/Footer"
import { useWeb3Context } from "./components/WithWeb3Context"
import { LIQUIDITY_REWARD_PAIRS } from "./constants/constants"

const liquidityRewardNotificationFunc = (address) => {
  for (const [pairName, value] of Object.entries(LIQUIDITY_REWARD_PAIRS)) {
    store.dispatch({
      type: `liquidity_rewards/${pairName}_liquidity_rewards_earned_notification`,
      payload: { liquidityRewardPairName: pairName, address },
    })
  }
}

const App = () => (
  <Provider store={store}>
    <Messages>
      <Web3ContextProvider>
        <ModalContextProvider>
          <Router>
            <AppLayout />
          </Router>
        </ModalContextProvider>
      </Web3ContextProvider>
    </Messages>
  </Provider>
)

const AppLayout = () => {
  const { yourAddress, provider } = useWeb3Context()

  useEffect(() => {
    const isActive = yourAddress && provider
    if (isActive) {
      for (const [pairName, value] of Object.entries(LIQUIDITY_REWARD_PAIRS)) {
        store.dispatch({
          type: `liquidity_rewards/KEEP_ETH_notification_interval_active`,
          payload: { liquidityRewardPairName: pairName },
        })
      }
      const liquidityRewardsNotificationInterval = setInterval(() => {
        liquidityRewardNotificationFunc(yourAddress)
      }, 420000) // every 7 minutes
      return () => {
        clearInterval(liquidityRewardsNotificationInterval)
      }
    }
  }, [yourAddress, provider])

  return (
    <>
      <AppHeader />
      <section className="app__content">
        <Routing />
      </section>
      <Footer className="app__footer" />
    </>
  )
}

const AppHeader = () => {
  return (
    <header className="app__header">
      <Link to="/" className="app-logo">
        <Icons.KeepDashboardLogo />
      </Link>
      <SideMenu />
    </header>
  )
}
export default App
