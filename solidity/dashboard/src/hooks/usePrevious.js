import { useRef, useEffect } from "react"

export const usePrevious = (value) => {
  const valueRef = useRef(value)
  useEffect(() => {
    valueRef.current = value
  }, [value])

  return valueRef.current
}
