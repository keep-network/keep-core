import { useEffect, useRef, useContext, useState } from "react"
import { Web3Context } from "../components/WithWeb3Context"
import { usePrevious } from "./usePrevious"

export const useSubscribeToContractEvent = (
  contractName,
  eventName,
  subscribeToEventCallback = () => {}
) => {
  const event = useRef(null)
  const contract = useRef(useContext(Web3Context)[contractName])
  const [latestEvent, setLatestEvent] = useState({})
  const previousEvent = usePrevious(latestEvent)

  useEffect(() => {
    const subscribeToEvent = (error, event) => {
      if (error) {
        return
      }
      setLatestEvent(event)
    }
    event.current = contract.current.events[eventName](subscribeToEvent)

    return () => {
      event.current.unsubscribe()
    }
  }, [eventName])

  useEffect(() => {
    if (previousEvent.transactionHash === latestEvent.transactionHash) {
      return
    }
    subscribeToEventCallback({ ...latestEvent })
  })

  return { latestEvent }
}
