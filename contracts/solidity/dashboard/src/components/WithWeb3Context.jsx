
import React from 'react'

export const Web3Context = React.createContext({ 
  yourAddress: '',
  networkType: '',
  token: { options: { address: '' } },
  stakingContract: { options: { address: '' } },
  grantContract: { options: { address: '' } },
  utils: {},
  eth: {},
  error: '',
  dataIsReady: false,
});

const WithWeb3Context = (Component) => {
  return (props) => (
    <Web3Context.Consumer>
      {web3 =>  <Component {...props} web3={web3} />}
    </Web3Context.Consumer>
  )
}

export default WithWeb3Context