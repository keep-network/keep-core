import React from 'react'
import { displayAmount } from '../utils'

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
      <h5>Totals</h5>
      <section>
        <span className="text-label text-white text-bg-warning">OWNED</span>
        <h2 className="balance">{displayAmount(keepBalance)}</h2>
        <div className="text-samll">
          Staked Owned Tokens: {displayAmount(stakingBalance)}
          <p className="text-smaller text-grey">Tokens you own that are delegated to an operator and doing work on the network.</p>
        </div>
        <div className="text-samll">
          Pending Undelegated Tokens: {displayAmount(pendingUndelegationBalance)}
          <p className="text-smaller text-grey">Stake undelegated from an operator. Estimated {undelegationPeriod} number of blocks until available.</p>
        </div>
        <hr />
      </section>
      <section>
        <span className="text-label text-white text-bg-warning">GRANTED</span>
        <h2 className="balance">{displayAmount(grantBalance)}</h2>
        <div className="text-samll">
          Staked Tokens: {displayAmount(tokenGrantsStakeBalance)}
          <p className="text-smaller text-grey">Tokens you were granted that are delegated to an operator and doing work on the network.</p>
        </div>
      </section>
    </section>
  )
}

export default TokensOverview
