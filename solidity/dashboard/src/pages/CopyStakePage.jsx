import React from "react"
import StepNav from "../components/StepNav"
import {
  CopyStakeStepO,
  CopyStakeStep1,
  CopyStakeStep2,
  CopyStakeStep3,
  CopyStakeStep4,
} from "../components/copy-stake-steps"
import {
  INCREMENT_STEP,
  DECREMENT_STEP,
  SET_STRATEGY,
  SET_DELEGATION,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST,
} from "../actions"
import { connect } from "react-redux"

const copyStakeSteps = ["balance", "upgrade", "review", "complete"]

const CopyStakePage = ({
  incrementStep,
  decrementStep,
  setStrategy,
  setDelegation,
  selectedDelegation,
  step,
  selectedStrategy,
  fetchOldDelegations,
  oldDelegations,
  oldDelegationsFetching,
}) => {
  const renderContent = () => {
    const defaultProps = { incrementStep, decrementStep }
    switch (step) {
      case 0:
      default:
        return <CopyStakeStepO {...defaultProps} />
      case 1:
        return (
          <CopyStakeStep1
            {...defaultProps}
            delegations={oldDelegations}
            isFetching={oldDelegationsFetching}
            selectedDelegation={selectedDelegation}
            onSelectDelegation={setDelegation}
            fetchDelegations={fetchOldDelegations}
          />
        )
      case 2:
        return (
          <CopyStakeStep2
            {...defaultProps}
            setStrategy={setStrategy}
            selectedStrategy={selectedStrategy}
          />
        )
      case 3:
        return (
          <CopyStakeStep3
            {...defaultProps}
            strategy={selectedStrategy}
            delegation={selectedDelegation || {}}
          />
        )
      case 4:
        return <CopyStakeStep4 />
    }
  }

  return (
    <div className="copy-stake__layout">
      <nav className="copy-stake__nav">
        <div className="copy-stake__nav__indicator">
          <StepNav steps={copyStakeSteps} activeStep={step} />
        </div>
      </nav>
      <main className="copy-stake__content-wrapper">
        <div className="copy-stake__content">{renderContent()}</div>
      </main>
    </div>
  )
}

const mapStateToProps = ({ copyStake }) => {
  return copyStake
}

const mapDispatchToProps = (dispatch) => {
  return {
    incrementStep: () => dispatch({ type: INCREMENT_STEP }),
    decrementStep: () => dispatch({ type: DECREMENT_STEP }),
    setStrategy: (strategy) =>
      dispatch({ type: SET_STRATEGY, payload: strategy }),
    setDelegation: (delegation) =>
      dispatch({ type: SET_DELEGATION, payload: delegation }),
    fetchOldDelegations: () =>
      dispatch({ type: FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST }),
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(CopyStakePage)
