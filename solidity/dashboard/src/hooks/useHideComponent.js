import React from "react"
import { useEffect } from "react"

export const useHideComponent = ({ hide = false }) => {
  const [isVisible, setIsVisible] = React.useState(true)

  useEffect(() => {
    setIsVisible(!hide)
  }, [hide])

  const hideComponent = () => {
    setIsVisible(false)
  }

  return [isVisible, hideComponent]
}
