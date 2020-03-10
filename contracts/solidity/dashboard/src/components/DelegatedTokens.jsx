import React, { useContext } from 'react'
import AddressShortcut from './AddressShortcut'
import { operatorService } from '../services/token-staking.service'
import { useFetchData } from '../hooks/useFetchData'
import { LoadingOverlay } from './Loadable'
import { displayAmount } from '../utils/general.utils'
import { Web3Context } from './WithWeb3Context'
import UndelegateStakeButton from './UndelegateStakeButton'
import { PENDING_STATUS, COMPLETE_STATUS } from '../constants/constants'
import Banner, { BANNER_TYPE } from './Banner'

const DelegatedTokens = (props) => {
  const { yourAddress } = useContext(Web3Context)
  const [state, setData] = useFetchData(operatorService.fetchDelegatedTokensData, {})
  const { isFetching, data: {
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
    undelegationStatus,
    isUndelegationFromGrant,
  } } = state

  const undelegationSuccessCallback = () => {
    setData({ ...state.data, undelegationStatus: PENDING_STATUS })
  }

  const renderUndelegationStatus = () => {
    if (undelegationStatus === PENDING_STATUS) {
      return (
        <div className="self-start">
          <Banner
            type={BANNER_TYPE.PENDING}
            title='Undelegation is pending'
          />
        </div>
      )
    } else if (undelegationStatus === COMPLETE_STATUS) {
      return (
        <div className="self-start">
          <Banner
            type={BANNER_TYPE.SUCCESS}
            title='Undelegation completed'
            withIcon
          />
        </div>
      )
    } else {
      return (
        <UndelegateStakeButton
          btnText='undelegate'
          btnClassName="btn btn-primary btn-lg flex-1"
          operator={yourAddress}
          successCallback={undelegationSuccessCallback}
          isFromGrant={isUndelegationFromGrant}
        />
      )
    }
  }
  return (
    <section id="delegated-tokens" className="flex row space-between">
      <section id="delegated-tokens-summary" className="tile flex column">
        <LoadingOverlay isFetching={isFetching} >
          <h3 className="text-grey-60 mb-1">Delegated Tokens</h3>
          <h2 className="balance">
            {stakedBalance && `${displayAmount(stakedBalance)}`} KEEP
          </h2>
          <h6 className="text-grey-70">owner&nbsp;
            <AddressShortcut
              address={ownerAddress}
              classNames='text-small text-normal text-darker-grey'
            />
          </h6>
          <h6 className="text-grey-70">beneficiary&nbsp;
            <AddressShortcut
              address={beneficiaryAddress}
              classNames='text-small text-normal text-darker-grey'
            />
          </h6>
          <h6 className="text-grey-70">authorizer&nbsp;
            <AddressShortcut
              address={authorizerAddress}
              classNames='text-small text-normal text-darker-grey'
            />
          </h6>
        </LoadingOverlay>
      </section>
      <section id="delegated-form-section" className="tile flex column ">
        <LoadingOverlay isFetching={isFetching} classNames="flex flex-1 column" >
          <h3 className="text-grey-60">Undelegate All Tokens</h3>
          <div className="text-big text-grey-70 mt-1 mb-1">
            Click undelegate below to return all of your delegated KEEP tokens to their original owner address.
          </div>
          <div className="flex" style={{ marginTop: 'auto' }}>
            {renderUndelegationStatus()}
          </div>
        </LoadingOverlay>
      </section>
    </section>
  )
}

export default DelegatedTokens
