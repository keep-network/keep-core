import React from "react"
import {
  claimTokensFromWithdrawal,
  withdrawAssetPool,
} from "../../actions/coverage-pool"
import PendingWithdrawalsView from "./PendingWithdrawalsView"
import { useDispatch } from "react-redux"

const PendingWithdrawals = () => {
  const dispatch = useDispatch()
  const onClaimTokensSubmitButtonClick = async (awaitingPromise) => {
    dispatch(claimTokensFromWithdrawal(awaitingPromise))
  }

  const onReinitiateWithdrawal = async (awaitingPromise) => {
    dispatch(withdrawAssetPool("0", awaitingPromise))
  }

  return (
    <PendingWithdrawalsView
      onClaimTokensSubmitButtonClick={onClaimTokensSubmitButtonClick}
      onReinitiateWithdrawal={onReinitiateWithdrawal}
    />
  )
}

export default PendingWithdrawals
