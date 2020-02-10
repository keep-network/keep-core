import React, { useContext } from 'react'
import { Web3Context } from './WithWeb3Context'

export const Web3Status = (props) => {
  const { web3, isFetching, yourAddress, connectAppWithAccount, error } = useContext(Web3Context)

  const renderStatus = () => {
    if (isFetching) {
      return (
        <div className="web3-status loading">
          Loading...
        </div>
      )
    }

    if (error) {
      return (
        <div className="web3-status alert">
          {error}
        </div>
      )
    }

    if (!web3) {
      return (
        <div className="web3-status alert flex">
          <div>
            <div className="title">
              Install the MetaMask browser extension
            </div>
            <div className="sub-title">
              You can then use the dapp in your current browser.
            </div>
          </div>
          <a href="http://metamask.io" target="_blank" rel="noopener noreferrer" className="btn btn-primary btn-sm text-white">INSTALL METAMASK</a>
        </div>
      )
    }

    if (!yourAddress) {
      return (
        <div className="web3-status notify">
          <span onClick={connectAppWithAccount}>
            Please log in and connect with dApp
          </span>
        </div>
      )
    }

    return (
      <div className="web3-status success">
        You are logged in securely to MetaMask.
      </div>
    )
  }

  return (
    <div className="web3">
      {renderStatus()}
    </div>
  )
}
