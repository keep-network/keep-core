import React from "react"
import { Route, Switch, Redirect } from "react-router-dom"
import TBTCRewardsPage from "./TBTCRewardsPage"
import ReawrdsPage from "./RewardsPage"
import PageWrapper from "../components/PageWrapper"

const subLinks = [
  { title: "Random Beacon", path: "/rewards/random-beacon" },
  { title: "tBtc", path: "/rewards/tbtc" },
]

const RewardsPageContainer = (props) => {
  return (
    // TODO: Rename to `earings` in the followup PR, according to the Figma
    // views.
    <PageWrapper title="Rewards" subLinks={subLinks}>
      <Switch>
        <Route exact path="/rewards/random-beacon" component={ReawrdsPage} />
        <Route exact path="/rewards/tbtc" component={TBTCRewardsPage} />
        <Redirect to="/rewards/random-beacon" />
      </Switch>
    </PageWrapper>
  )
}

export default RewardsPageContainer
