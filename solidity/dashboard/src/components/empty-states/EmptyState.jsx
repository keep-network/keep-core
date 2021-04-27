import React from "react"
import ResourceTooltip from "../ResourceTooltip"
import ConnectWalletBtn from "../ConnectWalletBtn";

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
  return <ConnectWalletBtn text={text} btnClassName={btnClassName} />
}

export default EmptyState
