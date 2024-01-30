import React from "react"
import { Provider } from "react-redux"
import { BrowserRouter as Router } from "react-router-dom"
import store from "./store"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Routing from "./components/Routing"
import { Messages } from "./components/Message"
import { SideMenu } from "./components/SideMenu"
import * as Icons from "./components/Icons"
import Footer from "./components/Footer"
import useSubscribeToConnectorEvents from "./hooks/useSubscribeToConnectorEvents"
import useAutoConnect from "./hooks/useAutoConnect"
import useAutoWalletAddressInjectIntoUrl from "./hooks/useAutoWalletAddressInjectIntoUrl"
import useModalWindowForMobileUsers from "./hooks/useModalWindowForMobileUsers"
import { ModalRoot } from "./components/modal"
import NavLink from "./components/NavLink"
import { useShowLegacyDappModal } from "./hooks/useShowLegacyDappModal"
import {
  useSubscribeToCopyStakeEvents,
  useSubscribeToCovPoolsAuctionClosedEvent,
  useSubscribeToCovPoolsAuctionCreatedEvent,
  useSubscribeToCovPoolsWithdrawalCompletedEvent,
  useSubscribeToCovPoolsWithdrawalInitiatedEvent,
  useSubscribeToDepositWithdrawEvent,
  useSubscribeToDepositedEvent,
  useSubscribeToECDSARewardsClaimedEvent,
  useSubscribeToERC20TransferEvent,
  useSubscribeToOperatorUndelegatedEvent,
  useSubscribeToRecoveredStakeEvent,
  useSubscribeToStakedEvents,
  useSubscribeToTokenGrantWithdrawnEvent,
  useSubscribeToTopUpCompletedEvent,
  useSubscribeToTopUpInitiatedEvent,
  useSubscribeToUndelegatedEvents,
} from "./hooks/subscribtions"
import { useSubscribeToThresholdStakeKeepEvent } from "./hooks/subscribtions/useSubscribeToThresholdStakeKeepEvent"

const App = () => (
  <Provider store={store}>
    <Router basename={`${process.env.PUBLIC_URL}`}>
      <Messages>
        <Web3ContextProvider>
          <ModalRoot />
          <AppLayout />
        </Web3ContextProvider>
      </Messages>
    </Router>
  </Provider>
)

const AppLayout = () => {
  useAutoConnect()
  useAutoWalletAddressInjectIntoUrl()
  useSubscribeToConnectorEvents()
  useModalWindowForMobileUsers()
  useShowLegacyDappModal()

  useSubscribeToERC20TransferEvent()
  useSubscribeToStakedEvents()
  useSubscribeToUndelegatedEvents()
  useSubscribeToRecoveredStakeEvent()
  useSubscribeToTokenGrantWithdrawnEvent()
  useSubscribeToDepositWithdrawEvent()
  useSubscribeToDepositedEvent()
  useSubscribeToTopUpInitiatedEvent()
  useSubscribeToTopUpCompletedEvent()
  useSubscribeToECDSARewardsClaimedEvent()
  useSubscribeToOperatorUndelegatedEvent()
  useSubscribeToCovPoolsWithdrawalInitiatedEvent()
  useSubscribeToCovPoolsWithdrawalCompletedEvent()
  useSubscribeToCovPoolsAuctionCreatedEvent()
  useSubscribeToCovPoolsAuctionClosedEvent()
  useSubscribeToThresholdStakeKeepEvent()
  useSubscribeToCopyStakeEvents()

  return (
    <>
      <AppHeader />
      <section className="app__content">
        <div className="bg-yellow-400 text-center text-yellow-900">
          This is a legacy dashboard. Only a small part of the features are
          still supported.
        </div>
        <Routing />
      </section>
    </>
  )
}

const styles = {
  legacyBadge: {
    wrapper: {
      marginTop: "1.5rem",
      padding: "0 2.75rem",
      borderRadius: "4px",
    },
    text: {
      fontWeight: 600,
    },
  },
}

const AppHeader = () => {
  return (
    <header className="app__header">
      <div
        className="bg-yellow-300 text-yellow-900 self-center"
        style={styles.legacyBadge.wrapper}
      >
        <h6 style={styles.legacyBadge.text}>LEGACY</h6>
      </div>
      <NavLink to="/overview" className="app-logo">
        <Icons.KeepDashboardLogo />
      </NavLink>
      <SideMenu />
      <Footer className="app__footer" />
    </header>
  )
}
export default App
