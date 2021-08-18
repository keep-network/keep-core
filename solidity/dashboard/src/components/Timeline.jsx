import React, { useContext } from "react"
import OnlyIf from "./OnlyIf"

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
  active = true,
  ...props
}) => {
  const activeClass = active ? "timeline-default-card--active" : ""
  return (
    <div
      className={`timeline-default-card ${activeClass} ${className}`}
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
  active = true,
  className = "",
  ...props
}) => {
  return (
    <span
      className={`breakpoint__dot ${
        lineBreaker ? "breakpoint__dot--breaker" : ""
      } ${active ? "breakpoint__dot--active" : ""} 
      ${lineBreakerColor ? `breakpoint__dot--breaker-${lineBreakerColor}` : ""}
      ${className}`}
      {...props}
    >
      <OnlyIf condition={!lineBreaker}>{children}</OnlyIf>
    </span>
  )
}

Timeline.BreakpointLine = ({ active = false, className = "" }) => {
  return (
    <span
      className={`breakpoint__line ${
        active ? "breakpoint__line--active" : ""
      } ${className}`}
    />
  )
}

export default Timeline
