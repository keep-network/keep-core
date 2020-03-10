import React from 'react'
import DelegatedTokens from '../components/DelegatedTokens'
import PendingUndelegation from '../components/PendingUndelegation'
import SlashedTokens from '../components/SlashedTokens'
import AuthorizationInfo from '../components/AuthorizationInfo'
import { useSubscribeToContractEvent } from '../hooks/useSubscribeToContractEvent'
import { TOKEN_STAKING_CONTRACT_NAME_EVENTS } from '../constants/constants'

const OperatorPage = (props) => {
  const { latestEvent } =
    useSubscribeToContractEvent(TOKEN_STAKING_CONTRACT_NAME_EVENTS, 'Undelegated')

  return (
    <>
      <h2 className="mb-2">My Token Operations</h2>
      <DelegatedTokens />
      <PendingUndelegation latestUnstakeEvent={latestEvent} />
      <AuthorizationInfo />
      {/* <SlashedTokens /> */}
    </>

  )
}

export default OperatorPage
