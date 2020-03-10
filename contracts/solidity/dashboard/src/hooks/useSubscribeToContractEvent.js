import { useEffect, useRef, useContext, useState } from 'react'
import { Web3Context } from '../components/WithWeb3Context'
import { isEmptyObj } from '../utils/general.utils'

export const useSubscribeToContractEvent = (contractName, eventName, subscribeToEventCallback = () => {}) => {
  const web3Context = useContext(Web3Context)
  const event = useRef(null)
  const contract = web3Context[contractName]
  const [latestEvent, setLatestEvent] = useState({})

  useEffect(() => {
    const subscribeToEvent = (error, event) => {
      if (error) {
        return
      }
      setLatestEvent(event)
    }
    event.current = contract.events[eventName](subscribeToEvent)

    return () => {
      event.current.unsubscribe((error, suscces) => console.log('unsub', error, suscces))
    }
  }, [])

  useEffect(() => {
    if (isEmptyObj(latestEvent)) {
      return
    }
    subscribeToEventCallback(latestEvent)
  }, [latestEvent])

  return { latestEvent }
}
