import React from 'react'

const Loadable = ({ text, isFetching, children }) => (
  isFetching ? <div>{text}</div> : children
)

Loadable.defaultProps = {
  text: 'Loading...',
  isFetching: true,
}

export const ClockIndicator = ({ color }) => (<div className={`indicator-clock ${color}`}/>)
ClockIndicator.defaultProps = {
  color: '',
}

export const LoadingOverlay = React.memo(({ isFetching, classNames, children }) => {
  return (
    <div className={`loading-overlay-container ${classNames}`}>
      {children}
      <div className={`loading-overlay${isFetching ? '' : ' hidden'}`}>
        Loading...
      </div>
    </div>
  )
})

export default Loadable
