import React, { useContext } from 'react'
import { Web3Context } from './WithWeb3Context'
import { ContractsDataContext } from './ContractsDataContextProvider'

const Alerts = (props) => {
    
    const data = useContext(ContractsDataContext)
    const { error } = useContext(Web3Context)
    const { isOperator, isOperatorOfStakedTokenGrant, stakedGrant, stakeOwner } = data

    return (
        <>
        { error && <div className="alert alert-danger m-5" role="alert">{error}</div> }

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
};

export default Alerts