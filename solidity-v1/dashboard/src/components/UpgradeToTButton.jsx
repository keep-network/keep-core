import React from "react"
import { LINK } from "../constants/constants"

const UpgradeToTButton = ({ className }) => {
  return (
    <a
      href={LINK.thresholdDapp}
      rel="noopener noreferrer"
      target="_blank"
      className={`btn ${className}`}
    >
      upgrade to t ↗
    </a>
  )
}

export default UpgradeToTButton
