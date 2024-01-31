import { useDispatch } from "react-redux"
import { MODAL_TYPES } from "../../constants/constants"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { showModal } from "../../actions/modal"
import { TOKEN_GRANT_CONTRACT_NAME } from "../../lib/keep/contracts"
import { KeepExplorerMode } from "../../contracts"

export const useSubscribeToTokenGrantWithdrawnEvent = () => {
  const dispatch = useDispatch()
  const {
    [TOKEN_GRANT_CONTRACT_NAME]: { instance: grantContract },
  } = KeepExplorerMode

  useSubscribeToExplorerModeContractEvent(
    TOKEN_GRANT_CONTRACT_NAME,
    "TokenGrantWithdrawn",
    async (event) => {
      try {
        const {
          transactionHash,
          returnValues: { grantId, amount },
        } = event

        const availableToStake = await grantContract.methods
          .availableToStake(grantId)
          .call()
        dispatch({
          type: "token-grant/grant_withdrawn",
          payload: { grantId, amount, availableToStake },
        })

        dispatch(
          showModal({
            modalType: MODAL_TYPES.GrantTokensWithdrawn,
            modalProps: { txHash: transactionHash },
          })
        )
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode TokenGrantWithdrawn event`,
          error
        )
      }
    }
  )
}
