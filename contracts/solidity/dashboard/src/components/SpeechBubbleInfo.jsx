import React from 'react'
import * as Icons from './Icons'

const childStyle = { marginLeft: '1rem' }
const wrapperStyle = { marginTop: '0.8rem' }

export const SpeechBubbleInfo = ({ children, className }) => {
  return (
    <div className={`flex flex-row ${className}`} style={wrapperStyle}>
      <Icons.SpeechBubble />
      <div className="text-small text-darker-grey" style={childStyle}>
        {children}
      </div>
    </div>
  )
}

SpeechBubbleInfo.defaultProps = {
  className: '',
}

export default React.memo(SpeechBubbleInfo)
