import React from "react"

const Chip = ({
  icon,
  text,
  size = "small",
  color = "emphasis",
  className = "",
}) => {
  return (
    <span
      className={`chip--${color} chip--${icon ? "icon" : size} ${className}`}
    >
      {icon ? icon : text}
    </span>
  )
}

export default React.memo(Chip)
