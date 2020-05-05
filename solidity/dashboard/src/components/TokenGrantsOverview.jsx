import React, { useEffect, useState } from "react"
import TokenGrantOverview from "./TokenGrantOverview"
import Dropdown from "./Dropdown"
import SelectedGrantDropdown from "./SelectedGrantDropdown"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent"
import { isEmptyObj } from "../utils/general.utils"
import { displayAmount } from "../utils/token.utils"
import { TOKEN_GRANT_CONTRACT_NAME } from "../constants/constants"
import { findIndexAndObject } from "../utils/array.utils"
import { useTokensPageContext } from "../contexts/TokensPageContext"

const TokenGrantsOverview = (props) => {
  const {
    grants,
    grantTokenBalance,
    refreshGrantTokenBalance,
    refreshKeepTokenBalance,
    grantStaked,
    grantWithdrawn,
  } = useTokensPageContext()
  const [selectedGrant, setSelectedGrant] = useState({})

  const subscribeToStakedEventCallback = (stakedEvent) => {
    const {
      returnValues: { grantId, amount },
    } = stakedEvent
    grantStaked(grantId, amount)
  }

  const subscribeToWithdrawanEventCallback = (withdrawanEvent) => {
    const {
      returnValues: { grantId, amount },
    } = withdrawanEvent
    grantWithdrawn(grantId, amount)
    refreshGrantTokenBalance()
    refreshKeepTokenBalance()
  }

  useSubscribeToContractEvent(
    TOKEN_GRANT_CONTRACT_NAME,
    "TokenGrantStaked",
    subscribeToStakedEventCallback
  )
  useSubscribeToContractEvent(
    TOKEN_GRANT_CONTRACT_NAME,
    "TokenGrantWithdrawn",
    subscribeToWithdrawanEventCallback
  )

  useEffect(() => {
    if (isEmptyObj(selectedGrant) && grants.length > 0) {
      setSelectedGrant(grants[0])
    } else if (!isEmptyObj(selectedGrant)) {
      const { obj: updatedGrant } = findIndexAndObject(
        "id",
        selectedGrant.id,
        grants
      )
      setSelectedGrant(updatedGrant)
    }
  }, [grants, selectedGrant])

  const onSelect = (selectedItem) => {
    setSelectedGrant(selectedItem)
  }

  return (
    <section>
      <h4 className="text-grey-60">Granted Tokens</h4>
      <h2 className="balance">{displayAmount(grantTokenBalance)}</h2>
      <div style={grants.length === 0 ? { display: "none" } : {}}>
        {grants.length > 1 && (
          <Dropdown
            onSelect={onSelect}
            options={grants}
            valuePropertyName="id"
            labelPropertyName="id"
            selectedItem={selectedGrant}
            labelPrefix="Grant ID"
            noItemSelectedText="Select Grant"
            label="Choose Grant"
            selectedItemComponent={
              <SelectedGrantDropdown grant={selectedGrant} />
            }
          />
        )}
        <TokenGrantOverview selectedGrant={selectedGrant} />
      </div>
    </section>
  )
}

export default React.memo(TokenGrantsOverview)
