import React from "react"

const Skeleton = ({
  tag = "h2",
  color = "grey-30",
  width = "100%",
  className = "",
  styles = {},
  shining = false,
}) => {
  const Tag = tag

  return (
    <Tag
      className={`skeleton--${color} ${shining ? "shining" : ""} ${className}`}
      style={{ width, ...styles }}
    />
  )
}

export default Skeleton
