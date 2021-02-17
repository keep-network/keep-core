import React, { useEffect } from "react"
import { Link } from "react-router-dom"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Routing from "./components/Routing"
import { Messages } from "./components/Message"
import { SideMenu } from "./components/SideMenu"
import { BrowserRouter as Router } from "react-router-dom"
import { Provider, useDispatch } from "react-redux"
import store from "./store"
import { ModalContextProvider } from "./components/Modal"
import * as Icons from "./components/Icons"
import Footer from "./components/Footer"
import { useWeb3Address } from "./components/WithWeb3Context"
import { usePrevious } from "./hooks/usePrevious"
import { isSameEthAddress } from "./utils/general.utils"

const App = () => (
  <Provider store={store}>
    <Router>
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
  const yourAddress = useWeb3Address()
  const previousAddress = usePrevious(yourAddress)
  const dispatch = useDispatch()

  useEffect(() => {
    if (previousAddress && !isSameEthAddress(previousAddress, yourAddress)) {
      dispatch({ type: "restart_saga" })
    }
  })

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
