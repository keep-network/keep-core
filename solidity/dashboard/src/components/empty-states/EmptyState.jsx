import React from "react"
import WalletOptions from "../WalletOptions"
import Tooltip from "../Tooltip"
import ResourceTooltip from "../ResourceTooltip"

const EmptyState = ({ children }) => {
  return <section className="empty-state__wrapper">{children}</section>
}

EmptyState.Title = ({
  text,
  tag = "h2",
  wrapperClassName = "",
  classname = "",
  tooltipProps = null,
}) => {
  const Tag = tag
  return (
    <header className={`empty-state__header ${wrapperClassName}`}>
      <Tag className={`empty-state__header__title ${classname}`}>{text}</Tag>
      {tooltipProps && (
        <ResourceTooltip
          tooltipClassName="empty-state__header__tooltip"
          {...tooltipProps}
        />
      )}
    </header>
  )
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
