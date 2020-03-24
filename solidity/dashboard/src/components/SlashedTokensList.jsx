import React from 'react'
import { formatDate, displayAmount } from '../utils'
import AddressShortcut from './AddressShortcut'

const SlashedTokensList = ({ slashedTokens }) => {
  return (
    <div className="slashed-tokens-list">
      <div className="flex flex-1">
        <span className="text-label flex flex-1">AMOUNT</span>
        <span className="text-label flex-2">SLASH EXPLANATION</span>
      </div>
      <ul className="flex flex-column flex-1">
        {slashedTokens.map(renderSlashedTokenItem)}
      </ul>
    </div>
  )
}

const renderSlashedTokenItem = (item) => <SlashedTokeItem key={item.id} {...item} />

const SlashedTokeItem = React.memo(({ amount, date, typeOfPunishment, groupPublicKey }) => (
  <li className="flex flex-row flex-1" >
    <div className="text-big text-dark-red flex-1">
      - {amount && `${displayAmount(amount)}`} KEEP
    </div>
    <div className="details flex-2">
      <div className="text-big text-darker-grey">
        Group <AddressShortcut address={groupPublicKey} classNames="text-big text-darker-grey" />&nbsp;
        {typeOfPunishment === 0 ?
          'was selected to do work and not enough members participated.' :
          'key was leaked. Private key was published outside of the members of the signing group.'
        }
      </div>
      <div className="text-small text-grey">
        {formatDate(date)}
      </div>
    </div>
  </li>
))


export default SlashedTokensList
