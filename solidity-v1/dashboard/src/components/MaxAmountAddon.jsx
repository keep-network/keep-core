import React from "react"
import * as Icons from "./Icons"
import Tag from "./Tag"

const MaxAmountAddon = ({ onClick, text, ...otherProps }) => {
  return (
    <Tag
      IconComponent={Icons.Plus}
      text={text}
      onClick={onClick}
      {...otherProps}
    />
  )
}

export default MaxAmountAddon
