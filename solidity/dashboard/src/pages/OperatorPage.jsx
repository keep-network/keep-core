import React from 'react'
import DelegatedTokens from '../components/DelegatedTokens'
import PendingUndelegation from '../components/PendingUndelegation'
import SlashedTokens from '../components/SlashedTokens'
import AuthorizationInfo from '../components/AuthorizationInfo'
import { useSubscribeToContractEvent } from '../hooks/useSubscribeToContractEvent'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import PageWrapper from '../components/PageWrapper'

const OperatorPage = (props) => {
  const { latestEvent } =
    useSubscribeToContractEvent(TOKEN_STAKING_CONTRACT_NAME, 'Undelegated')

  return (
    <PageWrapper title="Operations">
      <DelegatedTokens />
      <PendingUndelegation latestUnstakeEvent={latestEvent} />
      <SlashedTokens />
      <AuthorizationInfo />
    </PageWrapper>

  )
}

export default OperatorPage
