import React from 'react'
import { shortenAddress } from '../utils'

const AddressShortcut = ({ address }) => {
  return (
    <span className='address-shortcut'>
      { shortenAddress(address) }
    </span>
  )
}

export default React.memo(AddressShortcut)
