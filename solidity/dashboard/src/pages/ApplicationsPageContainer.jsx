import React from "react"
import { Route, Switch, Redirect } from "react-router-dom"
import TBTCApplicationPage from "./TBTCApplicationPage"
import KeepRandomBeaconApplicationPage from "./KeepRandomBeaconApplicationPage"
import PageWrapper from "../components/PageWrapper"

const subLinks = [
  { title: "Random Beacon", path: "/applications/random-beacon" },
  { title: "tBTC", path: "/applications/tbtc" },
]
const ApplicationsPageContainer = (props) => {
  return (
    <PageWrapper title="Applications" subLinks={subLinks}>
      <Switch>
        <Route
          exact
          path="/applications/tbtc"
          component={TBTCApplicationPage}
        />
        <Route
          exact
          path="/applications/random-beacon"
          component={KeepRandomBeaconApplicationPage}
        />
        <Redirect to="/applications/random-beacon" />
      </Switch>
    </PageWrapper>
  )
}

export default ApplicationsPageContainer
