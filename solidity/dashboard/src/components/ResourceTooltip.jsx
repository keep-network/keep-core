import React from "react"
import { Link } from "react-router-dom"
import Tooltip from "./Tooltip"
import * as Icons from "./Icons"
import { colors } from "../constants/colors"

export const ResourceTooltipContent = ({
  title,
  content,
  redirectLink = "/resources/quick-terminology",
  btnText = "go to resources",
  withRedirectButton = true,
}) => {
  return (
    <>
      <Tooltip.Header
        text={title}
        icon={Icons.Tooltip}
        iconProps={{ color: colors.grey70, backgroundColor: colors.mint20 }}
      />
      <Tooltip.Divider />
      <Tooltip.Content>{content}</Tooltip.Content>
      {withRedirectButton && (
        <Link to={redirectLink} className="btn btn-secondary btn-sm mt-2">
          {btnText}
        </Link>
      )}
    </>
  )
}

const ResourceTooltip = ({
  iconColor = colors.grey70,
  iconBackgroundColor = colors.mint20,
  tooltipClassName = "",
  ...restProps
}) => {
  return (
    <Tooltip
      className={tooltipClassName}
      triggerComponent={() => (
        <Icons.Tooltip
          color={iconColor}
          backgroundColor={iconBackgroundColor}
        />
      )}
    >
      <ResourceTooltipContent {...restProps} />
    </Tooltip>
  )
}

export default ResourceTooltip
