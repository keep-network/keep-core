import React from "react"
import Skeleton from "./Skeleton"

const CircleSkeleton = ({ width, height }) => {
  return <Skeleton styles={{ width, height, borderRadius: "50%" }} />
}

export default CircleSkeleton
