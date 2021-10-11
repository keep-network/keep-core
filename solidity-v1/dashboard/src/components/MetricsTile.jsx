import React from "react"
import Tooltip from "./Tooltip"
import * as Icons from "./Icons"

const MetricsTile = ({ className, children, style = {} }) => {
  return (
    <div className={`tile tile--metrics ${className}`} style={style}>
      {children}
    </div>
  )
}

MetricsTile.Tooltip = ({ className, children, ...restTooltipProps }) => {
  return (
    <Tooltip
      simple
      delay={0}
      triggerComponent={Icons.MoreInfo}
      className={`tile--metrics__tooltip ${className}`}
      {...restTooltipProps}
    >
      {children}
    </Tooltip>
  )
}

export default MetricsTile
