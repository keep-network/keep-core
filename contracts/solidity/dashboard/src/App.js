import React, { Fragment } from 'react'
import Web3ContextProvider from './components/Web3ContextProvider';
import Footer from './components/Footer';
import Header from './components/Header';
import Routing from './components/Routing'

const App = () => (
  <Web3ContextProvider>
     <div className='main'>
      <Fragment>
        <Header />
        <Routing />
        <Footer />
      </Fragment>
    </div>
  </Web3ContextProvider>
)

export default App