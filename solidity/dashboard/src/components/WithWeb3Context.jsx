import React, { useContext } from "react"

export const Web3Context = React.createContext({
  web3: null,
  // Points to the currently selected account and it's used as a seneder of the transaction.
  yourAddress: "",
  networkType: "",
  token: { options: { address: "" } },
  stakingContract: { options: { address: "" } },
  grantContract: { options: { address: "" } },
  utils: {},
  eth: {},
  error: "",
})

export const useWeb3Context = () => {
  const web3Context = useContext(Web3Context)

  if (!web3Context) {
    throw new Error("Web3Context not found")
  }

  return web3Context
}

export const useWeb3Address = () => {
  return useWeb3Context().yourAddress
}

const withWeb3Context = (Component) => {
  return (props) => (
    <Web3Context.Consumer>
      {(web3) => <Component {...props} web3={web3} />}
    </Web3Context.Consumer>
  )
}

export default withWeb3Context
