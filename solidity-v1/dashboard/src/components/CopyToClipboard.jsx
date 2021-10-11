import { useState, useEffect } from "react"
import copy from "copy-to-clipboard"

const CopyToClipboard = (props) => {
  const [copyStatus, setCopyStatus] = useState(props.defaultCopyText)

  useEffect(() => {
    setCopyStatus(props.defaultCopyText)
  }, [props.defaultCopyText])

  const copyToClipboard = () => {
    copy(props.toCopy)
      ? setCopyStatus("Copied!")
      : setCopyStatus(`Cannot copy value: ${props.toCopy}!`)
  }

  const reset = () => {
    setCopyStatus(props.defaultCopyText)
  }

  return props.render({ copyStatus, copyToClipboard, reset })
}

CopyToClipboard.defaultProps = {
  defaultCopyText: "Copy to clipboard",
}

export default CopyToClipboard
