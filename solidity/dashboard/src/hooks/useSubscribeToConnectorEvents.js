import { useEffect } from "react"
import { useDispatch } from "react-redux"
import { useWeb3Context } from "../components/WithWeb3Context"

const useSubscribeToConnectorEvents = () => {
  const dispatch = useDispatch()
  const { isConnected, connector, yourAddress } = useWeb3Context()

  useEffect(() => {
    const accountChangedHandler = (address) => {
      dispatch({ type: "app/account_changed", payload: { address } })
    }

    const disconnectHandler = () => {
      dispatch({ type: "app/logout" })
    }

    if (isConnected && connector) {
      dispatch({ type: "app/login", payload: { address: yourAddress } })
      connector.on("accountsChanged", accountChangedHandler)
      connector.once("disconnect", disconnectHandler)
    }

    return () => {
      if (connector) {
        connector.removeListener("accountsChanged", accountChangedHandler)
        connector.removeListener("disconnect", disconnectHandler)
      }
    }
  }, [isConnected, connector, dispatch, yourAddress])
}

export default useSubscribeToConnectorEvents
