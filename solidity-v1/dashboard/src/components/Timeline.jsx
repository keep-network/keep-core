import React, { useContext } from "react"
import OnlyIf from "./OnlyIf"

export const TIMELINE_ELEMENT_STATUS = {
  INACTIVE: 0,
  SEMI_ACTIVE: 1,
  ACTIVE: 2,
}

const TimelineContext = React.createContext({
  // May be useful in the future, eg. if we would like to be able to display the
  // timeline events on alternating sides.
  position: "left",
})

const useTimelineContext = () => {
  const context = useContext(TimelineContext)
  if (!context) {
    throw new Error("TimelineContext used outside of Timeline component")
  }
  return context
}

const Timeline = ({ children, position = "left", className = "" }) => {
  return (
    <TimelineContext.Provider value={{ position: position }}>
      <ul className={`timeline ${className}`}>{children}</ul>
    </TimelineContext.Provider>
  )
}

const TimelineElement = ({ position, children, className = "" }) => {
  const { position: timelinePosition } = useTimelineContext()
  const _position = position || timelinePosition
  return (
    <li
      className={`timeline__element timeline__element--${_position} ${className}`}
    >
      {children}
    </li>
  )
}

Timeline.Element = TimelineElement

Timeline.Content = ({ children, className = "", ...props }) => {
  return (
    <div className={`element__content ${className}`} {...props}>
      {children}
    </div>
  )
}

Timeline.ElementDefaultCard = ({
  children,
  className = "",
  status = TIMELINE_ELEMENT_STATUS.ACTIVE,
  ...props
}) => {
  let statusClass = ""
  if (status === TIMELINE_ELEMENT_STATUS.ACTIVE) {
    statusClass = "timeline-default-card--active"
  } else if (status === TIMELINE_ELEMENT_STATUS.SEMI_ACTIVE) {
    statusClass = "timeline-default-card--semi-active"
  }

  return (
    <div
      className={`timeline-default-card ${statusClass} ${className}`}
      {...props}
    >
      {children}
    </div>
  )
}

Timeline.Breakpoint = ({ children, className = "" }) => {
  return <div className={`element__breakpoint ${className}`}>{children}</div>
}

Timeline.BreakpointDot = ({
  children,
  lineBreaker = false,
  lineBreakerColor = "grey-30",
  status = TIMELINE_ELEMENT_STATUS.ACTIVE,
  className = "",
  ...props
}) => {
  let statusClass = ""
  if (status === TIMELINE_ELEMENT_STATUS.ACTIVE) {
    statusClass = "breakpoint__dot--active"
  } else if (status === TIMELINE_ELEMENT_STATUS.SEMI_ACTIVE) {
    statusClass = "breakpoint__dot--semi-active"
  }

  return (
    <span
      className={`breakpoint__dot ${
        lineBreaker ? "breakpoint__dot--breaker" : ""
      } ${statusClass} 
      ${lineBreakerColor ? `breakpoint__dot--breaker-${lineBreakerColor}` : ""}
      ${className}`}
      {...props}
    >
      <OnlyIf condition={!lineBreaker}>{children}</OnlyIf>
    </span>
  )
}

Timeline.BreakpointLine = ({
  status = TIMELINE_ELEMENT_STATUS.INACTIVE,
  className = "",
}) => {
  let statusClass = ""
  if (status === TIMELINE_ELEMENT_STATUS.ACTIVE) {
    statusClass = "breakpoint__line--active"
  } else if (status === TIMELINE_ELEMENT_STATUS.SEMI_ACTIVE) {
    statusClass = "breakpoint__line--semi-active"
  }

  return <span className={`breakpoint__line ${statusClass} ${className}`} />
}

export default Timeline
