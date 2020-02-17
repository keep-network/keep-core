import React, { useMemo } from 'react'
import { formatDate } from '../utils'

const Timeline = ({ title, breakepoints }) => {
  const breakePointComponents = useMemo(() => {
    return breakepoints.map(renderBreakPoint)
  }, [breakepoints])

  return (
    <React.Fragment>
      <div className="text-title text-darker-grey">{title}</div>
      <section className="timeline">
        {breakePointComponents}
      </section>
    </React.Fragment>
  )
}

const TimelineBreakPoint = ({ label, date, dotColorClassName }) => {
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

const renderBreakPoint = (item, index) => <TimelineBreakPoint key={index} {...item}/>

export default React.memo(Timeline)
