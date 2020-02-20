import React, { useMemo } from 'react'
import { formatDate } from '../utils'

const Timeline = ({ title, breakpoints, footer }) => {
  const breakpointComponents = useMemo(() => {
    return breakpoints.map(renderBreakpoint)
  }, [breakpoints])

  return (
    <React.Fragment>
      <div className="text-title text-darker-grey">{title}</div>
      <section className="timeline">
        {breakpointComponents}
      </section>
      {footer}
    </React.Fragment>
  )
}

const TimelineBreakpoint = ({ label, date, dotColorClassName }) => {
  return (
    <div className={`breakpoint ${dotColorClassName || ''}`}>
      <div className="breakpoint-content">
        <div className="text-big text-black">{label}</div>
        <div className="text-small text-grey">
          {formatDate(date)}
        </div>
      </div>
    </div>
  )
}

const renderBreakpoint = (item, index) => <TimelineBreakpoint key={index} {...item}/>

export default React.memo(Timeline)
