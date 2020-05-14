import React from "react"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Footer from "./components/Footer"
import Header from "./components/Header"
import Routing from "./components/Routing"
import ContractsDataContextProvider from "./components/ContractsDataContextProvider"
import { Messages } from "./components/Message"
import { SideMenu, SideMenuProvider } from "./components/SideMenu"
import { BrowserRouter as Router } from "react-router-dom"

const App = () => (
  <Messages>
    <Web3ContextProvider>
      <ContractsDataContextProvider>
        <SideMenuProvider>
          <Router>
            <main>
              <Header />
              <aside>
                <SideMenu />
              </aside>
              <content>
                <Routing />
              </content>
              <Footer />
            </main>
          </Router>
        </SideMenuProvider>
      </ContractsDataContextProvider>
    </Web3ContextProvider>
  </Messages>
)

export default App
