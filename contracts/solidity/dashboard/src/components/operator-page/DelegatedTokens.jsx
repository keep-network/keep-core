import React, { useEffect, useContext } from 'react'
import AddressShortcut from '../AddressShortcut'
import InlineForm from '../InlineForm'
import { operatorService } from './service'
import { useFetchData } from '../../hooks/useFetchData'
import { LoadingOverlay } from '../Loadable'
import { displayAmount } from '../../utils'
import { Web3Context } from '../WithWeb3Context'

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
    <section id="delegated-tokens" className="tile">
      <LoadingOverlay isFetching={isFetching}>
        <h5>Total Delegated Tokens</h5>
        <div className="flex flex-row">
          <div className="delegated-tokens-summary flex flex-column" style={{ flex: '1' }} >
            <h2 className="balance">
              {stakedBalance && `${displayAmount(stakedBalance)} K`}
            </h2>
            <div>
              <h6 className="text-darker-grey">OWNER&nbsp;
                <AddressShortcut address={ownerAddress} classNames='text-big text-darker-grey' />
              </h6>
              <h6 className="text-darker-grey">BENEFICIARY&nbsp;
                <AddressShortcut address={beneficiaryAddress} classNames='text-big text-darker-grey' />
              </h6>
            </div>
          </div>
          <InlineForm inputProps={{ placeholder: 'Amount' }} classNames="undelegation-form" />
        </div>
      </LoadingOverlay>
    </section>
  )
}

export default DelegatedTokens
