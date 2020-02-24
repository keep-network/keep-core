import React, { useContext } from 'react'
import AddressShortcut from './AddressShortcut'
import { operatorService } from '../services/token-staking.service'
import { useFetchData } from '../hooks/useFetchData'
import { LoadingOverlay } from './Loadable'
import { displayAmount } from '../utils'
import { Web3Context } from './WithWeb3Context'
import UndelegateStakeButton from './UndelegateStakeButton'

const DelegatedTokens = (props) => {
  const { yourAddress } = useContext(Web3Context)
  const [state] = useFetchData(operatorService.fetchDelegatedTokensData, {})
  const { isFetching, data: {
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
  } } = state

  return (
    <section id="delegated-tokens" className="flex flex-row-space-between">
      <LoadingOverlay isFetching={isFetching} >
        <section id="delegated-tokens-summary" className="tile flex flex-column">
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
        </section>
      </LoadingOverlay>
      <section id="delegated-form-section" className="tile flex flex-column ">
        <h3 className="text-darker-grey">Undelegate All Tokens</h3>
        <div className="text-big text-darker-grey mt-1 mb-1">
          Click undelegate below to return all of your delegated KEEP tokens to their original owner address.
        </div>
        <UndelegateStakeButton
          btnText='undelegate all my tokens'
          btnClassName="btn btn-primary btn-large"
          operator={yourAddress}
        />
      </section>
    </section>
  )
}

export default DelegatedTokens
