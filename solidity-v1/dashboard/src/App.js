import React from "react"
import { Provider } from "react-redux"
import { BrowserRouter as Router, Link } from "react-router-dom"
import store from "./store"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Routing from "./components/Routing"
import { Messages } from "./components/Message"
import { SideMenu } from "./components/SideMenu"
// import { ModalContextProvider } from "./components/Modal"
import * as Icons from "./components/Icons"
import Footer from "./components/Footer"
import useSubscribeToConnectorEvents from "./hooks/useSubscribeToConnectorEvents"
import useAutoConnect from "./hooks/useAutoConnect"
import useAutoWalletAddressInjectIntoUrl from "./hooks/useAutoWalletAddressInjectIntoUrl"
import useModalWindowForMobileUsers from "./hooks/useModalWindowForMobileUsers"
import { ModalRoot } from "./components/modal-component"
import { useModal } from "./hooks/useModal"
import { MODAL_TYPES } from "./constants/constants"
import Button from "./components/Button"

const App = () => (
  <Provider store={store}>
    <Router basename={`${process.env.PUBLIC_URL}`}>
      <Messages>
        <Web3ContextProvider>
          {/* <ModalContextProvider> */}
          <ModalRoot />
          <AppLayout />
          {/* </ModalContextProvider> */}
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
  const { openConfirmationModal } = useModal()

  const a = async () => {
    const result = await openConfirmationModal(MODAL_TYPES.Example, {
      name: "Raf",
    })
    console.log("result", result)
  }

  return (
    <>
      <Button onClick={a}>open confirmation modal</Button>
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
