import React from "react"

const MasonryFlexContainer = ({ maxHeight, className = "", children }) => {
  const style = {
    maxHeight: maxHeight,
  }
  return (
    <div className={`masonry-flex-container ${className}`} style={style}>
      {children}
    </div>
  )
}

export default MasonryFlexContainer
