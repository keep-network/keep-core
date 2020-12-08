import React from "react"
import Skeleton from "./Skeleton"

const CircleSkeleton = ({
  width,
  height,
  shining = false,
  color = "grey-20",
}) => {
  return (
    <Skeleton
      shining={shining}
      color={color}
      styles={{ width, height, borderRadius: "50%" }}
    />
  )
}

export default CircleSkeleton
