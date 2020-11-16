import React from "react"
import Chip from "../../components/Chip"
import TokenAmount from "../../components/TokenAmount"
import { SubmitButton } from "../../components/Button"
import * as Icons from "../../components/Icons"
import { LoadingOverlay } from "../../components/Loadable"
import { DataTableSkeleton } from "../../components/skeletons"
import RandomBeaconRewardsTable from "../../components/RandomBeaconRewardsTable"

const RandomBeaconRewardsPage = () => {
  return (
    <>
      <section className="tile rewards__overview--random-beacon">
        <div className="rewards__overview__balance">
          <h2 className="h2--alt text-grey-70">Random Beacon Rewards</h2>
          <div className="flex row center">
            <TokenAmount amount="0" currencySymbol="KEEP" />
            <Chip
              text="ALLOCATED"
              size="tiny"
              color="success"
              className="ml-1"
            />
          </div>
        </div>
        <div className="rewards__overview__period">
          <h5 className="text-grey-70">current rewards period</h5>
          <span className="rewards-period__date">11/15/2020 - 12/15/2020</span>
          <div className="rewards-period__remaining-periods">
            <Icons.Time width="16" height="16" className="time-icon--grey-30" />
            {/* TODO tooltip */}
            <h4>9 rewards periods remaining</h4>
          </div>
        </div>
        <div className="rewards__overview__withdraw">
          <SubmitButton
            className="btn btn-primary btn-lg"
            onSubmitAction={() => console.log("submit btn")}
          >
            withdraw all
          </SubmitButton>
        </div>
      </section>
      <LoadingOverlay
        isFetching={false}
        skeletonComponent={
          <DataTableSkeleton columns={4} subtitleWidth="30%" />
        }
      >
        <section className="tile rewards__datatable--random-beacon">
          <RandomBeaconRewardsTable
            data={[
              {
                amount: 30,
                status: "ALLOCATED",
                period: "10/15/2020 - 11/15/2020",
                groupPublicKey: "0x086813525A7dC7dafFf015Cdf03896Fd276eab60",
              },
            ]}
          />
        </section>
      </LoadingOverlay>
    </>
  )
}

const Component = () => <></>

RandomBeaconRewardsPage.route = {
  title: "Random Beacon",
  path: "/rewards/random-beacon",
  exact: true,
  withConnectWalletGuard: false,
  // TODO: empty state page component
  emptyStateComponent: Component,
}

export default RandomBeaconRewardsPage
