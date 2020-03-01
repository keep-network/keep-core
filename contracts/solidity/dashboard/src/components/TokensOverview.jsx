import React from 'react'
import { displayAmount } from '../utils'
import TokenGrantOverview from './TokenGrantOverview'

const TokensOverview = ({
  undelegationPeriod,
  keepBalance,
  stakingBalance,
  pendingUndelegationBalance,
  grantBalance,
  tokenGrantsStakeBalance,
}) => {
  return (
    <section id="tokens-overview" className="tile">
      <section>
        <h4 className="text-grey-60">Owned Tokens</h4>
        <h2 className="balance">{displayAmount(keepBalance)}</h2>
        <div className="text-samll">
          Staked Owned Tokens: {displayAmount(stakingBalance)}
          <p className="text-smaller text-grey-30">Tokens you own that are delegated to an operator and doing work on the network.</p>
        </div>
        <div className="text-samll">
          Pending Undelegated Tokens: {displayAmount(pendingUndelegationBalance)}
          <p className="text-smaller text-grey-30">Stake undelegated from an operator. Estimated {undelegationPeriod} number of blocks until available.</p>
        </div>
        <hr />
      </section>
      <section>
        <h4 className="text-grey-60">Granted Tokens</h4>
        <h2 className="balance">{displayAmount(grantBalance)}</h2>
        <TokenGrantOverview />
      </section>
    </section>
  )
}

export default TokensOverview
