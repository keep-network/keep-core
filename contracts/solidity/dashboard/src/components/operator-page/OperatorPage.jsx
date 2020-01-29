import React from 'react'
import DelegatedTokens from './DelegatedTokens'
import PendingUndelegation from './PendingUndelegation'
import SlashedTokens from './SlashedTokens'

const OperatorPage = (props) => {
  return (
    <>
      <h3>My Token Operations</h3>
      <DelegatedTokens />
      <PendingUndelegation />
      <SlashedTokens />
    </>

  )
}

export default OperatorPage
