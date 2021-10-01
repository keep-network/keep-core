import React from "react"

const Badge = ({ type = "primary", style, text }) => {
  return (
    <div className={`badge--${type}`} style={style}>
      <span className="badge__text">{text}</span>
    </div>
  )
}

export default React.memo(Badge)
