import { useDispatch } from "react-redux"
import { riskManagerAuctionClosedEventEmitted } from "../../actions/coverage-pool"
import { RISK_MANAGER_V1_CONTRACT_NAME } from "../../lib/keep/contracts"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"

export const useSubscribeToCovPoolsAuctionClosedEvent = () => {
  const dispatch = useDispatch()

  useSubscribeToExplorerModeContractEvent(
    RISK_MANAGER_V1_CONTRACT_NAME,
    "AuctionClosed",
    (event) => {
      dispatch(riskManagerAuctionClosedEventEmitted(event))
    }
  )
}
