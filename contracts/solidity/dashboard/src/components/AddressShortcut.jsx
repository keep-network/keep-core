import React from 'react'
import { shortenAddress } from '../utils'

const AddressShortcut = ({ address, classNames }) => {
  return (
    <span className={`address-shortcut ${classNames}`}>
      { shortenAddress(address) }
    </span>
  )
}

export default React.memo(AddressShortcut)
