import React from "react"
import WalletOptions from "../WalletOptions"
import Tooltip from "../Tooltip"

const EmptyState = ({ children }) => {
  return <section className="empty-state__wrapper">{children}</section>
}

EmptyState.Title = ({ text, tag = "h2", classname = "" }) => {
  const Tag = tag
  return <Tag className={`empty-state__title ${classname}`}>{text}</Tag>
}

EmptyState.Subtitle = ({ text, tag = "h3", className = "" }) => {
  const Tag = tag
  return <Tag className={`empty-state__subtitle ${className}`}>{text}</Tag>
}

EmptyState.Skeleton = ({ children, className = "" }) => {
  return <div className={`empty-state__skeleton ${className}`}>{children}</div>
}

EmptyState.ConnectWalletBtn = ({
  text = "connect wallet",
  btnClassName = "",
}) => {
  // TODO connect to a wallet onClick
  return (
    <Tooltip
      direction="top"
      simple
      className="empty-state__wallet-options-tooltip"
      triggerComponent={() => (
        <span
          className={`btn btn-primary btn-lg empty-state__connect-wallet-btn ${btnClassName}`}
        >
          {text}
        </span>
      )}
    >
      <WalletOptions />
    </Tooltip>
  )
}

export default EmptyState
