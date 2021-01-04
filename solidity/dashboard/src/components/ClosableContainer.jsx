import React, { useEffect } from "react"
import * as Icons from "./Icons"

const ClosableContainer = ({ hide = false, className, children }) => {
  const [showComponent, setShowComponent] = React.useState(true)

  useEffect(() => {
    setShowComponent(!hide)
  }, [hide])

  const hideComponent = () => {
    setShowComponent(false)
  }

  return showComponent ? (
    <div className={`closable-container ${className}`}>
      <Icons.Cross
        className={`closable-container__close-icon`}
        onClick={hideComponent}
      />
      {children}
    </div>
  ) : (
    <></>
  )
}

export default ClosableContainer
