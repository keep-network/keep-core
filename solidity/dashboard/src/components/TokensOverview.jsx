import React, { useMemo } from "react"
import TokenGrantOverview from "./TokenGrantOverview"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import OwnedTokensOverview from "./OwnedTokensOverview.jsx"
import { add } from "../utils/arithmetics.utils"

const TokensOverview = (props) => {
  const {
    tokensContext,
    selectedGrant,
    keepTokenBalance,
    ownedTokensUndelegationsBalance,
    ownedTokensDelegationsBalance,
    getGrantStakedAmount,
  } = useTokensPageContext()

  const selectedGrantStakedAmount = useMemo(() => {
    return getGrantStakedAmount(selectedGrant.id)
  }, [getGrantStakedAmount, selectedGrant.id])

  return tokensContext === "granted" ? (
    <TokenGrantOverview
      selectedGrant={selectedGrant}
      selectedGrantStakedAmount={selectedGrantStakedAmount}
    />
  ) : (
    <OwnedTokensOverview
      keepBalance={keepTokenBalance}
      stakedBalance={add(
        ownedTokensUndelegationsBalance,
        ownedTokensDelegationsBalance
      )}
    />
  )
}

export default TokensOverview
