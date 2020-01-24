import React from 'react'

const Loadable = ({ text, isFetching, children }) => (
  isFetching ? <div>{text}</div> : children
)

Loadable.defaultProps = {
  text: 'Loading...',
  isFetching: true,
}

export const ClockIndicator = (props) => (<div className='indicator-clock'/>)

export const LoadingOverlay = ({ isFetching, classNames, children }) => {
  if (!isFetching) {
    return children
  }

  return (
    <div className={`loading-overlay-container ${classNames}`}>
      <div className='loading-overlay'>
        Loading...
      </div>
      {children}
    </div>
  )
}

export default Loadable
