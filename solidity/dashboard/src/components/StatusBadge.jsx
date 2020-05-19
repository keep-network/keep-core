import React from "react"
import { PENDING_STATUS, COMPLETE_STATUS } from "../constants/constants"
import * as Icons from "./Icons"

export const BADGE_STATUS = {
  [PENDING_STATUS]: {
    textClassName: "text-grey-70 text-normal",
    bgClassName: "bg-pending",
    icon: <Icons.PendingBadge />,
  },
  [COMPLETE_STATUS]: {
    textClassName: "text-success",
    bgClassName: "bg-success",
    icon: <Icons.OKBadge />,
  },
  DISABLED: { textClassName: "text-grey-50", bgClassName: "bg-grey-10" },
}

const badgeStyle = { padding: "0.1rem 0.5rem", borderRadius: "100px" }

const StatusBadge = ({ status, text, className, onlyIcon }) => {
  return onlyIcon ? (
    <span className="flex row center">
      {status.icon}
      <span style={{ marginLeft: "0.5rem" }}>{text}</span>
    </span>
  ) : (
    <span
      className={`${status.textClassName} ${status.bgClassName} text-label text-normal ${className}`}
      style={badgeStyle}
    >
      {text}
    </span>
  )
}

StatusBadge.defaultProps = {
  onlyIcon: false,
}

export default React.memo(StatusBadge)
