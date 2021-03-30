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
import { useWeb3Context } from "./components/WithWeb3Context"

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
  const { isConnected, connector, yourAddress } = useWeb3Context()
  const dispatch = useDispatch()

  useEffect(() => {
    const accountChangedHandler = (address) => {
      dispatch({ type: "app/account_changed", payload: { address } })
    }

    const disconnectHandler = () => {
      dispatch({ type: "app/logout" })
    }

    if (isConnected && connector) {
      dispatch({ type: "app/login", payload: { address: yourAddress } })
      connector.on("accountsChanged", accountChangedHandler)
      connector.once("disconnect", disconnectHandler)
    }

    return () => {
      if (connector) {
        connector.removeListener("accountsChanged", accountChangedHandler)
        connector.removeListener("disconnect", disconnectHandler)
      }
    }
  }, [isConnected, connector, dispatch, yourAddress])

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
