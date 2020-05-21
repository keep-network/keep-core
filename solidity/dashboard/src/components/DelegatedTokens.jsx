import React, { useContext } from "react"
import AddressShortcut from "./AddressShortcut"
import { operatorService } from "../services/token-staking.service"
import { useFetchData } from "../hooks/useFetchData"
import { LoadingOverlay } from "./Loadable"
import { Web3Context } from "./WithWeb3Context"
import UndelegateStakeButton from "./UndelegateStakeButton"
import { PENDING_STATUS, COMPLETE_STATUS } from "../constants/constants"
import Banner, { BANNER_TYPE } from "./Banner"
import moment from "moment"
import TokenAmount from "./TokenAmount"

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
      undelegationStatus,
      isUndelegationFromGrant,
      isInInitializationPeriod,
      undelegationPeriod,
      isManagedGrant,
      managedGrantContractInstance,
    },
  } = state

  const undelegationSuccessCallback = () => {
    setData({ ...state.data, undelegationStatus: PENDING_STATUS })
  }

  const cancelSuccessCallback = () => {
    setData({
      ...state.data,
      stakedBalance: "0",
      undelegationStatus: COMPLETE_STATUS,
    })
  }

  const renderUndelegationStatus = () => {
    if (undelegationStatus === PENDING_STATUS) {
      const undelegationPeriodRelativeTime = moment()
        .add(undelegationPeriod, "seconds")
        .fromNow(true)
      const title = `Undelegation is pending. Estimated to complete in ${undelegationPeriodRelativeTime}.`
      return (
        <div className="self-start">
          <Banner type={BANNER_TYPE.PENDING} title={title} withIcon />
        </div>
      )
    } else if (undelegationStatus === COMPLETE_STATUS) {
      return (
        <div className="self-start">
          <Banner
            type={BANNER_TYPE.SUCCESS}
            title="Undelegation completed"
            withIcon
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
            isInInitializationPeriod
              ? cancelSuccessCallback
              : undelegationSuccessCallback
          }
          isFromGrant={isUndelegationFromGrant}
          isInInitializationPeriod={isInInitializationPeriod}
          isManagedGrant={isManagedGrant}
          managedGrantContractInstance={managedGrantContractInstance}
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
