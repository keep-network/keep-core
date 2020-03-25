import React from 'react'
import Button from './Button'

export const SeeAllButton = ({ showAll, onClickCallback, previewDataCount, dataLength }) => {
  if (dataLength <= previewDataCount || dataLength === 0) {
    return null
  }

  return (
    <Button
      className="btn btn-secondary see-all-btn"
      onClick={onClickCallback}
    >
      {showAll ? 'see less' : `see all (${dataLength - previewDataCount})`}
    </Button>
  )
}
