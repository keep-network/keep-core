import React from "react"
import "../css/card-container.less"

const CardContainer = ({ children }) => {
  return (
    <div className={'card-container-root'}>
      {children}
    </div>
  )
}

export default CardContainer
