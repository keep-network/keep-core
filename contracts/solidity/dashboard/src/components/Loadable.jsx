import React from 'react'

const Loadable = ({ text, isFetching, children }) => (
  isFetching ? <div>{text}</div> : children
)

Loadable.defaultProps = {
  text: 'Loading...',
  isFetching: true,
}

export const ClockIndicator = (props) => (<div className='indicator-clock'/>)

export default Loadable
