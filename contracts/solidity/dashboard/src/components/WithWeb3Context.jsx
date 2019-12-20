import React from 'react'

export const Web3Context = React.createContext({
  web3: null,
  yourAddress: '',
  networkType: '',
  token: { options: { address: '' } },
  stakingContract: { options: { address: '' } },
  grantContract: { options: { address: '' } },
  utils: {},
  eth: {},
  error: '',
  eventToken: { options: { address: '' } },
  eventStakingContract: { options: { address: '' } },
  eventGrantContract: { options: { address: '' } },
})

const withWeb3Context = (Component) => {
  return (props) => (
    <Web3Context.Consumer>
      {({ eventToken, eventStakingContract, eventGrantContract, ...web3 }) => (
        <Component
          {...props}
          web3={web3}
          web3EventProvider={{ eventGrantContract, eventStakingContract, eventToken }}
        />
      )}
    </Web3Context.Consumer>
  )
}

export default withWeb3Context
