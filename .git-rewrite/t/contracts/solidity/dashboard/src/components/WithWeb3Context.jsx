
import React from 'react'

export const Web3Context = React.createContext()

const WithWeb3Context = (Component) => {
  return (props) => (
    <Web3Context.Consumer>
      {web3 =>  <Component {...props} web3={web3} />}
    </Web3Context.Consumer>
  )
}

export default WithWeb3Context