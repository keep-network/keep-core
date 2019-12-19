import React from 'react'
import { Row, Col } from 'react-bootstrap'
import StakingForm from './StakingForm'
import StakingDelegateForm from './StakingDelegateForm'
import { withContractsDataContext } from './ContractsDataContextProvider';

const StakeTab = (props) => {
    const { tokenBalance } = props;
    return(
        <>
            <h3>Stake Delegation</h3>
            <p>
                Keep network does not require token owners to perform the day-to-day operations of staking 
                with the private keys holding the tokens. This is achieved by stake delegation, where different
                addresses hold different responsibilities and cold storage is supported to the highest extent practicable.
            </p>
            <StakingDelegateForm tokenBalance={tokenBalance} />
            <hr></hr>
            <h3>Stake Delegation (Simplified)</h3>
            <p>
                Simplified arrangement where you operate and receive rewards under one account.
            </p>
            <StakingForm btnText="Stake" action="stake" />
        </>
    )
}

export default withContractsDataContext(StakeTab)