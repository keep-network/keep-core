import React, { useContext } from 'react'
import { ContractsDataContext } from './ContractsDataContextProvider'

const Alerts = (props) => {
  const { isOperator, isOperatorOfStakedTokenGrant, stakedGrant, stakeOwner } = useContext(ContractsDataContext)

  return (
    <>
      {isOperator && !isOperatorOfStakedTokenGrant &&
            <div className="alert alert-info m-5" role="alert">You are registered as an operator for {stakeOwner}</div>
      }

      {isOperatorOfStakedTokenGrant &&
            <div className="alert alert-info m-5" role="alert">
              You are registered as a staked token grant operator for {stakedGrant.grantee}
            </div>
      }
    </>
  )
}

export default Alerts
