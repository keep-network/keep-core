import React, { useContext } from "react"
import { Web3Context } from "./WithWeb3Context"
import Banner, { BANNER_TYPE } from "./Banner"
import { WALLETS } from "../constants/constants"

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
    if (!provider || !web3) {
      return <Banner type={BANNER_TYPE.ERROR} title="Please select a wallet" />
    }

    if (isFetching) {
      return <Banner type={BANNER_TYPE.DISABLED} title="Loading ..." />
    }

    if (error) {
      return <Banner type={BANNER_TYPE.ERROR} title={error} />
    }

    if (!yourAddress) {
      return (
        <Banner
          titleClassName="text-link"
          type={BANNER_TYPE.PENDING}
          title="Please log in and connect with dApp"
          onTitleClick={connectAppWithAccount}
        />
      )
    }

    return (
      <Banner
        type={BANNER_TYPE.SUCCESS}
        title={`You are logged into ${WALLETS[provider].label}.`}
        withIcon
      />
    )
  }

  return <div className="web3">{renderStatus()}</div>
}
