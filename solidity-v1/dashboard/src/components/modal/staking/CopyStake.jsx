import React from "react"
import StepNav from "../../StepNav"
import {
  CopyStakeStepO,
  CopyStakeStep1,
  CopyStakeStep2,
  CopyStakeStep3,
  CopyStakeStep4,
} from "../../copy-stake-steps"
import {
  INCREMENT_STEP,
  DECREMENT_STEP,
  SET_STRATEGY,
  SET_DELEGATION,
  RESET_COPY_STAKE_FLOW,
} from "../../../actions"
import { connect } from "react-redux"
import * as Icons from "../../Icons"
import { Modal, ModalContent, ModalCloseButton } from "../Modal"

const copyStakeSteps = ["balance", "upgrade", "review", "complete"]

const CopyStakeComponent = ({
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
  resetSteps,
  onClose: closeModal,
  ...restProps
}) => {
  const onClose = () => {
    closeModal()
    resetSteps()
  }

  const onSubmit = () => {
    if (
      selectedStrategy === "WAIT_FLOW" &&
      selectedDelegation.canRecoverStake
    ) {
      restProps.recoverOldStake(selectedDelegation)
    } else if (
      selectedStrategy === "WAIT_FLOW" &&
      !selectedDelegation.isUndelegating
    ) {
      restProps.undelegateOldStake(selectedDelegation)
    } else if (selectedStrategy === "COPY_STAKE_FLOW") {
      restProps.copyStake(selectedDelegation)
    }
  }

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
            incrementStep={onSubmit}
            strategy={selectedStrategy}
            delegation={selectedDelegation || {}}
            isProcessing={restProps.isProcessing}
          />
        )
      case 4:
        return (
          <CopyStakeStep4
            onClose={onClose}
            strategy={selectedStrategy}
            undelegationPeriod={restProps.oldUndelegationPeriod}
            selectedDelegation={selectedDelegation}
          />
        )
    }
  }

  return (
    <Modal size="full" isOpen onClose={onClose}>
      <ModalContent>
        <ModalCloseButton>
          <div className="flex row center">
            <Icons.Cross width={15} height={15} />
            <span className="h5 text-grey-60" style={{ marginLeft: "0.5rem" }}>
              close
            </span>
          </div>
        </ModalCloseButton>
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
      </ModalContent>
    </Modal>
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
    undelegateOldStake: (delegation) =>
      dispatch({ type: "copy-stake/undelegate_request", payload: delegation }),
    recoverOldStake: (delegation) =>
      dispatch({ type: "copy-stake/recover_request", payload: delegation }),
    resetSteps: () => dispatch({ type: RESET_COPY_STAKE_FLOW }),
    copyStake: (delegation) =>
      dispatch({ type: "copy-stake/copy-stake_request", payload: delegation }),
  }
}

export const CopyStake = connect(
  mapStateToProps,
  mapDispatchToProps
)(CopyStakeComponent)
