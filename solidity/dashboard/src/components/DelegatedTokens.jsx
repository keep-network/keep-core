import React, { useContext } from "react"
import AddressShortcut from "./AddressShortcut"
import { operatorService } from "../services/token-staking.service"
import { useFetchData } from "../hooks/useFetchData"
import { LoadingOverlay } from "./Loadable"
import { displayAmount } from "../utils/token.utils"
import { Web3Context } from "./WithWeb3Context"
import UndelegateStakeButton from "./UndelegateStakeButton"
import { PENDING_STATUS, COMPLETE_STATUS } from "../constants/constants"
import Banner, { BANNER_TYPE } from "./Banner"
import moment from "moment"
import Tile from "./Tile"

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
          btnText="undelegate"
          btnClassName="btn btn-primary btn-lg flex-1"
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
      <section id="delegated-tokens" className="flex row space-between">
        <Tile
          title="Total Balance"
          id="delegated-tokens-summary"
          className="tile flex column"
        >
          <h2 className="balance">
            {stakedBalance && `${displayAmount(stakedBalance)}`} KEEP
          </h2>
          <h6 className="text-grey-70">
            owner&nbsp;
            <AddressShortcut
              address={ownerAddress}
              classNames="text-small text-normal text-darker-grey"
            />
          </h6>
          <h6 className="text-grey-70">
            beneficiary&nbsp;
            <AddressShortcut
              address={beneficiaryAddress}
              classNames="text-small text-normal text-darker-grey"
            />
          </h6>
          <h6 className="text-grey-70">
            authorizer&nbsp;
            <AddressShortcut
              address={authorizerAddress}
              classNames="text-small text-normal text-darker-grey"
            />
          </h6>
        </Tile>
        <Tile
          title="Undelegate All Tokens"
          id="delegated-form-section"
          className="tile flex column "
        >
          <div className="text-big text-grey-70 mt-1 mb-1">
            Click undelegate below to return all of your delegated KEEP tokens
            to their original owner address.
          </div>
          <div className="flex" style={{ marginTop: "auto" }}>
            {renderUndelegationStatus()}
          </div>
        </Tile>
      </section>
    </LoadingOverlay>
  )
}

export default DelegatedTokens
