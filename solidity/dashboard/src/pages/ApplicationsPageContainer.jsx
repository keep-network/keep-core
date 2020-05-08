import React from "react"
import { Route, Switch, Redirect } from "react-router-dom"
import ApplicationOverviewPage from "./ApplicationOverviewPage"
import TBTCApplicationPage from "./TBTCApplicationPage"
import KeepRandomBeaconApplicationPage from "./KeepRandomBeaconApplicationPage"

const ApplicationsPageContainer = (props) => {
  return (
    <Switch>
      <Route
        exact
        path="/applications/overview"
        component={ApplicationOverviewPage}
      />
      <Route exact path="/applications/tbtc" component={TBTCApplicationPage} />
      <Route
        exact
        path="/applications/random-beacon"
        component={KeepRandomBeaconApplicationPage}
      />
      <Redirect to="/applications/overview" />
    </Switch>
  )
}

export default ApplicationsPageContainer
