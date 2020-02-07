import React, { useEffect, useContext } from 'react'
import AddressShortcut from './AddressShortcut'
import { operatorService } from '../services/token-staking.service'
import { useFetchData } from '../hooks/useFetchData'
import { LoadingOverlay } from './Loadable'
import { displayAmount } from '../utils'
import { Web3Context } from './WithWeb3Context'
import UndelegateForm from './UndelegateForm'

const DelegatedTokens = ({ latestUnstakeEvent }) => {
  const { utils } = useContext(Web3Context)
  const [state, setData] = useFetchData(operatorService.fetchDelegatedTokensData, {})
  const { isFetching, data: { stakedBalance, ownerAddress, beneficiaryAddress } } = state

  useEffect(() => {
    if (latestUnstakeEvent) {
      const { returnValues: { value } } = latestUnstakeEvent
      const updatedStakeBalance = utils.toBN(stakedBalance).sub(utils.toBN(value))
      setData({ stakedBalance: updatedStakeBalance, ownerAddress, beneficiaryAddress })
    }
  }, [latestUnstakeEvent])

  return (
    <section id="delegated-tokens" className="flex flex-row-space-between">
      <LoadingOverlay isFetching={isFetching} >
        <section id="delegated-tokens-summary" className="tile flex flex-column">
          <h5>My Delegated Tokens</h5>
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
        <h5 className="flex flex-1">Undelegate Tokens</h5>
        <p className="text-warning border flex flex-1">
          Starting an undelegation restarts the amount of time, or undelegation period, until tokens are returned to the owner
        </p>
        <UndelegateForm />
      </section>
    </section>
  )
}

export default DelegatedTokens
