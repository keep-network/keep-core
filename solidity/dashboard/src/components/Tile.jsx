import React from 'react'

const Tile = ({ title, children, ...sectionProps }) => {
    return (
        <section className="tile" {...sectionProps}>
            <h4 className="mb-1 text-grey-70">{title}</h4>
            {children}
        </section>
    )
}

export default Tile
