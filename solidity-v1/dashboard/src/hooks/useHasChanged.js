import { usePrevious } from "./usePrevious"

/**
 * Checks if value has changed since the last component re-render
 * @param {*} val - current value
 * @return {boolean}
 */
const useHasChanged = (val) => {
  const prevVal = usePrevious(val)
  return prevVal !== val
}

export default useHasChanged
