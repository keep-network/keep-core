import React from 'react'
import Button from './Button'

export const SeeAllButton = ({ showAll, onClickCallback, previewDataCount, dataLength }) => {
  if (dataLength >= previewDataCount || dataLength === 0) {
    return null
  }

  return (
    <Button
      className="btn btn-default see-all-btn"
      onClick={onClickCallback}
    >
      {showAll ? 'SEE LESS' : `SEE ALL (${dataLength - previewDataCount})`}
    </Button>
  )
}
