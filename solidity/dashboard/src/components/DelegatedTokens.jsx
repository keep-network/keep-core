import React, { useContext } from "react"
import AddressShortcut from "./AddressShortcut"
import { operatorService } from "../services/token-staking.service"
import { useFetchData } from "../hooks/useFetchData"
import { LoadingOverlay } from "./Loadable"
import { Web3Context } from "./WithWeb3Context"
import UndelegateStakeButton from "./UndelegateStakeButton"
import Banner, { BANNER_TYPE } from "./Banner"
import moment from "moment"
import TokenAmount from "./TokenAmount"
import { formatDate } from "../utils/general.utils"

const DelegatedTokens = (props) => {
  const { yourAddress } = useContext(Web3Context)
  const [state, setData] = useFetchData(
    operatorService.fetchDelegatedTokensData,
    {}
  )
  const {
    isFetching,
    data: {
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
    },
  } = state

  const undelegationSuccessCallback = () => {
    setData({ ...state.data, delegationStatus: "UNDELEGATED" })
  }

  const cancelSuccessCallback = () => {
    setData({
      ...state.data,
      stakedBalance: "0",
      delegationStatus: "CANCELED",
    })
  }

  const getUndelegationBannerTitle = () => {
    const undelegationPeriodRelativeTime = moment()
      .add(undelegationPeriod, "seconds")
      .fromNow(true)
    const isUndelegationPeriodOver = moment().isAfter(undelegationCompletedAt)
    const title = isUndelegationPeriodOver
      ? `Completed at ${formatDate(undelegationCompletedAt)}`
      : `Estimated to complete in ${undelegationPeriodRelativeTime}.`

    return `Undelegation is pending. ${title}`
  }

  const renderUndelegationStatus = () => {
    if (delegationStatus) {
      const title =
        delegationStatus === "UNDELEGATED"
          ? getUndelegationBannerTitle()
          : `Delegation ${delegationStatus}`
      let bannerType = BANNER_TYPE.PENDING
      if (delegationStatus === "CANCELED") {
        bannerType = BANNER_TYPE.DISABLED
      } else if (delegationStatus === "RECOVERED") {
        bannerType = BANNER_TYPE.SUCCESS
      }

      return (
        <div className="self-start">
          <Banner type={bannerType} title={title} withIcon />
        </div>
      )
    } else {
      return (
        <UndelegateStakeButton
          btnText="undelegate tokens"
          btnClassName="btn btn-primary btn-lg self-start"
          operator={yourAddress}
          successCallback={
            isInInitializationPeriod
              ? cancelSuccessCallback
              : undelegationSuccessCallback
          }
          isFromGrant={isUndelegationFromGrant}
          isInInitializationPeriod={isInInitializationPeriod}
          isManagedGrant={isManagedGrant}
          managedGrantContractInstance={managedGrantContractInstance}
          disabled={stakedBalance === "0" || !stakedBalance}
        />
      )
    }
  }

  return (
    <LoadingOverlay isFetching={isFetching}>
      <section className="flex row wrap">
        <section className="tile delegation-overview">
          <h2 className="text-grey-70">Total Balance</h2>
          <TokenAmount amount={stakedBalance} />
          <DelegationAddress address={ownerAddress} label={"owner"} />
          <DelegationAddress
            address={beneficiaryAddress}
            label={"beneficiary"}
          />
          <DelegationAddress address={authorizerAddress} label={"authorizer"} />
        </section>
        <section className="tile flex column undelegation-section">
          <h4 className="text-grey-70">Undelegate</h4>
          <div
            className="text-small text-grey-70 mt-1"
            style={{ marginBottom: "auto" }}
          >
            Click undelegate below to return all of your delegated KEEP tokens
            to their original owner address.
          </div>
          {renderUndelegationStatus()}
        </section>
      </section>
    </LoadingOverlay>
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
