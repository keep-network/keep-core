import React, { useEffect, useState } from 'react'
import TokenGrantOverview from './TokenGrantOverview'
import Dropdown from './Dropdown'
import SelectedGrantDropdown from './SelectedGrantDropdown'
import { useSubscribeToContractEvent } from '../hooks/useSubscribeToContractEvent'
import { displayAmount, isEmptyObj } from '../utils/general.utils'
import { TOKEN_GRANT_CONTRACT_NAME } from '../constants/constants'
import { findIndexAndObject } from '../utils/array.utils'
import { useTokensPageContext } from '../contexts/TokensPageContext'

const TokenGrantsOverview = ({delegatedTokens}) => {
  const {
    grants,
    grantTokenBalance,
    refreshGrantTokenBalance,
    refreshKeepTokenBalance,
    grantStaked,
    grantWithdrawn,
  } = useTokensPageContext()
  const [selectedGrant, setSelectedGrant] = useState({})
  const [currentTokenAmount, setCurrentTokenAmount] = useState(0)

  const { latestEvent: stakedEvent } = useSubscribeToContractEvent(TOKEN_GRANT_CONTRACT_NAME, 'TokenGrantStaked')
  const { latestEvent: withdrawanEvent } = useSubscribeToContractEvent(TOKEN_GRANT_CONTRACT_NAME, 'TokenGrantWithdrawn')

  useEffect(() => {
    if (isEmptyObj(selectedGrant) && grants.length > 0) {
      onSelect(grants[0])
    } else if (!isEmptyObj(selectedGrant)) {
      const { obj: updatedGrant } = findIndexAndObject('id', selectedGrant.id, grants)
      setSelectedGrant(updatedGrant)
    }
  }, [grants])

  const onSelect = (selectedItem) => {
    setSelectedGrant(selectedItem)

    const result = delegatedTokens.filter((grant) => grant.grantId === selectedItem.id)
    var totalBalance = result.reduce((total, grant) => {
        return total + grant.amount
    })

    setCurrentTokenAmount(totalBalance.amount)
  }

  useEffect(() => {
    if (isEmptyObj(stakedEvent)) {
      return
    }
    const { returnValues: { grantId, amount } } = stakedEvent
    grantStaked(grantId, amount)
  }, [stakedEvent.transactionHash])

  useEffect(() => {
    if (isEmptyObj(withdrawanEvent)) {
      return
    }
    const { returnValues: { grantId, amount } } = withdrawanEvent
    grantWithdrawn(grantId, amount)
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
        <TokenGrantOverview selectedGrant={selectedGrant} delegatedTokens={currentTokenAmount}/>
      </div>
    </section>
  )
}

export default React.memo(TokenGrantsOverview)
