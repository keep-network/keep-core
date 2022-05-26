import React from "react"
import { colors } from "../constants/colors"
import NavLink from "./NavLink"

const wrapperStyle = {
  width: "325px",
  border: `1px solid ${colors.grey60}`,
  padding: "2rem",
}

const ApplicationBox = ({
  icon,
  name,
  websiteUrl,
  websiteName,
  description,
  btnLink,
}) => {
  return (
    <section
      className="flex column center space-between mt-1"
      style={wrapperStyle}
    >
      {icon}
      <h2 className="mt-2">{name}</h2>
      <a
        href={websiteUrl}
        className="h4"
        rel="noopener noreferrer"
        target="_blank"
      >
        {websiteName}
      </a>
      <span className="text-small text-center text-grey-60 mt-3">
        {description}
      </span>
      <NavLink to={btnLink} className="btn btn-primary mt-2">
        manage
      </NavLink>
    </section>
  )
}

export default React.memo(ApplicationBox)

const emptyBoxWrapperStyle = {
  border: `1px solid ${colors.grey20}`,
  padding: "2rem",
  width: "325px",
  minHeight: "425px",
}
export const EmptyApplicationBox = React.memo(() => (
  <section
    className="flex column full-center mt-1"
    style={emptyBoxWrapperStyle}
  >
    <h4 className="text-grey-30 mb-1">Coming Soon</h4>
    <span className="text-grey-30 text-small text-center">
      Future applications will be listed here.
    </span>
  </section>
))
