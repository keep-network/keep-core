import React from "react"
import {
  authorizeSortitionPoolContract,
  authorizeOperatorContract,
  deauthorizeSortitionPoolContract,
} from "../../actions/web3"
import { connect } from "react-redux"
import EmptyStatePage from "./EmptyStatePage"

const ThresholdApplicationPage = ({
  authorizeSortitionPoolContract,
  authorizeOperatorContract,
  deauthorizeSortitionPoolContract,
}) => {
  return <div>test</div>
}

const mapDispatchToProps = {
  authorizeSortitionPoolContract,
  authorizeOperatorContract,
  deauthorizeSortitionPoolContract,
}

const ConnectedThresholdApplicationPage = connect(
  null,
  mapDispatchToProps
)(ThresholdApplicationPage)

ConnectedThresholdApplicationPage.route = {
  title: "Threshold",
  path: "/applications/threshold",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

export default ConnectedThresholdApplicationPage
