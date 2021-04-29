import React, { useEffect, useState, useRef } from "react"
import { CSSTransition } from "react-transition-group"

const Tooltip = ({
  triggerComponent: TriggerComponent,
  children,
  direction = "bottom",
  simple = false,
  delay = 300,
  className = "",
  shouldShowTooltip = true,
  contentWrapperStyles = {},
  tooltipContentWrapperClassName = "",
}) => {
  const timeout = useRef(null)
  const [active, setActive] = useState(false)

  useEffect(() => {
    return () => {
      if (timeout.current) {
        clearTimeout(timeout.current)
      }
    }
  })

  const showTooltip = () => {
    setActive(true)
    clearTimeout(timeout.current)
  }

  const hideTooltip = () => {
    timeout.current = setTimeout(() => {
      setActive(false)
    }, delay)
  }

  const _wrapperClassName = `tooltip__content-wrapper ${
    tooltipContentWrapperClassName ? tooltipContentWrapperClassName : ""
  }`

  return (
    <div
      className={`tooltip${
        simple ? "--simple" : ""
      }--${direction} ${className}`}
    >
      <div
        className="tooltip__trigger"
        onMouseEnter={showTooltip}
        onMouseLeave={hideTooltip}
      >
        <TriggerComponent />
      </div>
      <CSSTransition
        in={shouldShowTooltip ? active : false}
        timeout={delay}
        classNames={_wrapperClassName}
        unmountOnExit
        onEnter={showTooltip}
        onExited={hideTooltip}
      >
        <div
          className={_wrapperClassName}
          style={contentWrapperStyles}
          onMouseEnter={showTooltip}
          onMouseLeave={hideTooltip}
        >
          {children}
        </div>
      </CSSTransition>
    </div>
  )
}

Tooltip.Divider = () => <hr className="tooltip__divider" />

Tooltip.Header = ({
  icon: IconComponent,
  text,
  className = "",
  iconProps = {},
}) => (
  <div className={`tooltip__header ${className}`}>
    <IconComponent className="tooltip__header__icon" {...iconProps} />
    <div className="tooltip__header__title">{text}</div>
  </div>
)

Tooltip.Content = ({ children }) => (
  <div className="tooltip__content">{children}</div>
)

export default Tooltip
