import React from "react"
import "../css/card.less"

const Card = ({ children }) => {
  return (
    <div className={'card-root'}>
      {children}
    </div>
  )
}

export default Card
