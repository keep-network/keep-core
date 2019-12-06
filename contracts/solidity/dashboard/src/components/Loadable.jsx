import React from 'react'

const Loadable = ({ text }) => (
    <div>{text}</div>
)

Loadable.defaultProps = {
    text: 'Loadging...'
}

export default Loadable