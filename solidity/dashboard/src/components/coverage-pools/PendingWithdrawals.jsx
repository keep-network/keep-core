import React from "react"
import { claimTokensFromWithdrawal } from "../../actions/coverage-pool"
import PendingWithdrawalsView from "./PendingWithdrawalsView"
import { useDispatch } from "react-redux"

const PendingWithdrawals = () => {
  const dispatch = useDispatch()
  const onClaimTokensSubmitButtonClick = async (awaitingPromise) => {
    dispatch(claimTokensFromWithdrawal(awaitingPromise))
  }

  const onReinitiateWithdrawal = async (awaitingPromise) => {
    console.log("reinitiate withdrawal")
  }

  return (
    <PendingWithdrawalsView
      onClaimTokensSubmitButtonClick={onClaimTokensSubmitButtonClick}
      onReinitiateWithdrawal={onReinitiateWithdrawal}
    />
  )
}

export default PendingWithdrawals
