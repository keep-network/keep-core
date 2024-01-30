import { useEffect, useRef, useState } from "react"
import { usePrevious } from "./usePrevious"
import { KeepExplorerMode } from "../contracts"
import { useWeb3Context } from "../components/WithWeb3Context"

export const useSubscribeToExplorerModeContractEvent = (
  contractName,
  eventName,
  subscribeToEventCallback = () => {},
  options = {},
  subscribeOnMainnetOnly = false
) => {
  const event = useRef(null)
  const { isConnected, chainId } = useWeb3Context()
  const contract = KeepExplorerMode[contractName]?.instance
  const contractRef = useRef(null)
  const [latestEvent, setLatestEvent] = useState({})
  const previousEvent = usePrevious(latestEvent)

  useEffect(() => {
    if (subscribeOnMainnetOnly && chainId?.toString() !== "1") {
      console.warn(
        `Subscribing to ${eventName} event from ${contractName} contract is only available on mainnet.`
      )
      return
    }
    if (!isConnected) return

    if (!contract) {
      console.error(
        `Failed subscribing to ${eventName} event: ${contractName} contract was not found in KeepExplorerMode lib.`
      )
      return
    }

    if (!contract?.options?.address) {
      console.error(
        `Failed subscribing to ${eventName} event: ${contractName} contract instance doesn't have an address set.`
      )
      return
    }

    contractRef.current = contract
    const subscribeToEvent = (error, event) => {
      if (error) {
        return
      }
      setLatestEvent(event)
    }

    try {
      event.current = contractRef.current.events[eventName](
        { fromBlock: "latest", ...options },
        subscribeToEvent
      )
    } catch (error) {
      console.error(`Failed subscribing to ${eventName}`, error)
    }

    return () => {
      if (event.current) {
        event.current.unsubscribe()
      }
    }
  }, [eventName, contract, isConnected, subscribeOnMainnetOnly, chainId])

  useEffect(() => {
    if (previousEvent.transactionHash === latestEvent.transactionHash) {
      return
    }
    subscribeToEventCallback({ ...latestEvent })
  }, [previousEvent, latestEvent])

  return { latestEvent }
}
