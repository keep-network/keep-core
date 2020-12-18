import React from "react"
import "../css/card-container.less"

const CardContainer = ({ className, children }) => {
  return (
    <div className={`card-container-root ${className}`}>
      {children}
    </div>
  )
}

export default CardContainer
