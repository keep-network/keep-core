import React from "react"
import { NotFound404 } from "./NotFound404"

const withAddress = (Component) => (props) => {
  const address = window.location.pathname.split("/")[1]
  // TODO: Use/Create function to check if it's valid eth address
  if (address.slice(0, 2) === "0x") {
    return <Component {...props} />
  }

  return <NotFound404 />
}

export default withAddress
