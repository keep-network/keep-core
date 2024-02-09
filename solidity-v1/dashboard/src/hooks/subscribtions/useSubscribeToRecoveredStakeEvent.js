import { useDispatch, useSelector } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { isSameEthAddress } from "../../utils/general.utils"
import { MODAL_TYPES } from "../../constants/constants"
import { sub } from "../../utils/arithmetics.utils"
import { showModal } from "../../actions/modal"
import { TOKEN_STAKING_CONTRACT_NAME } from "../../lib/keep/contracts"

export const useSubscribeToRecoveredStakeEvent = () => {
  const { undelegations } = useSelector((state) => state.staking)
  const dispatch = useDispatch()

  useSubscribeToExplorerModeContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "RecoveredStake",
    async (event) => {
      try {
        const {
          transactionHash,
          returnValues: { operator },
        } = event

        const recoveredUndelegation = undelegations.find((undelegation) =>
          isSameEthAddress(undelegation.operatorAddress, operator)
        )

        if (!recoveredUndelegation) {
          return
        }

        dispatch(
          showModal({
            modalType: MODAL_TYPES.StakingTokensClaimed,
            modalProps: { txHash: transactionHash },
          })
        )

        if (!recoveredUndelegation.isFromGrant) {
          dispatch({ type: "staking/remove_undelegation", payload: operator })

          dispatch({
            type: "staking/update_owned_undelegations_tokens_balance",
            payload: { operation: sub, value: recoveredUndelegation.amount },
          })
        }
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode RecoveredStake event`,
          error
        )
      }
    }
  )
}
