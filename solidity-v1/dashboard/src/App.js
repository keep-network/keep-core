import React from "react"
import { Provider } from "react-redux"
import { BrowserRouter as Router, Link } from "react-router-dom"
import store from "./store"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Routing from "./components/Routing"
import { Messages } from "./components/Message"
import { SideMenu } from "./components/SideMenu"
import { ModalContextProvider } from "./components/Modal"
import * as Icons from "./components/Icons"
import Footer from "./components/Footer"
import useSubscribeToConnectorEvents from "./hooks/useSubscribeToConnectorEvents"
import useAutoConnect from "./hooks/useAutoConnect"
import useAutoWalletAddressInjectIntoUrl from "./hooks/useAutoWalletAddressInjectIntoUrl"
import useModalWindowForMobileUsers from "./hooks/useModalWindowForMobileUsers"

const App = () => (
  <Provider store={store}>
    <Router basename={`${process.env.PUBLIC_URL}`}>
      <Messages>
        <Web3ContextProvider>
          <ModalContextProvider>
            <AppLayout />
          </ModalContextProvider>
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
