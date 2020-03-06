import React, { useContext } from 'react'
import AddressShortcut from './AddressShortcut'
import { operatorService } from '../services/token-staking.service'
import { useFetchData } from '../hooks/useFetchData'
import { LoadingOverlay } from './Loadable'
import { displayAmount } from '../utils'
import { Web3Context } from './WithWeb3Context'
import UndelegateStakeButton from './UndelegateStakeButton'
import RecoverStakeButton from './RecoverStakeButton'

const DelegatedTokens = (props) => {
  const { yourAddress } = useContext(Web3Context)
  const [state, setData] = useFetchData(operatorService.fetchDelegatedTokensData, {})
  const { isFetching, data: {
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
    undelegationStatus,
  } } = state

  const undelegationSuccessCallback = () => {
    setData({ ...state.data, undelegationStatus: 'PENDING' })
  }

  const renderUndelegationStatus = () => {
    if (undelegationStatus === 'PENDING') {
      // TODO create and render notification component in the future PR.
      return (
        <div className="text-warning text-normal text-bg-pending-light self-start">
          Undelegation is pending
        </div>
      )
    } else if (undelegationStatus === 'COMPLETED') {
      // TODO create and render notification component in the future PR.
      return (
        <div className="text-success text-normal text-bg-success-light self-start">
          Undelegation completed
        </div>
      )
    } else {
      return (
        <UndelegateStakeButton
          btnText='undelegate'
          btnClassName="btn btn-primary btn-large flex-1"
          operator={yourAddress}
          successCallback={undelegationSuccessCallback}
        />
      )
    }
  }
  return (
    <section id="delegated-tokens" className="flex flex-row-space-between">
      <section id="delegated-tokens-summary" className="tile flex flex-column">
        <LoadingOverlay isFetching={isFetching} className >
          <h3 className="text-darker-grey">Delegated Tokens</h3>
          <h2 className="balance">
            {stakedBalance && `${displayAmount(stakedBalance)}`} KEEP
          </h2>
          <h6 className="text-darker-grey">owner&nbsp;
            <AddressShortcut
              address={ownerAddress}
              classNames='text-small text-normal text-darker-grey'
            />
          </h6>
          <h6 className="text-darker-grey">beneficiary&nbsp;
            <AddressShortcut
              address={beneficiaryAddress}
              classNames='text-small text-normal text-darker-grey'
            />
          </h6>
          <h6 className="text-darker-grey">authorizer&nbsp;
            <AddressShortcut
              address={authorizerAddress}
              classNames='text-small text-normal text-darker-grey'
            />
          </h6>
        </LoadingOverlay>
      </section>
      <section id="delegated-form-section" className="tile">
        <LoadingOverlay isFetching={isFetching} className="flex flex-column">
          <h3 className="text-darker-grey">Undelegate All Tokens</h3>
          <div className="text-big text-darker-grey mt-1 mb-1">
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
