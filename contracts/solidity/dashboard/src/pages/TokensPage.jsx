import React from 'react'
import DelegateStakeForm from '../components/DelegateStakeForm'
import TokensOverview from '../components/TokensOverview'
import Grants from '../components/Grants'
import DelegatedTokens from '../components/DelegatedTokens'
import Undelegations from '../components/Undelegations'

const TokensPage = () => {

  return (
    <React.Fragment>
      <div className="flex flex-1 flex-row-space-between">
        <DelegateStakeForm />
        <TokensOverview />
      </div>
      <div className="flex flex-row">
        <Grants />
        <DelegatedTokens />
      </div>
      <Undelegations />
    </React.Fragment>
  )
}

export default TokensPage
