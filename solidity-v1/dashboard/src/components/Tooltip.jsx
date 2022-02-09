import React, { useEffect, useState, useRef } from "react"
import { CSSTransition } from "react-transition-group"

export const TOOLTIP_DIRECTION = {
  TOP: "top",
  BOTTOM: "bottom",
}

const Tooltip = ({
  triggerComponent: TriggerComponent,
  children,
  direction = TOOLTIP_DIRECTION.BOTTOM,
  simple = false,
  delay = 300,
  className = "",
  shouldShowTooltip = true,
  contentWrapperStyles = {},
  tooltipContentWrapperClassName = "",
}) => {
  const timeout = useRef(null)
  const [active, setActive] = useState(false)
  const contentRef = useRef(null)
  const arrowRef = useRef(null)

  useEffect(() => {
    return () => {
      if (timeout.current) {
        clearTimeout(timeout.current)
      }
    }
  })

  /**
   * Check if tooltip is out of the RIGHT side of the view.
   * If it is it will move the tooltip to the left.
   */
  const handleDropdownPosition = () => {
    if (!contentRef || !contentRef.current) return
    if (!arrowRef || !arrowRef.current) return
    const contentRect = contentRef.current.getBoundingClientRect()
    const contentRightX = contentRect.x + contentRect.width

    if (contentRightX > window.outerWidth) {
      contentRef.current.style.transform = `translateX(calc(-100% + 9px))`
      if (direction === TOOLTIP_DIRECTION.TOP) {
        contentRef.current.style.borderBottomRightRadius = "0"
      } else if (direction === TOOLTIP_DIRECTION.BOTTOM) {
        contentRef.current.style.borderTopRightRadius = "0"
      }

      arrowRef.current.style.left = "auto"
      arrowRef.current.style.right = "calc(0%)"
    }
  }

  const showTooltip = () => {
    handleDropdownPosition()
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
          ref={contentRef}
        >
          <div ref={arrowRef} className={"tooltip__arrow"} />
          {children}
        </div>
      </CSSTransition>
    </div>
  )
}

Tooltip.Divider = () => <hr className="tooltip__divider" />

Tooltip.Header = ({ text, className = "" }) => (
  <div className={`tooltip__header ${className}`}>
    <div className="tooltip__header__title">{text}</div>
  </div>
)

Tooltip.Content = ({ children }) => (
  <div className="tooltip__content">{children}</div>
)

export default Tooltip
