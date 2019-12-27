import React from 'react'
import Web3ContextProvider from './components/Web3ContextProvider'
import Footer from './components/Footer'
import Header from './components/Header'
import Routing from './components/Routing'
import ContractsDataContextProvider from './components/ContractsDataContextProvider'
import { Messages } from './components/Message'

const App = () => (
  <Messages>
    <Web3ContextProvider>
      <div className='main'>
        <Header />
        <ContractsDataContextProvider>
          <Routing />
        </ContractsDataContextProvider>
        <Footer />
      </div>
    </Web3ContextProvider>
  </Messages>
)

export default App
