import { useCallback } from "react"
import { useDispatch } from "react-redux"
import { useModal } from "./useModal"
import { ContractsLoaded } from "../contracts"
import { releaseTokens } from "../actions/web3"
import { isEmptyArray } from "../utils/array.utils"
import { MODAL_TYPES } from "../constants/constants"

const useReleaseTokens = () => {
  const dispatch = useDispatch()
  const { openConfirmationModal } = useModal()

  const onWithdrawTokens = useCallback(
    async (grant, awaitingPromise) => {
      const { escrowOperatorsToWithdraw } = grant

      if (!isEmptyArray(escrowOperatorsToWithdraw)) {
        const { tokenStakingEscrow } = await ContractsLoaded
        await openConfirmationModal(MODAL_TYPES.ConfirmReleaseTokensFromGrant, {
          escrowContractAddress: tokenStakingEscrow.options.address,
        })
      }
      dispatch(releaseTokens(grant, awaitingPromise))
    },
    [dispatch, openConfirmationModal]
  )

  return onWithdrawTokens
}

export default useReleaseTokens
