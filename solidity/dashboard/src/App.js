import React from "react"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Footer from "./components/Footer"
import Header from "./components/Header"
import Routing from "./components/Routing"
import ContractsDataContextProvider from "./components/ContractsDataContextProvider"
import { Messages } from "./components/Message"
import { SideMenu, SideMenuProvider } from "./components/SideMenu"
import { BrowserRouter as Router, Route, Switch } from "react-router-dom"
import CopyStakePage from "./pages/CopyStakePage"
import { ModalContextProvider } from "./components/Modal"

const App = () => (
  <Messages>
    <Web3ContextProvider>
      <ModalContextProvider>
        <ContractsDataContextProvider>
          <SideMenuProvider>
            <Router>
              <Switch>
                <Route exact path="/copy-stake">
                  <CopyStakePage />
                </Route>
                <Route path="/">
                  <main>
                    <Header />
                    <aside>
                      <SideMenu />
                    </aside>
                    <div className="content">
                      <Routing />
                    </div>
                    <Footer />
                  </main>
                </Route>
              </Switch>
            </Router>
          </SideMenuProvider>
        </ContractsDataContextProvider>
      </ModalContextProvider>
    </Web3ContextProvider>
  </Messages>
)

export default App
