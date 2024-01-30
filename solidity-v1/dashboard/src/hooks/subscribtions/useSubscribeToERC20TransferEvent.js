import { keepBalanceActions } from "../../actions"
import { useDispatch } from "react-redux"
import { useWeb3Context } from "../../components/WithWeb3Context"
import { KEEP_TOKEN_CONTRACT_NAME } from "../../lib/keep/contracts"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"

export const useSubscribeToERC20TransferEvent = () => {
  const dispatch = useDispatch()
  const {
    eth: { defaultAccount },
  } = useWeb3Context()

  const fromOptions = {
    filter: {
      from: defaultAccount,
    },
  }

  const toOptions = {
    filter: {
      to: defaultAccount,
    },
  }

  useSubscribeToExplorerModeContractEvent(
    KEEP_TOKEN_CONTRACT_NAME,
    "Transfer",
    (event) => {
      console.log("event: ", event)
      dispatch(keepBalanceActions.keepTokenTransferFromEventEmitted(event))
    },
    fromOptions
  )

  useSubscribeToExplorerModeContractEvent(
    KEEP_TOKEN_CONTRACT_NAME,
    "Transfer",
    (event) => {
      console.log("event: ", event)
      dispatch(keepBalanceActions.keepTokenTransferToEventEmitted(event))
    },
    toOptions
  )
}
