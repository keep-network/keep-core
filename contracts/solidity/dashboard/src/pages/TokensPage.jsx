import React, { useContext } from 'react'
import DelegateStakeForm from '../components/DelegateStakeForm'
import TokensOverview from '../components/TokensOverview'
import Undelegations from '../components/Undelegations'
import { Web3Context } from '../components/WithWeb3Context'
import { useShowMessage } from '../components/Message'
import { wait } from '../utils'

const TokensPage = () => {
  const { token } = useContext(Web3Context)
  const showMessage = useShowMessage()


  const handleSubmit = async (values, onTransactionHashCallback) => {
    try {
      console.log('values', values, onTransactionHashCallback)
      await Promise.all([wait(5000)])
      onTransactionHashCallback('hash heheeh')
    } catch (error) {
      console.log('error', error)
    }
  }

  return (
    <React.Fragment>
      <h3>My Tokens</h3>
      <div className="flex flex-1 flex-row-space-between flex-wrap">
        <DelegateStakeForm onSubmit={handleSubmit} minStake={100} availableTokens={3000}/>
        <TokensOverview />
      </div>
      <Undelegations />
    </React.Fragment>
  )
}

export default TokensPage
