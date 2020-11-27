import React, { useEffect, useState, useRef } from "react"
import { CSSTransition } from "react-transition-group"

const Tooltip = ({
  triggerComponent: TriggerComponent,
  children,
  direction = "bottom",
  simple = false,
  delay = 300,
  className = "",
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
        in={active}
        timeout={delay}
        classNames="tooltip__content-wrapper"
        unmountOnExit
        onEnter={showTooltip}
        onExited={hideTooltip}
      >
        <div
          className="tooltip__content-wrapper"
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
