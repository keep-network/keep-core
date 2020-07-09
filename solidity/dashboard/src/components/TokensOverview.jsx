import React, { useMemo } from "react"
import TokenGrantOverview from "./TokenGrantOverview"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import OwnedTokensOverview from "./OwnedTokensOverview.jsx"
import { add } from "../utils/arithmetics.utils"
import { LoadingOverlay } from "./Loadable"
import TokenOverviewSkeleton from "./skeletons/TokenOverviewSkeleton"

const TokensOverview = (props) => {
  const {
    tokensContext,
    selectedGrant,
    keepTokenBalance,
    ownedTokensUndelegationsBalance,
    ownedTokensDelegationsBalance,
    getGrantStakedAmount,
    isFetching,
    grantsAreFetching,
  } = useTokensPageContext()

  const selectedGrantStakedAmount = useMemo(() => {
    return getGrantStakedAmount(selectedGrant.id)
  }, [getGrantStakedAmount, selectedGrant.id])

  return (
    <section id="tokens-overview" className="tile">
      <LoadingOverlay
        isFetching={
          tokensContext === "granted" ? grantsAreFetching : isFetching
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
