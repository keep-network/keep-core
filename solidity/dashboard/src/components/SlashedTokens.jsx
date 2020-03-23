import React from 'react'
import SlashedTokensList from './SlashedTokensList'
import { LoadingOverlay } from './Loadable'

const SlashedTokens = (props) => {
  return (
    <LoadingOverlay isFetching={false}>
      <section id="slashed-tokens" className="tile">
        <h3 className="text-grey-60 mb-1">
          Slashed Tokens
        </h3>
        <span className="text-small text-darker-grey">
          A slash is a penalty for signing group misbehavior.
          A slash results in a removal of a portion of your delegated KEEP tokens.
          You can see a record below of all slashes.
        </span>
        <SlashedTokensList />
      </section>
    </LoadingOverlay>

  )
}

export default SlashedTokens
