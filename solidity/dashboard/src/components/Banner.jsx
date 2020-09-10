import React from "react"
import * as Icons from "./Icons"
import { colors } from "../constants/colors"
import ViewTxMsgContent from "./ViewTxMsgContent"

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

const Banner = ({
  type,
  title,
  onTitleClick,
  titleClassName,
  subtitle,
  withIcon,
  withCloseIcon,
  onCloseIcon,
  children,
  withTransactionHash,
}) => {
  return (
    <div className={`banner banner-${type.className}`}>
      {withIcon && <div className="banner-icon flex">{type.iconComponent}</div>}
      <div className="banner-content-wrapper">
        {title && (
          <div
            className={`banner-title ${titleClassName}`}
            onClick={onTitleClick}
          >
            {title}
          </div>
        )}
        {subtitle && (
          <div className="banner-subtitle">
            {withTransactionHash ? (
              <ViewTxMsgContent txHash={subtitle} />
            ) : (
              subtitle
            )}
          </div>
        )}
      </div>
      {withCloseIcon && (
        <div className="banner-close-icon" onClick={onCloseIcon}>
          <Icons.Cross color={colors.grey70} height={12} width={12} />
        </div>
      )}
      {children}
    </div>
  )
}

Banner.defaultProps = {
  onTitleClick: () => {},
  titleClassName: "",
  withIcon: false,
  withCloseIcon: false,
  onCloseIcon: () => {},
  children: null,
}

export default React.memo(Banner)
