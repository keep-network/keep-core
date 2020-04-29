import React, { useRef, useState, useEffect } from 'react'
import { shortenAddress } from '../utils/general.utils'


const AddressShortcut = ({ address, classNames }) => {
  const addressElement = useRef(null)
  const [copyStatus, setCopyStatus] = useState('Copy to clipboard')

  useEffect(() => {
    const copyEventListener = (event) => {
      event.preventDefault()
      if (event.clipboardData) {
        event.clipboardData.setData('text/plain', address)
      } else if (window.clipboardData) {
        window.clipboardData.setData('Text', address)
      } else {
        setCopyStatus(`Cannot copy value: ${address}!`)
      }
    }
    if (addressElement.current !== null) {
      addressElement.current.addEventListener('copy', copyEventListener)
      return () => {
        addressElement.current.removeEventListener('copy', copyEventListener)
      }
    }
  })

  const copyToClipboard = () => {
    try {
      if (document.selection) {
        const range = document.body.createTextRange()
        range.moveToElementText(addressElement.current)
        range.select().createTextRange()
      } else if (window.getSelection) {
        const range = document.createRange()
        range.selectNode(addressElement.current)
        window.getSelection().addRange(range)
      }
      document.execCommand('copy')
      setCopyStatus('Copied!')
    } catch (error) {
      setCopyStatus(`Cannot copy value: ${address}!`)
    }
  }

  return (
    <span
      onClick={copyToClipboard}
      onMouseOut={() => setCopyStatus('Copy to clipboard')}
      className={`address-shortcut tooltip address ${classNames}`}
    >
      <span className="tooltip-text address bottom">
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
