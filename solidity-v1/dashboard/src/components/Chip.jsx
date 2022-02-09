import React from "react"

const Chip = ({
  icon,
  text,
  size = "small",
  color = "emphasis",
  className = "",
  style = {},
}) => {
  return (
    <span
      className={`chip chip--${color} chip--${
        icon ? "icon" : size
      } ${className}`}
      style={style}
    >
      {icon ? icon : text}
    </span>
  )
}

export default React.memo(Chip)
