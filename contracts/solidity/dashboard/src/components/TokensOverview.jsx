import React from 'react'

const TokensOverview = () => {
  return (
    <section id="tokens-overview" className="tile">
      <h5>Totals</h5>
      <section>
        <span className="text-label text-white" style={{ background: '#AC6E16', padding: '.2rem' }}>OWNED</span>
        <h2 className="balance">200</h2>
        <div className="text-samll">
          Staked Owned Tokens: 200
          <p className="text-smaller text-grey">Tokens you own that are delegated to an operator and doing work on the network.</p>
        </div>
        <div className="text-samll">
          Pengind Undelegated Tokens: 0
          <p className="text-smaller text-grey">Stake undelegated from an operator. Estimated X number of blocks until available.</p>
        </div>
        <hr />
      </section>
      <section>
        <span className="text-label text-white" style={{ background: '#AC6E16', padding:  '.2rem' }}>GRANTED</span>
        <h2 className="balance">200</h2>
        <div className="text-samll">
          Staked Tokens: 0
          <p className="text-smaller text-grey">Tokens you were granted that are delegated to an operator and doing work on the network.</p>
        </div>
      </section>
    </section>
  )
}

export default TokensOverview
