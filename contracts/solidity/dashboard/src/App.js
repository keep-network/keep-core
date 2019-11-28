import React, { Fragment } from 'react'
import Web3ContextProvider from './components/Web3ContextProvider';
import Footer from './components/Footer';
import Header from './components/Header';
import Routing from './components/Routing'
import ContractsDataContextProvider from './components/ContractsDataContextProvider';

const App = () => (
  <Web3ContextProvider>
     <div className='main'>
      <Fragment>
        <Header />
        <ContractsDataContextProvider>
          <Routing />
        </ContractsDataContextProvider>
        <Footer />
      </Fragment>
    </div>
  </Web3ContextProvider>
)

export default App