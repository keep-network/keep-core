import React from 'react'

export const BADGE_STATUS = {
  PENDING: { textClassName: 'text-warning text-normal', bgClassName: 'text-bg-pending-light' },
  COMPLETED: { textClassName: 'text-success', bgClassName: 'text-bg-success-light' },
}

const StatusBadge = ({ status, text, className }) => {
  return (
    <div
      className={`${status.textClassName} ${status.bgClassName} text-label ${className}`}
      style={{ padding: '0.1rem 0.2rem' }}
    >
      {text}
    </div>
  )
}

export default React.memo(StatusBadge)
