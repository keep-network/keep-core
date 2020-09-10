import React from "react"
import { Route, Switch, Redirect } from "react-router-dom"
// import TBTCRewardsPage from "./TBTCRewardsPage"
import ReawrdsPage from "./RewardsPage"

const RewardsPageContainer = (props) => {
  return (
    <Switch>
      <Route exact path="/rewards/random-beacon" component={ReawrdsPage} />
      {/* <Route exact path="/rewards/tbtc" component={TBTCRewardsPage} /> */}
      <Redirect to="/rewards/random-beacon" />
    </Switch>
  )
}

export default RewardsPageContainer
