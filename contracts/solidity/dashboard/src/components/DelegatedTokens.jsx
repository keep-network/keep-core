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
  const { isFetching, data: { stakedBalance, ownerAddress, beneficiaryAddress } } = state

  return (
    <section id="delegated-tokens" className="flex flex-row-space-between">
      <LoadingOverlay isFetching={isFetching} >
        <section id="delegated-tokens-summary" className="tile flex flex-column">
          <h5>DelegatedTokens</h5>
          <h2 className="balance">
            {stakedBalance && `${displayAmount(stakedBalance)}`}
          </h2>
          <h6 className="text-darker-grey">OWNER&nbsp;
            <AddressShortcut address={ownerAddress} classNames='text-big text-darker-grey' />
          </h6>
          <h6 className="text-darker-grey">BENEFICIARY&nbsp;
            <AddressShortcut address={beneficiaryAddress} classNames='text-big text-darker-grey' />
          </h6>
        </section>
      </LoadingOverlay>
      <section id="delegated-form-section" className="tile flex flex-column ">
        <h5 className="flex">Undelegate All Tokens</h5>
        <p className="text-warning border flex">
          This will return all of your delegated tokens to the original owner address.
        </p>
        <h6 className="text-label">pending undelegation</h6>
        <UndelegateStakeButton
          btnText='undelegate all my tokens'
          btnClassName="btn btn-default btn-large"
          operator={yourAddress}
        />
      </section>
    </section>
  )
}

export default DelegatedTokens
