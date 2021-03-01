import React, { useEffect, useCallback } from "react"

export const useHideComponent = ({ hide = false }) => {
  const [isVisible, setIsVisible] = React.useState(true)

  useEffect(() => {
    setIsVisible(!hide)
  }, [hide])

  const hideComponent = useCallback(() => {
    setIsVisible(false)
  }, [])

  return [isVisible, hideComponent]
}
