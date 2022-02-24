import React from "react"
import { Link } from "react-router-dom"
import Tooltip from "./Tooltip"
import * as Icons from "./Icons"
import { colors } from "../constants/colors"

export const ResourceTooltipContent = ({
  title,
  content,
  redirectLink = "/resources/quick-terminology",
  linkText = "Learn more in Resources",
  withRedirectLink = true,
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
      {withRedirectLink && (
        <Link to={redirectLink} className="internal text-small">
          {linkText}
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
      triggerComponent={() => <Icons.MoreInfo className={"resource-tooltip"} />}
    >
      <ResourceTooltipContent {...restProps} />
    </Tooltip>
  )
}

export default ResourceTooltip
