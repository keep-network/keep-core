import { useDispatch } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { isSameEthAddress } from "../../utils/general.utils"
import {
  TOKEN_GRANT_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
} from "../../lib/keep/contracts"
import { KeepExplorerMode } from "../../contracts"

export const useSubscribeToDepositWithdrawEvent = () => {
  const dispatch = useDispatch()
  const {
    [TOKEN_STAKING_ESCROW_CONTRACT_NAME]: { instance: tokenStakingEscrow },
    [TOKEN_GRANT_CONTRACT_NAME]: { instance: grantContract },
    web3,
  } = KeepExplorerMode
  const defaultAccount = web3?.lib?.eth?.defaultAccount

  useSubscribeToExplorerModeContractEvent(
    TOKEN_STAKING_ESCROW_CONTRACT_NAME,
    "DepositWithdrawn",
    async (event) => {
      try {
        const {
          returnValues: { grantee, operator, amount },
        } = event

        // A `grantee` param in the `DepositWithdrawn` event always points to the "right" grantee address.
        // No needed additional check if it's about a managed grant.
        if (!isSameEthAddress(grantee, defaultAccount)) {
          return
        }

        const grantId = await tokenStakingEscrow.methods
          .depositGrantId(operator)
          .call()

        const availableToStake = await grantContract.methods
          .availableToStake(grantId)
          .call()

        dispatch({
          type: "token-grant/grant_withdrawn",
          payload: { grantId, amount, operator, availableToStake },
        })
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode DepositWithdrawn event`,
          error
        )
      }
    }
  )
}
