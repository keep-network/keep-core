import React from "react"

const DoubleIcon = ({ MainIcon, SecondaryIcon, className }) => {
  return (
    <div className={`double-icon-container ${className}`}>
      <SecondaryIcon className={`secondary-icon`} />
      <MainIcon className={`main-icon`} />
    </div>
  )
}

export default DoubleIcon
