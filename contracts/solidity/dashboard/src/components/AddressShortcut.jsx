import React, { useRef, useState } from 'react'
import { shortenAddress } from '../utils/general.utils'


const AddressShortcut = ({ address, classNames }) => {
  const addressElement = useRef(null)
  const [copyStatus, setCopyStatus] = useState('Copy to clipboard')

  const copyToClipboard = () => {
    try {
      if (document.selection) {
        const range = document.body.createTextRange()
        range.moveToElementText(addressElement.current)
        range.select().createTextRange()
        document.execCommand('copy')
      } else if (window.getSelection) {
        const range = document.createRange()
        range.selectNode(addressElement.current)
        window.getSelection().addRange(range)
        document.execCommand('copy')
      }
      setCopyStatus('Copied!')
    } catch (error) {
      setCopyStatus(`Cannot copy value: ${address}!`)
    }
  }

  return (
    <span
      onClick={copyToClipboard}
      onMouseOut={() => setCopyStatus('Copy to clipboard')}
      className={`address-shortcut tooltip ${classNames}`}
    >
      <span className="tooltip-text bottom">
        {copyStatus}
      </span>
      <span
        className="full-address"
        ref={addressElement}
      >
        {address}
      </span>
      { shortenAddress(address) }
    </span>
  )
}

export default React.memo(AddressShortcut)
