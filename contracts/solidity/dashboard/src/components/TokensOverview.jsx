import React, { useState, useEffect } from 'react'
import { displayAmount } from '../utils'
import TokenGrantOverview from './TokenGrantOverview'
import Dropdown from './Dropdown'
import SelectedGrantDropdown from './SelectedGrantDropdown'
import { useFetchData } from '../hooks/useFetchData'
import { tokenGrantsService } from '../services/token-grants.service'

const TokensOverview = ({
  undelegationPeriod,
  keepBalance,
  stakingBalance,
  pendingUndelegationBalance,
  grantBalance,
  tokenGrantsStakeBalance,
}) => {
  const [state] = useFetchData(tokenGrantsService.fetchGrants, [])
  const { data } = state
  const [selectedGrant, setSelectedGrant] = useState({})

  useEffect(() => {
    if (selectedGrant && data.length > 0) {
      setSelectedGrant(data[0])
    }
  }, [data])

  const onSelect = (selectedItem) => {
    setSelectedGrant(selectedItem)
  }

  return (
    <section id="tokens-overview" className="tile">
      <section>
        <h4 className="text-grey-60">Granted Tokens</h4>
        <h2 className="balance">{displayAmount(grantBalance)}</h2>
        <div style={data.length === 0 ? { display: 'none' } : {}}>
          {
            data.length > 1 &&
              <Dropdown
                onSelect={onSelect}
                options={data}
                valuePropertyName='id'
                labelPropertyName='id'
                selectedItem={selectedGrant}
                labelPrefix='Grant ID'
                noItemSelectedText='Select Grant'
                label="Choose Grant"
                selectedItemComponent={<SelectedGrantDropdown grant={selectedGrant} />}
              />
          }
          <TokenGrantOverview selectedGrant={selectedGrant} />
        </div>
      </section>
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
          <p className="text-smaller text-grey-30">Stake undelegated from an operator. Estimated {undelegationPeriod} number of blocks until available.</p>
        </div>
      </section>
    </section>
  )
}

export default TokensOverview
