import React, { useContext } from "react"
import AddressShortcut from "./AddressShortcut"
import { Web3Context } from "./WithWeb3Context"
import UndelegateStakeButton from "./UndelegateStakeButton"
import Banner from "./Banner"
import moment from "moment"
import TokenAmount from "./TokenAmount"
import { formatDate } from "../utils/general.utils"
import * as Icons from "./Icons"

const DelegatedTokens = ({ data, setData }) => {
  const { yourAddress } = useContext(Web3Context)

  const {
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
    isUndelegationFromGrant,
    isInInitializationPeriod,
    undelegationPeriod,
    isManagedGrant,
    managedGrantContractInstance,
    delegationStatus,
    undelegationCompletedAt,
  } = data

  const cancelSuccessCallback = () => {
    setData({
      ...data,
      stakedBalance: "0",
      delegationStatus: "CANCELED",
    })
  }

  const getUndelegationBannerTitle = () => {
    const undelegationPeriodRelativeTime = moment()
      .add(undelegationPeriod, "seconds")
      .fromNow(true)
    const isUndelegationPeriodOver = moment().isAfter(undelegationCompletedAt)
    return isUndelegationPeriodOver
      ? `Undelegation period is over at ${formatDate(undelegationCompletedAt)}`
      : `Undelegation is pending. Estimated to complete in ${undelegationPeriodRelativeTime}.`
  }

  const renderUndelegationStatus = () => {
    if (delegationStatus) {
      const title =
        delegationStatus === "UNDELEGATED"
          ? getUndelegationBannerTitle()
          : "Undelegation completed"

      const bannerClassName =
        delegationStatus === "UNDELEGATED" ? "bg-pending" : "bg-success"

      return (
        <div className="self-start">
          <Banner
            inline
            title={title}
            className={bannerClassName}
            icon={Icons.Time}
          />
        </div>
      )
    } else {
      return (
        <UndelegateStakeButton
          btnText="undelegate tokens"
          btnClassName="btn btn-primary btn-lg self-start"
          operator={yourAddress}
          successCallback={
            isInInitializationPeriod ? cancelSuccessCallback : () => {}
          }
          isFromGrant={isUndelegationFromGrant}
          isInInitializationPeriod={isInInitializationPeriod}
          isManagedGrant={isManagedGrant}
          managedGrantContractInstance={managedGrantContractInstance}
          disabled={stakedBalance === "0" || !stakedBalance}
          undelegationPeriod={undelegationPeriod}
        />
      )
    }
  }

  return (
    <section className="flex row wrap">
      <section className="tile delegation-overview">
        <h2 className="text-grey-70">Total Balance</h2>
        <TokenAmount amount={stakedBalance} />
        <DelegationAddress address={ownerAddress} label={"owner"} />
        <DelegationAddress address={beneficiaryAddress} label={"beneficiary"} />
        <DelegationAddress address={authorizerAddress} label={"authorizer"} />
      </section>
      <section className="tile flex column undelegation-section">
        <h4 className="text-grey-70">Undelegate</h4>
        <div
          className="text-small text-grey-70 mt-1"
          style={{ marginBottom: "auto" }}
        >
          Click undelegate below to return all of your delegated KEEP tokens to
          their original owner address.
        </div>
        {renderUndelegationStatus()}
      </section>
    </section>
  )
}

const DelegationAddress = React.memo(({ address, label }) => (
  <h6 className="text-grey-50" style={{ marginTop: "0.5rem" }}>
    {label}&nbsp;
    <AddressShortcut
      address={address}
      classNames="h6 text-normal text-grey-50"
    />
  </h6>
))

export default DelegatedTokens
