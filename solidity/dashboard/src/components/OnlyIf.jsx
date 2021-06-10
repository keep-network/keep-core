import React from "react"

const OnlyIf = ({ condition, children }) => {
  return condition ? children : <></>
}

export default OnlyIf
