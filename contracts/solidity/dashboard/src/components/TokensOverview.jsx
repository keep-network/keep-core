import React from 'react'

const TokensOverview = () => {
  return (
    <section className="tile">
      <h5>Totals</h5>
      <section>
        <h5>Owned Tokens: 200</h5>
        <span className="text-small">Tokens that you own are in your wallet.</span>
        <div style={{ borderLeft: '1px solid black', paddingLeft: '1rem' }}>
          <p className="text-samll">Staked Owned Tokens: 200</p>
          <span className="text-small text-darker-grey">Tokens you own that are delegated to an operator and doing work on the network.</span>
          <p className="text-samll">Pengind Undelegated Tokens: 0</p>
          <span className="text-small text-darker-grey">Stake undelegated from an operator. Estimated X number of blocks until available.</span>
        </div>
      </section>
      <section>
        <h5>Granted Tokens: 3000</h5>
        <span className="text-small">Tokens that are locked in a grant with a vesting schedule.</span>
        <div style={{ borderLeft: '1px solid black', paddingLeft: '1rem' }}>
          <p className="text-samll">Staked Granted Tokens: 0</p>
          <span className="text-small text-darker-grey">Tokens you own that are delegated to an operator and doing work on the network.</span>
        </div>
      </section>
    </section>
  )
}

export default TokensOverview
