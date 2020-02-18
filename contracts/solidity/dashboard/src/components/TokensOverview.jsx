import React from 'react'
import { displayAmount } from '../utils'
import { useModal } from '../hooks/useModal'
import TokenGrantVestingSchedule from './TokenGrantVestingSchedule'
import TokenGrantOverview from './TokenGrantOverview'

const TokensOverview = ({
  undelegationPeriod,
  keepBalance,
  stakingBalance,
  pendingUndelegationBalance,
  grantBalance,
  tokenGrantsStakeBalance,
}) => {
  const { showModal, ModalComponent } = useModal()

  return (
    <section id="tokens-overview" className="tile">
      <ModalComponent title="Vesting Schedule Summary">
        <TokenGrantVestingSchedule />
      </ModalComponent>
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
        <h3 className="text-darker-grey">Granted Tokens</h3>
        <h2 className="balance">{displayAmount(grantBalance)}</h2>
        <TokenGrantOverview />
      </section>
    </section>
  )
}

export default TokensOverview
