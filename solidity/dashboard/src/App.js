import React from "react"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Header, { HeaderContextProvider } from "./components/Header"
import Routing from "./components/Routing"
import ContractsDataContextProvider from "./components/ContractsDataContextProvider"
import { Messages } from "./components/Message"
import { SideMenu, SideMenuProvider } from "./components/SideMenu"
import { BrowserRouter as Router } from "react-router-dom"
import { Provider } from "react-redux"
import store from "./store"
import { ModalContextProvider } from "./components/Modal"

const App = () => (
  <Provider store={store}>
    <Messages>
      <Web3ContextProvider>
        <ModalContextProvider>
          <ContractsDataContextProvider>
            <HeaderContextProvider>
              <SideMenuProvider>
                <Router>
                  <main>
                    <Header />
                    <SideMenu />
                    <div className="content-wrapper">
                      <div className="content">
                        <Routing />
                      </div>
                    </div>
                  </main>
                </Router>
              </SideMenuProvider>
            </HeaderContextProvider>
          </ContractsDataContextProvider>
        </ModalContextProvider>
      </Web3ContextProvider>
    </Messages>
  </Provider>
)

export default App
