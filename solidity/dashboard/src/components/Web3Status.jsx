import React, { useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import Banner, { BANNER_TYPE } from './Banner'

export const Web3Status = () => {
  const {
    web3,
    isFetching,
    yourAddress,
    connectAppWithAccount,
    error,
    provider,
  } = useContext(Web3Context)

  const renderStatus = () => {
    if (!provider) {
      return (
        <Banner
          type={BANNER_TYPE.ERROR}
          title='Please select a wallet'
        />
      )
    }

    if (isFetching) {
      return (
        <Banner
          type={BANNER_TYPE.DISABLED}
          title='Loading ...'
        />
      )
    }

    if (error) {
      return (
        <Banner
          type={BANNER_TYPE.ERROR}
          title={error}
        />
      )
    }

    if (!web3) {
      return (
        <Banner
          type={BANNER_TYPE.ERROR}
          title='Install the MetaMask browser extension'
          subtitle='You can then use the dapp in your current browser.'
        >
          <a
            href="http://metamask.io"
            target="_blank"
            rel="noopener noreferrer"
            className="btn btn-transparent btn-xs ml-1">
              install metamask
          </a>
        </Banner>
      )
    }

    if (!yourAddress) {
      return (
        <Banner
          titleClassName="text-link"
          type={BANNER_TYPE.PENDING}
          title='Please log in and connect with dApp'
          onTitleClick={connectAppWithAccount}
        />
      )
    }

    return (
      <Banner
        type={BANNER_TYPE.SUCCESS}
        title='You are logged in securely to MetaMask.'
        withIcon
      />
    )
  }

  return (
    <div className="web3">
      {renderStatus()}
    </div>
  )
}
