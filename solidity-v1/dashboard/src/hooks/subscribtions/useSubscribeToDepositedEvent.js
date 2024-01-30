import { useDispatch, useSelector } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { MODAL_TYPES } from "../../constants/constants"
import { showModal } from "../../actions/modal"
import {
  TOKEN_GRANT_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
} from "../../lib/keep/contracts"
import { KeepExplorerMode } from "../../contracts"

export const useSubscribeToDepositedEvent = () => {
  const grants = useSelector((state) => state.tokenGrants.grants)
  const dispatch = useDispatch()
  const {
    [TOKEN_GRANT_CONTRACT_NAME]: { instance: grantContract },
  } = KeepExplorerMode

  useSubscribeToExplorerModeContractEvent(
    TOKEN_STAKING_ESCROW_CONTRACT_NAME,
    "Deposited",
    async (event) => {
      try {
        const {
          transactionHash,
          returnValues: { operator, grantId, amount },
        } = event

        if (grants.find((grant) => grant.id === grantId)) {
          dispatch({ type: "staking/remove_delegation", payload: operator })
          dispatch({ type: "staking/remove_undelegation", payload: operator })

          const availableToWitdrawGrant = await grantContract.methods
            .withdrawable(grantId)
            .call()

          dispatch(
            showModal({
              modalType: MODAL_TYPES.StakingTokensClaimed,
              modalProps: { txHash: transactionHash },
            })
          )

          dispatch({
            type: "token-grant/grant_deposited",
            payload: {
              grantId,
              availableToWitdrawGrant,
              amount,
              operator,
            },
          })
        }
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode Deposited event`,
          error
        )
      }
    }
  )
}
