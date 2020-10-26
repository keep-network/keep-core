import React from "react"
import { Link } from "react-router-dom"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Routing from "./components/Routing"
import ContractsDataContextProvider from "./components/ContractsDataContextProvider"
import { Messages } from "./components/Message"
import { SideMenu } from "./components/SideMenu"
import { BrowserRouter as Router } from "react-router-dom"
import { Provider } from "react-redux"
import store from "./store"
import { ModalContextProvider } from "./components/Modal"
import * as Icons from "./components/Icons"
import Footer from "./components/Footer"

const App = () => (
  <Provider store={store}>
    <Messages>
      <Web3ContextProvider>
        <ModalContextProvider>
          <ContractsDataContextProvider>
            <Router>
              <AppLayout />
            </Router>
          </ContractsDataContextProvider>
        </ModalContextProvider>
      </Web3ContextProvider>
    </Messages>
  </Provider>
)

const AppLayout = () => {
  return (
    <div className="app-layout">
      <div className="app-layout__left">
        <Link to="/" className="app-logo">
          <Icons.KeepDashboardLogo />
        </Link>
        <SideMenu />
        <Footer />
      </div>
      <div className="app-layout__center">
        <Routing />
      </div>
    </div>
  )
}

export default App
