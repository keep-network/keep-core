import { useEffect } from "react"
import { useSelector } from "react-redux"
import { usePrevious } from "./usePrevious"

export const useDispatchWhenAccountChanged = (dispatchFn) => {
  const isReady = useSelector((state) => state.app.isReady)
  const previousIsAppReady = usePrevious(isReady)

  useEffect(() => {
    dispatchFn()
  }, [dispatchFn])

  useEffect(() => {
    if (!previousIsAppReady && isReady) {
      dispatchFn()
    }
  })
}
