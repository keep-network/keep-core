import React from "react"
import * as Icons from "./Icons"

const Banner = ({ inline, className, children, style = {}, ...restProps }) => {
  return (
    <div
      className={`banner${inline ? "--inline" : ""} ${className}`}
      style={style}
    >
      {inline ? (
        <>
          <Banner.Icon icon={restProps.icon} />
          <Banner.Title>{restProps.title}</Banner.Title>
        </>
      ) : (
        children
      )}
    </div>
  )
}

Banner.Title = ({ onClick, children, className = "" }) => {
  return (
    <div className={`banner__title ${className}`} onClick={onClick}>
      {children}
    </div>
  )
}

Banner.Description = ({ children, className = "", ...restProps }) => {
  return (
    <div className={`banner__description ${className}`} {...restProps}>
      {children}
    </div>
  )
}

Banner.Action = ({ children, onClick, icon, className = "" }) => {
  return (
    <div className={`banner__action ${className}`} onClick={onClick}>
      {icon && <Icons.KeepOutline className="banner__action__icon" />}
      {children}
    </div>
  )
}

Banner.CloseIcon = ({
  icon: IconComponent = Icons.Cross,
  className = "",
  onClick,
}) => {
  return (
    <IconComponent
      className={`banner__close-icon ${className}`}
      onClick={onClick}
    />
  )
}

Banner.Icon = ({ icon: IconComponent, className = "", ...iconProps }) => {
  return (
    <IconComponent className={`banner__icon ${className}`} {...iconProps} />
  )
}

export default Banner
