import React from 'react'

export const BADGE_STATUS = {
  PENDING: { textClassName: 'text-warning text-normal', bgClassName: 'text-bg-pending-light' },
  COMPLETED: { textClassName: 'text-success', bgClassName: 'text-bg-success-light' },
}

const badgeStyle = { padding: '0.1rem 0.5rem', borderRadius: '100px' }

const StatusBadge = ({ status, text, className }) => {
  return (
    <div
      className={`${status.textClassName} ${status.bgClassName} text-label ${className}`}
      style={badgeStyle}
    >
      {text}
    </div>
  )
}

export default React.memo(StatusBadge)
