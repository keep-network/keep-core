import React from "react"
import { Route, Switch, Redirect } from "react-router-dom"
import LiquidationsPage from "./LiquidationsPage"
import LiquidateDepositPage from "./LiquidateDepositPage"

const LiquidationsPageContainer = (props) => {

  return (
    <Switch>
      <Route
        exact
        path="/liquidations"
        component={LiquidationsPage}
      />
      <Route
        exact
        path="/liquidations/:depositAddress"
        component={LiquidateDepositPage}
      />
      <Redirect to="/liquidations" />
    </Switch>
  )
}

export default LiquidationsPageContainer
