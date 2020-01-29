import React from 'react'
import DelegatedTokens from './DelegatedTokens'
import PendingUndelegation from './PendingUndelegation'

const OperatorPage = (props) => {

    return (
        <>
            <DelegatedTokens />
            <PendingUndelegation />
        </>
        
    )
}