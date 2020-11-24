import React, { useMemo } from "react"
import TokenGrantOverview from "./TokenGrantOverview"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import OwnedTokensOverview from "./OwnedTokensOverview.jsx"
import { add } from "../utils/arithmetics.utils"
import { LoadingOverlay } from "./Loadable"
import TokenOverviewSkeleton from "./skeletons/TokenOverviewSkeleton"
import { useSelector } from "react-redux"

const TokensOverview = (props) => {
  const {
    tokensContext,
    selectedGrant,
    getGrantStakedAmount,
  } = useTokensPageContext()

  const {
    ownedTokensUndelegationsBalance,
    ownedTokensDelegationsBalance,
    isDelegationDataFetching,
  } = useSelector((state) => state.staking)

  const { isFetching: areGrantsFetching } = useSelector(
    (state) => state.tokenGrants
  )
  const { value: keepTokenBalance } = useSelector(
    (state) => state.keepTokenBalance
  )

  const selectedGrantStakedAmount = useMemo(() => {
    return getGrantStakedAmount(selectedGrant.id)
  }, [getGrantStakedAmount, selectedGrant.id])

  return (
    <section id="tokens-overview" className="tile">
      <LoadingOverlay
        isFetching={
          tokensContext === "granted"
            ? areGrantsFetching
            : isDelegationDataFetching
        }
        skeletonComponent={
          <TokenOverviewSkeleton items={tokensContext === "granted" ? 2 : 1} />
        }
      >
        {tokensContext === "granted" ? (
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
        )}
      </LoadingOverlay>
    </section>
  )
}

export default TokensOverview
