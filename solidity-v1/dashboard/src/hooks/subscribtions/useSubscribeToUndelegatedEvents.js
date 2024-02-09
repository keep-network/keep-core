import { useDispatch, useSelector } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { isSameEthAddress } from "../../utils/general.utils"
import { showModal } from "../../actions/modal"
import { MODAL_TYPES } from "../../constants/constants"
import moment from "moment"
import { REMOVE_STAKE_FROM_THRESHOLD_AUTH_DATA } from "../../actions"
import { add, sub } from "../../utils/arithmetics.utils"
import { TOKEN_STAKING_CONTRACT_NAME } from "../../lib/keep/contracts"

export const useSubscribeToUndelegatedEvents = () => {
  const { delegations, undelegationPeriod } = useSelector(
    (state) => state.staking
  )
  const dispatch = useDispatch()

  useSubscribeToExplorerModeContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "Undelegated",
    async (event) => {
      try {
        const {
          transactionHash,
          returnValues: { operator, undelegatedAt },
        } = event

        const delegation = delegations.find(({ operatorAddress }) =>
          isSameEthAddress(operatorAddress, operator)
        )

        if (!delegation) {
          return
        }

        // If the delegation exists, we create a undelegation based on the existing delegation.
        dispatch(
          showModal({
            modalType: MODAL_TYPES.UndelegationInitiated,
            modalProps: {
              txHash: transactionHash,
              undelegatedAt,
              undelegationPeriod,
            },
          })
        )

        const undelegation = {
          ...delegation,
          undelegatedAt: moment.unix(undelegatedAt),
          undelegationCompleteAt: moment
            .unix(undelegatedAt)
            .add(undelegationPeriod, "seconds"),
          canRecoverStake: false,
        }

        if (!undelegation.isFromGrant) {
          dispatch({
            type: "staking/update_owned_delegated_tokens_balance",
            payload: { operation: sub, value: undelegation.amount },
          })
          dispatch({
            type: "staking/update_owned_undelegations_tokens_balance",
            payload: { operation: add, value: undelegation.amount },
          })
        }

        dispatch({ type: "staking/remove_delegation", payload: operator })
        dispatch({ type: "staking/add_undelegation", payload: undelegation })
        dispatch({
          type: REMOVE_STAKE_FROM_THRESHOLD_AUTH_DATA,
          payload: operator,
        })
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode Undelegated event`,
          error
        )
      }
    }
  )
}
