import React from "react"
import * as Icons from "./Icons"
import Tooltip from "./Tooltip"

export const SpeechBubbleTooltip = ({
  text,
  title,
  iconColor,
  iconBackgroundColor,
}) => {
  return (
    <Tooltip
      triggerComponent={() => (
        <Icons.Tooltip
          color={iconColor}
          backgroundColor={iconBackgroundColor}
        />
      )}
      title={title}
      content={text}
    />
  )
}

SpeechBubbleTooltip.defaultProps = {
  className: "",
}

export default React.memo(SpeechBubbleTooltip)
