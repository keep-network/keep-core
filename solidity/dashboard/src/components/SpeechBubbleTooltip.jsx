import React from 'react'
import * as Icons from './Icons'

export const SpeechBubbleTooltip = ({
  text,
  title,
  iconColor,
  iconBackgroundColor
}) => {
  return (
    <div className="flex row">
      <div className='tooltip'>
        <Icons.Tooltip color={iconColor} backgroundColor={iconBackgroundColor} />
        <span className="tooltip-text top">
          {text}
        </span>
      </div>
      {title && 
        <span
          style={{ marginLeft: '0.5rem'}}
          className="text-grey-60 text-caption">
            {title}
        </span>
      }
    </div>
  )
}

SpeechBubbleTooltip.defaultProps = {
  className: '',
}

export default React.memo(SpeechBubbleTooltip)
