import React from 'react'

export const BADGE_STATUS = {
  PENDING: { textClassName: 'text-warning', bgClassName: 'text-bg-pending-light' },
  COMPLETE: { textClassName: 'text-success', bgClassName: 'text-bg-success-light' },
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
