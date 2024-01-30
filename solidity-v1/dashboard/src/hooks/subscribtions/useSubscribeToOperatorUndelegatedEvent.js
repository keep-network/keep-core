import { useDispatch, useSelector } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { MODAL_TYPES } from "../../constants/constants"
import { isSameEthAddress } from "../../utils/general.utils"
import { showModal } from "../../actions/modal"
import moment from "moment"
import { OPERATOR_DELEGATION_UNDELEGATED } from "../../actions"
import {
  TOKEN_GRANT_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
} from "../../lib/keep/contracts"
import { KeepExplorerMode } from "../../contracts"

export const useSubscribeToOperatorUndelegatedEvent = () => {
  const { undelegationPeriod } = useSelector((state) => state.operator)
  const dispatch = useDispatch()
  const {
    [TOKEN_GRANT_CONTRACT_NAME]: { instance: grantContract },
    web3,
  } = KeepExplorerMode
  const defaultAccount = web3?.lib?.eth?.defaultAccount

  useSubscribeToExplorerModeContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "Undelegated",
    async (event) => {
      try {
        const {
          transactionHash,
          returnValues: { operator, undelegatedAt },
        } = event

        if (!isSameEthAddress(defaultAccount, operator)) {
          return
        }

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

        const undelegationCompletedAt = moment
          .unix(undelegatedAt)
          .add(undelegationPeriod, "seconds")

        dispatch({
          type: OPERATOR_DELEGATION_UNDELEGATED,
          payload: { undelegationCompletedAt, delegationStatus: "UNDELEGATED" },
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
