import React from 'react'
import { addressToShortcut } from '../utils'

const AddressShortcut = ({ address }) => {
  return (
    <span className='address-shortcut'>
      { addressToShortcut(address) }
    </span>
  )
}

export default React.memo(AddressShortcut)
