import React, { useMemo } from "react"
import {
  TokenGrantDetails,
  TokenGrantStakedDetails,
  TokenGrantUnlockingdDetails,
} from "../components/TokenGrantOverview"
import TokenAmount from "../components/TokenAmount"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import { LoadingOverlay } from "../components/Loadable"
import { TokenGrantSkeletonOverview } from "../components/skeletons/TokenOverviewSkeleton"
import { add } from "../utils/arithmetics.utils"

const TokenGrantsPage = () => {
  const { grants, grantsAreFetching } = useTokensPageContext()
  const totalGrantAmount = useMemo(() => {
    return grants.map(({ amount }) => amount).reduce(add, "0")
  }, [grants])

  return (
    <>
      <TokenAmount
        wrapperClassName="mb-2"
        amount={totalGrantAmount}
        amountClassName="h2 text-grey-40"
        currencyIconProps={{ className: "keep-outline grey-40" }}
        displayWithMetricSuffix={false}
      />

      <LoadingOverlay
        isFetching={grantsAreFetching}
        skeletonComponent={<TokenGrantSkeletonOverview />}
      >
        {grants.map(renderTokenGrantOverview)}
      </LoadingOverlay>
    </>
  )
}

const renderTokenGrantOverview = (tokenGrant) => (
  <TokenGrantOverview key={tokenGrant.id} tokenGrant={tokenGrant} />
)

const TokenGrantOverview = React.memo(({ tokenGrant }) => {
  const { getGrantStakedAmount } = useTokensPageContext()

  const selectedGrantStakedAmount = useMemo(() => {
    return getGrantStakedAmount(tokenGrant.id)
  }, [getGrantStakedAmount, tokenGrant.id])

  return (
    <section
      key={tokenGrant.id}
      className="tile token-grant-overview"
      style={{ marginBottom: "1.2rem" }}
    >
      <div className="grant-amount">
        <TokenGrantDetails title="Grant Amount" selectedGrant={tokenGrant} />
      </div>
      <div className="unlocking-details">
        <TokenGrantUnlockingdDetails selectedGrant={tokenGrant} />
      </div>
      <div className="staked-details">
        <TokenGrantStakedDetails
          selectedGrant={tokenGrant}
          stakedAmount={selectedGrantStakedAmount}
        />
      </div>
    </section>
  )
})

export default React.memo(TokenGrantsPage)
