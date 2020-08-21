import React, { useCallback, useReducer, useMemo, useContext } from "react"
import StepNav from "../components/StepNav"
import {
  CopyStakeStepO,
  CopyStakeStep1,
  CopyStakeStep2,
  CopyStakeStep3,
  CopyStakeStep4,
} from "../components/copy-stake-steps"
import copyStakeReducer, {
  copyStakeInitialData,
  INCREMENT_STEP,
  DECREMENT_STEP,
  SET_STRATEGY,
  SET_DELEGATION,
} from "../reducers/copy-stake"

const copyStakeSteps = ["balance", "upgrade", "review", "complete"]

const CopyStakePage = () => {
  const { fetchDelegations, state, dispatch } = useCopyStakeContext()

  const incrementStep = useCallback(() => {
    dispatch({ type: INCREMENT_STEP })
  }, [dispatch])

  const decrementStep = useCallback(() => {
    dispatch({ type: DECREMENT_STEP })
  }, [dispatch])

  const setStrategy = useCallback(
    (strategy) => {
      dispatch({ type: SET_STRATEGY, payload: strategy })
    },
    [dispatch]
  )

  const setDelegation = useCallback(
    (delegation) => {
      dispatch({ type: SET_DELEGATION, payload: delegation })
    },
    [dispatch]
  )

  const renderContent = () => {
    const defaultProps = { incrementStep, decrementStep }
    console.log("statee", state.step)
    switch (state.step) {
      case 0:
      default:
        return <CopyStakeStepO {...defaultProps} />
      case 1:
        return (
          <CopyStakeStep1
            {...defaultProps}
            delegations={[
              {
                operatorAddress: "0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca",
                authorizerAddress: "0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca",
                beneficiary: "0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca",
                amount: "100000000000000000000000",
              },
              {
                operatorAddress: "0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca",
                authorizerAddress: "0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca",
                beneficiary: "0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca",
                amount: "100000000000000000000000",
              },
            ]}
            selectedDelegation={state.selectedDelegation}
            onSelectDelegation={setDelegation}
          />
        )
      case 2:
        return (
          <CopyStakeStep2
            {...defaultProps}
            setStrategy={setStrategy}
            selectedStrategy={state.selectedStrategy}
          />
        )
      case 3:
        return (
          <CopyStakeStep3
            {...defaultProps}
            strategy={state.selectedStrategy}
            delegation={state.selectedDelegation || {}}
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
          <StepNav steps={copyStakeSteps} activeStep={state.step} />
        </div>
      </nav>
      <main className="copy-stake__content-wrapper">
        <div className="copy-stake__content">{renderContent()}</div>
      </main>
    </div>
  )
}

const CopyStakeProvider = (props) => {
  const [state, dispatch] = useReducer(copyStakeReducer, copyStakeInitialData)

  const contextValue = useMemo(() => {
    return { state, dispatch }
  }, [state, dispatch])

  const fetchOldDelegations = () => {}

  return (
    <CopyStakeContext.Provider value={{ ...contextValue, fetchOldDelegations }}>
      <CopyStakePage />
    </CopyStakeContext.Provider>
  )
}

const CopyStakeContext = React.createContext({
  dispatch: () => {},
  ...copyStakeInitialData,
})

const useCopyStakeContext = () => {
  return useContext(CopyStakeContext)
}

export default CopyStakeProvider
