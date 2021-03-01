import React from "react"

const Card = ({ className, children }) => {
  return <div className={`card-root ${className}`}>{children}</div>
}

export default Card
