import { useEffect } from "react"
import { useSelector } from "react-redux"
import { useWeb3Address } from "../components/WithWeb3Context"

/**
 * Custom hook that triggers a callback when account changed and the redux store
 * has been restarted and ready to handle actions.
 *
 * @param {Function} callback - A Function to trigger. Note: Make sure the
 * function is memoized otherwise it will be called every time the component
 * updates.
 */
export const useDispatchWhenAccountChanged = (callback) => {
  const isReady = useSelector((state) => state.app.isReady)
  const address = useWeb3Address()

  useEffect(() => {
    callback()
  }, [callback, address, isReady])
}
