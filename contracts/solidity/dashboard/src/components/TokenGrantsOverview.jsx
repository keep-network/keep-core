import React, { useEffect, useState } from 'react'
import TokenGrantOverview from './TokenGrantOverview'
import Dropdown from './Dropdown'
import SelectedGrantDropdown from './SelectedGrantDropdown'
import { useSubscribeToContractEvent } from '../hooks/useSubscribeToContractEvent'
import { displayAmount, isEmptyObj } from '../utils/general.utils'
import { TOKEN_GRANT_CONTRACT_NAME } from '../constants/constants'
import { findIndexAndObject } from '../utils/array.utils'
import { useTokensPageContext, GRANT_STAKED, GRANT_WITHDRAWN } from '../contexts/TokensPageContext'

const TokenGrantsOverview = (props) => {
  const {
    grants,
    grantTokenBalance,
    dispatch,
    refreshGrantTokenBalance,
    refreshKeepTokenBalance,
  } = useTokensPageContext()
  const [selectedGrant, setSelectedGrant] = useState({})
  const { latestEvent: stakedEvent } = useSubscribeToContractEvent(TOKEN_GRANT_CONTRACT_NAME, 'TokenGrantStaked')
  const { latestEvent: withdrawanEvent } = useSubscribeToContractEvent(TOKEN_GRANT_CONTRACT_NAME, 'TokenGrantWithdrawn')

  useEffect(() => {
    if (isEmptyObj(selectedGrant) && grants.length > 0) {
      setSelectedGrant(grants[0])
    } else if (!isEmptyObj(selectedGrant)) {
      const { obj: updatedGrant } = findIndexAndObject('id', selectedGrant.id, grants)
      setSelectedGrant(updatedGrant)
    }
  }, [grants])

  const onSelect = (selectedItem) => {
    setSelectedGrant(selectedItem)
  }

  useEffect(() => {
    if (isEmptyObj(stakedEvent)) {
      return
    }
    const { returnValues: { grantId, amount } } = stakedEvent
    dispatch({ type: GRANT_STAKED, payload: { grantId, amount } })
  }, [stakedEvent.transactionHash])

  useEffect(() => {
    if (isEmptyObj(withdrawanEvent)) {
      return
    }
    const { returnValues: { grantId, amount } } = withdrawanEvent
    dispatch({ type: GRANT_WITHDRAWN, payload: { grantId, amount } })
    refreshGrantTokenBalance()
    refreshKeepTokenBalance()
  }, [withdrawanEvent.transactionHash])

  return (
    <section>
      <h4 className="text-grey-60">Granted Tokens</h4>
      <h2 className="balance">{displayAmount(grantTokenBalance)}</h2>
      <div style={grants.length === 0 ? { display: 'none' } : {}}>
        {
          grants.length > 1 &&
              <Dropdown
                onSelect={onSelect}
                options={grants}
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
  )
}

export default React.memo(TokenGrantsOverview)
