import React from 'react'
import { formatDate, displayAmount } from '../utils/general.utils'
import AddressShortcut from './AddressShortcut'

const SlashedTokensList = ({ slashedTokens }) => {
  return (
    <div className="slashed-tokens-list">
      <div className="flex flex-1">
        <span className="text-label flex flex-1">amount</span>
        <span className="text-label flex-2">details</span>
      </div>
      <ul className="flex column flex-1">
        {slashedTokens.map(renderSlashedTokenItem)}
      </ul>
    </div>
  )
}

const renderSlashedTokenItem = (item) => <SlashedTokeItem key={item.id} {...item} />

const SlashedTokeItem = React.memo(({ amount, date, event, groupPublicKey }) => (
  <li className="flex row flex-1" >
    <div className="text-big flex-1">
      <span className="text-error">{amount && `-${displayAmount(amount)} `}</span>
      <span className="text-grey-40">KEEP</span>
    </div>
    <div className="details flex-2">
      <div className="text-big text-grey-70">
        Group <AddressShortcut address={groupPublicKey} classNames="text-big text-grey-70" />&nbsp;
        {event === 'UnauthorizedSigningReported' ?
          'key was leaked. Private key was published outside of the members of the signing group.' :
          'was selected to do work and not enough members participated.'
        }
      </div>
      <div className="text-small text-grey-50">
        {formatDate(date)}
      </div>
    </div>
  </li>
))


export default SlashedTokensList
