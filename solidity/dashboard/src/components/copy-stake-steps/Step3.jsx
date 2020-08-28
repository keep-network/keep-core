import React from "react"
import TokenAmount from "../TokenAmount"
import Button from "../Button"
import { KeepLoadingIndicator } from "../Icons"

const styles = {
  title: {
    marginBottom: "2.5rem",
  },
  subtitle: {
    textAlign: "justify",
    marginBottom: "2.5rem",
  },
  addressWrapper: {
    marginBottom: "0.5rem",
  },
}

const subtitle = {
  COPY_STAKE_FLOW:
    "This stake balance will be copied to a newly upgraded delegation. Your current stake will start the undelegation process from the old staking contract.",
  WAIT_FLOW:
    "The total balance of the following stake will start the undelegation process from the old staking contract.",
}

const CopyStakeStep3 = ({
  delegation,
  strategy,
  incrementStep,
  decrementStep,
  isProcessing,
}) => {
  const getSubtitle = () => {
    if (strategy === "WAIT_FLOW") {
      if (
        delegation &&
        delegation.isUndelegation &&
        !delegation.canRecoverStake
      ) {
        return `Your delegation has been already undelegated. You can recover stake after
          ${delegation.undelegationCompleteAt.fromNow(true)}`
      } else if (delegation && delegation.canRecoverStake) {
        return `Your delegation has been already undelegated. Hit confirm recover to transfer back.`
      } else {
        return subtitle[strategy]
      }
    } else if (strategy === "COPY_STAKE_FLOW") {
      return subtitle[strategy]
    }
  }

  const getSubmitBtnText = () => {
    if (
      strategy === "WAIT_FLOW" &&
      delegation &&
      (delegation.canRecoverStake || delegation.isUndelegation)
    ) {
      return "confirm recover"
    } else if (
      strategy === "WAIT_FLOW" &&
      delegation &&
      !delegation.isUndelegation
    ) {
      return "confirm undelegation"
    } else {
      return "confirm upgrade"
    }
  }

  const isBtnDisabled = () => {
    if (
      strategy === "WAIT_FLOW" &&
      delegation.isUndelegation &&
      !delegation.canRecoverStake
    ) {
      return true
    }

    return false
  }

  return isProcessing ? (
    <KeepLoadingIndicator width={300} height={300} />
  ) : (
    <>
      <h2 style={styles.title}>Review your stake details below.</h2>
      <h3 className="text-grey-70" style={styles.subtitle}>
        {getSubtitle()}
      </h3>
      <section className="tile" style={{ width: "100%" }}>
        <h3>
          Stake balance to
          {strategy === "COPY_STAKE_FLOW" ? " Copy" : " Undelegate"}
        </h3>
        <TokenAmount
          amount={delegation.amount}
          currencySymbol="KEEP"
          wrapperClassName="mb-1"
        />
        <Address address={delegation.authorizerAddress} label="authorizer" />
        <Address address={delegation.operatorAddress} label="operator" />
        <Address address={delegation.beneficiary} label="beneficiary" />
      </section>
      <div className="flex row space-between self-end">
        <Button
          onClick={decrementStep}
          className="btn btn-transparent btn-lg mr-2"
        >
          back
        </Button>
        <Button
          disabled={isBtnDisabled()}
          onClick={incrementStep}
          className="btn btn-primary btn-lg"
        >
          {getSubmitBtnText()}
        </Button>
      </div>
    </>
  )
}

const Address = ({ label, address }) => (
  <div className="flex row center" style={styles.addressWrapper}>
    <h5 className="text-grey-70 flex-1">{label}</h5>
    <div className="text-big text-grey-50">{address}</div>
  </div>
)

export default CopyStakeStep3
