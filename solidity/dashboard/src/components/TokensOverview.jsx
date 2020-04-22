import React from 'react'
import { displayAmount } from '../utils/general.utils'
import TokenGrantsOverview from './TokenGrantsOverview'
import moment from 'moment'

const TokensOverview = ({
  undelegationPeriod,
  keepBalance,
  delegatedTokens,
  stakingBalance,
  pendingUndelegationBalance,
}) => {
  const estimatedUndelegationPeriod = moment().add(undelegationPeriod, 'seconds').fromNow(true)

  return (
    <section id="tokens-overview" className="tile">
      <TokenGrantsOverview delegatedTokens={delegatedTokens} />
      <hr />
      <section>
        <h4 className="text-grey-60">Owned Tokens</h4>
        <h2 className="balance">{displayAmount(keepBalance)}</h2>
        <div className="text-samll">
          Staked Owned Tokens: {displayAmount(stakingBalance)}
          <p className="text-smaller text-grey-30">Tokens you own that are delegated to an operator and doing work on the network.</p>
        </div>
        <div className="text-samll">
          Pending Undelegated Tokens: {displayAmount(pendingUndelegationBalance)}
          <p className="text-smaller text-grey-30">
            Stake undelegated from an operator. Estimated {estimatedUndelegationPeriod} until available.
          </p>
        </div>
      </section>
    </section>
  )
}

export default TokensOverview
