import React from 'react'
import SigningForm from './SigningForm'
import { withContractsDataContext } from './ContractsDataContextProvider'
import { Redirect } from 'react-router-dom'

const Siginig = (props) => {

    if(props.contractsDataIsFetching)
        return (<div>Loading...</div>)

    return (props.isOperator || props.isTokenHolder) ? <Redirect to='/overview' /> :
        <div className="signing">
            <div className="alert alert-info m-5" role="alert">Sorry, looks like you don't have any tokens to stake.</div>
                <h3>Become an operator</h3>
                <p>
                    To become an operator you must have a mutual agreement with the stake owner. This is achieved by creating
                    a signature of the stake owner address and sending it to the owner. Using the signature the owner can initiate
                    stake delegation and you will be able to participate in network operations on behalf of the stake owner.
                </p>
            <div className="signing-form well">
                <SigningForm description="Sign stake owner address" defaultMessageToSign="0x0" />
                <SigningForm
                    description="(Optional) Sign Token Grant contract address. This is required only for Token Grants stake operators"
                    defaultMessageToSign={props.message}
                />
            </div>
        </div>
}

export default withContractsDataContext(Siginig)