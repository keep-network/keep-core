import React, { useContext } from 'react'
import { Web3Context } from "./WithWeb3Context"

export const Web3Status = (props) => {
    const { web3, isFetching, yourAddress, connectAppWithAccount, error } = useContext(Web3Context)

    const renderStatus = () => {
        if (isFetching) {
            return (
                <div className="web3-status loading">
                    loading...
                </div>
            )
        }
        
        if(error) {
            return (
                <div className="web3-status alert">
                    {error}
                </div>
            )
        }
        
        if (!web3) {
            return (
                <div className="web3-status alert">
                    Web3 not detected. We suggest&nbsp;<a href="http://metamask.io" target="_blank" rel="noopener noreferrer">MetaMask</a>.
                </div>
            )
        }
        
        if (!yourAddress) {
            return (
                <div className="web3-status notify">
                    <span onClick={connectAppWithAccount}>
                        Please log in and connect with dApp
                    </span>
                </div>
            )
        }

        return (
            <div className="web3-status success">
                Account logged in
            </div>
        )
    }

    return (
        <div className="web3">
            {renderStatus()}
        </div>
    )
}