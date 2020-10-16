import React from "react"
import * as Icons from "./Icons"
import { colors } from "../constants/colors"

export const BANNER_TYPE = {
  SUCCESS: {
    className: "success",
    iconComponent: <Icons.OK color={colors.success} />,
  },
  PENDING: {
    className: "pending",
    iconComponent: (
      <Icons.PendingBadge bgColor={colors.bgPending} color={colors.pending} />
    ),
  },
  ERROR: {
    className: "error",
    iconComponent: <Icons.Cross color={colors.error} height={12} width={12} />,
  },
  DISABLED: { className: "disabled", iconComponent: null },
  NOTIFICATION: {
    className: "notification",
    iconComponent: <Icons.Alert width={24} height={24} />,
  },
}

const Banner = ({ inline, className, children, ...restProps }) => {
  return (
    <div className={`banner${inline ? "--inline" : ""} ${className}`}>
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

Banner.Description = ({ onClick, children, className = "" }) => {
  return <div className={`banner__description ${className}`}>{children}</div>
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

Banner.Icon = ({ icon: IconComponent, className = "" }) => {
  return <IconComponent className={`banner__icon ${className}`} />
}

export default Banner
