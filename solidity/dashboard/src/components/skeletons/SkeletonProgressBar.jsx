import React from "react"
import CircularProgressBar from "../CircularProgressBar"
import { colors } from "../../constants/colors"

const SkeletonProgressBar = ({
  circular = true,
  primaryColor = colors.grey30,
  secondaryColor = colors.grey10,
  fillingInPercentage = 50,
  barWidth = 24,
}) => {
  return (
    circular && (
      <CircularProgressBar
        value={fillingInPercentage}
        total={100}
        color={primaryColor}
        backgroundStroke={secondaryColor}
        barWidth={barWidth}
      />
    )
  )
}

export default SkeletonProgressBar
