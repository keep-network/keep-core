import React from "react"
import "../css/card.less"

const Card = ({ className, children }) => {
  return (
    <div className={`card-root ${className}`}>
      {children}
    </div>
  )
}

export default Card
