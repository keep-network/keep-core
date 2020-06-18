import React from "react"
import SpeechBubbleTooltip from "./SpeechBubbleTooltip"

const Tile = ({
  title,
  titleStyle,
  titleClassName,
  withTooltip,
  tooltipProps,
  subtitle,
  children,
  ...sectionProps
}) => {
  return (
    <section className="tile" {...sectionProps}>
      <div className="flex center">
        <h4 className={titleClassName} style={titleStyle}>
          {title}
        </h4>
        {withTooltip && <SpeechBubbleTooltip {...tooltipProps} />}
      </div>
      <div className="text-grey-40 text-small">{subtitle}</div>
      <div className="mt-1">{children}</div>
    </section>
  )
}

Tile.defaultProps = {
  withTooltip: false,
  titleClassName: "mr-1 text-grey-70",
}

export default Tile
