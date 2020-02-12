import React, { useState } from 'react'
import { tokenGrantsService } from '../services/token-grants.service'
import { useFetchData } from '../hooks/useFetchData'
import { formatDate, displayAmount } from '../utils'
import Dropdown from './Dropdown'

const TokenGrantOverview = (props) => {
  const [state] = useFetchData(tokenGrantsService.fetchGrants, [])
  const { data } = state
  const [selectedGrant, setSelectedGrant] = useState(data[0] || {})

  const onSelect = (selectedItem) => {
    setSelectedGrant(selectedItem)
  }

  return (
    <div className="token-grant-overview">
      <div className="text-big">
        Grant ID {selectedGrant.id}
      </div>
      <div className="flex flex-row-center flex-row-space-between">
        <h4 className="balance">{displayAmount(selectedGrant.amount)}&nbsp;KEEP</h4>
        <a href="" className="text-warning">Vesting schedule</a>
      </div>
      <div className="text-small text-grey">
        Issued on {formatDate(selectedGrant.start * 1000)}
      </div>
      <div>
        <Dropdown
          onSelect={onSelect}
          options={data}
          valuePropertyName='id'
          labelPropertyName='id'
          selectedItem={selectedGrant}
          labelPrefix='Grant ID'
        />
        <div className="flex flex-row-center">
          <div className="dot grey"/>
          {displayAmount(selectedGrant.vested)}&nbsp;KEEP&nbsp;<span className="text-small text-grey">Vested</span>
        </div>
        <div className="flex flex-row-center">
          <div className="dot brown"/>{displayAmount(selectedGrant.released)}&nbsp;KEEP&nbsp;<span className="text-small text-grey">Released</span>
        </div>
        <div className="flex flex-row-center">
          <div className="dot black"/>{displayAmount(selectedGrant.staked)}&nbsp;KEEP&nbsp;<span className="text-small text-grey">Staked</span>
        </div>
      </div>
    </div>
  )
}

export default TokenGrantOverview
