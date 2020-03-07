import React from 'react'
import { PENDING_STATUS, COMPLETE_STATUS } from '../constants/constants'

export const BADGE_STATUS = {
  [PENDING_STATUS]: { textClassName: 'text-pending text-normal', bgClassName: 'bg-pending' },
  [COMPLETE_STATUS]: { textClassName: 'text-success', bgClassName: 'bg-success' },
}

const badgeStyle = { padding: '0.1rem 0.5rem', borderRadius: '100px' }

const StatusBadge = ({ status, text, className }) => {
  return (
    <div
      className={`${status.textClassName} ${status.bgClassName} text-label text-normal ${className}`}
      style={badgeStyle}
    >
      {text}
    </div>
  )
}

export default React.memo(StatusBadge)
