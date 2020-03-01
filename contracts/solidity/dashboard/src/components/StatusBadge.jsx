import React from 'react'

export const BADGE_STATUS = {
  PENDING: { textClassName: 'text-pending text-normal', bgClassName: 'bg-pending' },
  COMPLETED: { textClassName: 'text-success', bgClassName: 'bg-success' },
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
