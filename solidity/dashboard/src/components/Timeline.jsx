import React, { useContext } from "react"

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

const Timeline = ({ children, position = "left" }) => {
  return (
    <TimelineContext.Provider value={{ position: position }}>
      <div className="timeline">{children}</div>
    </TimelineContext.Provider>
  )
}
const TimelineElement = ({ position, children, className = "" }) => {
  const { position: timelinePosition } = useTimelineContext()
  const _position = position || timelinePosition
  return (
    <div className={`timeline__element timeline__element--${_position} ${className}`}>
      {children}
    </div>
  )
}
Timeline.Element = TimelineElement

Timeline.ElementBreakpoint = ({
  children,
  className = "",
  // start, center
  // TODO add support for a start position.
  position = "center",
}) => {
  return (
    <div
      className={`element__breakpoint element__breakpoint--${position} ${className}`}
    >
      {children}
    </div>
  )
}

export default Timeline
