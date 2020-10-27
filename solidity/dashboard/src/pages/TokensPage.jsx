import React from "react"
import DelegateStakeForm from "../components/DelegateStakeForm"
import TokensOverview from "../components/TokensOverview"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import Tile from "../components/Tile"
import TokensContextSwitcher from "../components/TokensContextSwitcher"
import DelegationOverview from "../components/DelegationOverview"
import { useModal } from "../hooks/useModal"
import { connect } from "react-redux"
import moment from "moment"
import EmptyStateComponent from "./delegation/EmptyStatePage"

const confirmationModalOptions = (initializationPeriod) => ({
  modalOptions: { title: "Initiate Delegation" },
  title: "You’re about to delegate stake.",
  subtitle: `You’re delegating KEEP tokens. You will be able to cancel the delegation for up to ${moment()
    .add(initializationPeriod, "seconds")
    .fromNow(true)}. After that time, you can undelegate your stake.`,
  btnText: "delegate",
  confirmationText: "DELEGATE",
})

const TokensPage = ({ delegateStake, restProps }) => {
  const { openConfirmationModal } = useModal()

  const {
    keepTokenBalance,
    minimumStake,
    selectedGrant,
    tokensContext,
    initializationPeriod,
  } = useTokensPageContext()

  const handleSubmit = async (values, meta) => {
    await openConfirmationModal(confirmationModalOptions(initializationPeriod))
    const grantData =
      tokensContext === "owned"
        ? {}
        : { ...selectedGrant, grantId: selectedGrant.id }
    delegateStake(
      {
        ...values,
        ...grantData,
        amount: values.stakeTokens,
      },
      meta
    )
  }

  const getAvailableToStakeAmount = () => {
    if (tokensContext === "granted") {
      return selectedGrant.availableToStake
    }

    return keepTokenBalance
  }

  return (
    <>
      <TokensContextSwitcher />
      <div className="tokens-wrapper">
        <Tile
          title="Delegate Tokens"
          id="delegate-stake-section"
          withTooltip
          tooltipProps={{
            text: (
              <>
                <span className="text-bold">Delegation</span>&nbsp; sets aside
                an amount of KEEP to be staked by a trusted third party,
                referred to within the dApp as an operator.
              </>
            ),
          }}
        >
          <DelegateStakeForm
            onSubmit={handleSubmit}
            minStake={minimumStake}
            availableToStake={getAvailableToStakeAmount()}
          />
        </Tile>
        <TokensOverview />
      </div>
      <DelegationOverview />
    </>
  )
}

const mapDispatchToProps = (dispatch) => ({
  delegateStake: (values, meta) =>
    dispatch({
      type: "staking/delegate_request",
      payload: values,
      meta,
    }),
})

const ConnectedTokensPage = connect(null, mapDispatchToProps)(TokensPage)

ConnectedTokensPage.route = {
  title: "Delegate",
  path: "/tokens/delegate",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStateComponent,
}

export default ConnectedTokensPage
