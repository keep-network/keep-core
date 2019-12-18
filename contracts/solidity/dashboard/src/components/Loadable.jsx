import React from 'react'

const Loadable = ({ text }) => (
    <div>{text}</div>
)

Loadable.defaultProps = {
    text: 'Loading...'
}

export const ClockIndicator = (props) => (<div className='indicator-clock'/>)

export default Loadable