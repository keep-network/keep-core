import React from "react"
import moment from "moment"
import EmptyStateComponent from "./EmptyStatePage"
import DelegateStakeForm from "../../components/DelegateStakeForm"
import DelegatedTokensTable from "../../components/DelegatedTokensTable"
import Undelegations from "../../components/Undelegations"
import { DataTableSkeleton } from "../../components/skeletons"
import { LoadingOverlay } from "../../components/Loadable"
import Tag from "../../components/Tag"
import * as Icons from "../../components/Icons"
import {
  TokenGrantDetails,
  TokenGrantStakedDetails,
  TokenGrantUnlockingdDetails,
  TokenGrantWithdrawnTokensDetails,
} from "../../components/TokenGrantOverview"

const GrantedTokensPage = () => {
  // TODO: get grant from redux or react context
  const selectedGrant = {
    amount: 100,
    unlocked: 20,
    released: 30,
    id: 1,
    start: moment().unix(),
    cliff: 1,
    duration: 1200000,
  }

  return (
    <>
      <section className="granted-page__overview-layout">
        <section className="tile granted-page__overview__grant-details">
          <h4 className="mb-1">Grant Allocation</h4>
          <TokenGrantDetails selectedGrant={selectedGrant} />
        </section>
        <section className="tile granted-page__overview__staked-tokens">
          <h4 className="mb-2">Tokens Staked</h4>
          <TokenGrantStakedDetails
            selectedGrant={selectedGrant}
            stakedAmount={80}
          />
        </section>
        <section className="tile granted-page__overview__stake-form">
          <DelegateStakeForm />
        </section>
        <section className="tile granted-page__overview__withdraw-tokens">
          <h4 className="mb-2">Withdraw Unlocked Tokens</h4>
          <TokenGrantWithdrawnTokensDetails
            selectedGrant={selectedGrant}
            onWithdrawnBtn={() => console.log("on clikc withdraw tokens")}
          />
        </section>
        <section className="tile granted-page__overview__unlocked-tokens">
          <h4 className="mb-2">Tokens Unlocking Progress</h4>
          <TokenGrantUnlockingdDetails selectedGrant={selectedGrant} />
        </section>
      </section>
      <section>
        <header className="flex row center mb-2">
          <h2 className="h2--alt text-grey-60">Grant Activity</h2>
          <div className="flex row center ml-a">
            <Tag IconComponent={Icons.Grant} text="Grant ID" />
            <span className="ml-1 mr-2">1234</span>
            <Tag IconComponent={Icons.Time} text="Issued" />
            <span className="ml-1">01/01/2020</span>
          </div>
        </header>
        <LoadingOverlay
          isFetching={true}
          skeletonComponent={<DataTableSkeleton />}
        >
          <DelegatedTokensTable
            delegatedTokens={[]}
            //   cancelStakeSuccessCallback={cancelStakeSuccessCallback}
            // keepTokenBalance={keepToken.value}
            // undelegationPeriod={undelegationPeriod}
          />
        </LoadingOverlay>
        <LoadingOverlay
          isFetching={true}
          skeletonComponent={<DataTableSkeleton />}
        >
          <Undelegations undelegations={[]} />
        </LoadingOverlay>
      </section>
    </>
  )
}

GrantedTokensPage.route = {
  title: "Granted Tokens",
  path: "/delegation/grant",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStateComponent,
}

export { GrantedTokensPage }
