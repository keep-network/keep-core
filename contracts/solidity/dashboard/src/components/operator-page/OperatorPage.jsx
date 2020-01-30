import React, { useContext } from 'react'
import DelegatedTokens from './DelegatedTokens'
import PendingUndelegation from './PendingUndelegation'
import SlashedTokens from './SlashedTokens'
import { useSubscribeToContractEvent } from '../../hooks/useSubscribeToContractEvent'
import { TOKEN_STAKING_CONTRACT_NAME } from '../../constants/constants'
import { Web3Context } from '../WithWeb3Context'

const OperatorPage = (props) => {
  const { yourAddress } = useContext(Web3Context)
  const { latestEvent } =
    useSubscribeToContractEvent(TOKEN_STAKING_CONTRACT_NAME, 'InitiatedUnstake', { filter: { operator: yourAddress } })

  return (
    <>
      <h3>My Token Operations</h3>
      <DelegatedTokens latestUnstakeEvent={latestEvent} />
      <PendingUndelegation latestUnstakeEvent={latestEvent} />
      <SlashedTokens />
    </>

  )
}

export default OperatorPage
