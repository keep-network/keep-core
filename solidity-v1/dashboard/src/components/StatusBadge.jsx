import React from "react"
import { PENDING_STATUS, COMPLETE_STATUS } from "../constants/constants"
import * as Icons from "./Icons"
import ReactTooltip from "react-tooltip"
import OnlyIf from "./OnlyIf"

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
  ACTIVE: {
    textClassName: "text-grey-70 text-normal",
    bgClassName: "bg-success-light",
    icon: <Icons.OKBadge />,
  },
  ERROR: {
    textClassName: "text-black",
    bgClassName: "bg-error",
  },
}

const badgeStyle = { padding: "0.1rem 0.5rem", borderRadius: "100px" }

const StatusBadge = ({
  status,
  text,
  className,
  onlyIcon,
  bgClassName,
  withTooltip = false,
  tooltipId,
  tooltipProps = {},
}) => {
  return onlyIcon ? (
    <span className="flex row center">
      {status.icon}
      <span style={{ marginLeft: "0.5rem" }}>{text}</span>
    </span>
  ) : (
    <span
      className={`${status.textClassName} ${
        bgClassName || status.bgClassName
      } text-label text-normal ${className}`}
      style={badgeStyle}
      data-tip
      data-for={tooltipId}
    >
      {text}
      <OnlyIf condition={withTooltip}>
        <ReactTooltip id={tooltipId} {...tooltipProps}>
          <span>
            The stake amount is not yet confirmed. Click “Stake” to confirm the
            stake amount. This stake is not staked on Threshold until it is
            confirmed.
          </span>
        </ReactTooltip>
      </OnlyIf>
    </span>
  )
}

StatusBadge.defaultProps = {
  onlyIcon: false,
}

export default React.memo(StatusBadge)
