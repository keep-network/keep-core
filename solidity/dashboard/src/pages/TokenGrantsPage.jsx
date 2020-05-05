import React from "react"
import {
  TokenGrantDetails,
  TokenGrantStakedDetails,
  TokenGrantUnlockingdDetails,
} from "../components/TokenGrantOverview"
import PageWrapper from "../components/PageWrapper"
import TokenAmount from "../components/TokenAmount"
import { useTokensPageContext } from "../contexts/TokensPageContext"

const TokenGrantsPage = () => {
  const { grants, grantTokenBalance } = useTokensPageContext()

  return (
    <PageWrapper title="Token Grants" className="">
      <TokenAmount
        wrapperClassName="mb-2"
        amount={grantTokenBalance}
        amountClassName="h2 text-grey-40"
        currencyIconProps={{ className: "keep-outline grey-40" }}
        withMetricSuffix={false}
      />
      {grants.map(renderTokenGrantOverview)}
    </PageWrapper>
  )
}

const renderTokenGrantOverview = (tokenGrant) => {
  return (
    <section key={tokenGrant.id} className="tile token-grant-overview">
      <div className="grant-amount">
        <TokenGrantDetails title="Grant Amount" selectedGrant={tokenGrant} />
      </div>
      <div className="unlocking-details">
        <TokenGrantUnlockingdDetails selectedGrant={tokenGrant} />
      </div>
      <div className="staked-details">
        <TokenGrantStakedDetails selectedGrant={tokenGrant} />
      </div>
    </section>
  )
}

export default React.memo(TokenGrantsPage)
