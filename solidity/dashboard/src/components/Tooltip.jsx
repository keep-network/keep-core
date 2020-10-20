import React, { useEffect, useState, useRef } from "react"
import * as Icons from "./Icons"

const Tooltip = ({
  triggerComponent: TriggerComponent,
  title,
  content,
  children,
  direction = "bottom",
  icon: IconComponent = Icons.Tooltip,
  simple = false,
  delay = 300,
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
    timeout.current = setTimeout(() => {
      setActive(true)
    }, delay)
  }

  const hideTooltip = () => {
    clearInterval(timeout.current)
    setActive(false)
  }

  return (
    <div className={`tooltip${simple ? "--simple" : ""}--${direction}`}>
      <div
        className="tooltip__trigger"
        onMouseLeave={hideTooltip}
        onMouseEnter={showTooltip}
      >
        <TriggerComponent />
      </div>
      {active && (
        <div className="tooltip__content-wrapper">
          {children ? (
            children
          ) : (
            <>
              <Tooltip.Header text={title} icon={IconComponent} />
              <Tooltip.Divider />
              <Tooltip.Content>{content}</Tooltip.Content>
            </>
          )}
        </div>
      )}
    </div>
  )
}

Tooltip.Divider = () => <hr className="tooltip__divider" />

Tooltip.Header = ({ icon: IconComponent, text, className = "" }) => (
  <div className={`tooltip__header ${className}`}>
    <IconComponent className="tooltip__header__icon" />
    <div className="tooltip__header__title">{text}</div>
  </div>
)

Tooltip.Content = ({ children }) => (
  <div className="tooltip__content">{children}</div>
)

export default Tooltip
