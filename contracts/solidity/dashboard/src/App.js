import React from 'react'
import Web3ContextProvider from './components/Web3ContextProvider'
import Footer from './components/Footer'
import Header from './components/Header'
import Routing from './components/Routing'
import ContractsDataContextProvider from './components/ContractsDataContextProvider'
import { Messages } from './components/Message'
import { SideMenu, SideMenuProvider } from './components/SideMenu'
import { BrowserRouter as Router } from 'react-router-dom'

const App = () => (
  <Messages>
    <Web3ContextProvider>
      <ContractsDataContextProvider>
        <SideMenuProvider>
          <Router>
            <div className='main'>
              <Header />
              <SideMenu />
              <div className='content'>
                <Routing />
                <Footer />
              </div>
            </div>
          </Router>
        </SideMenuProvider>
      </ContractsDataContextProvider>
    </Web3ContextProvider>
  </Messages>
)

export default App
