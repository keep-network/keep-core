import React from 'react'
import SlashedTokensList from './SlashedTokensList'

const SlashedTokens = (props) => {
  return (
    <section id="slashed-tokens" className="tile">
      <h5>
          Slashed Tokens
        <span className="text-warning">Group misbehavior results in a slash of KEEP tokens.</span>
      </h5>
      <SlashedTokensList />
    </section>
  )
}

export default SlashedTokens
