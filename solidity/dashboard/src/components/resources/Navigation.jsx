import React from "react"
import { useLocation, useHistory } from "react-router-dom"
import { Link } from "react-scroll"

const Navigation = () => {
  return (
    <nav className="tile resources-nav">
      <h5 className="text-grey-50">Contents</h5>
      <ul>
        <ResourcesLink
          to={{ pathname: "/resources", hash: "#documentation" }}
          label="Documentation"
        />
        <ResourcesLink
          to={{ pathname: "/resources", hash: "#quick-terminology" }}
          label="Quick Terminology"
        />
        <ResourcesLink
          to={{ pathname: "/resources", hash: "#diagram" }}
          label="Delegation Diagram"
          offset={-800}
        />
      </ul>
    </nav>
  )
}

const ResourcesLink = ({ label, to, offset = -100 }) => {
  const { hash } = useLocation()
  const history = useHistory()

  return (
    <li>
      <Link
        className={`text-small${
          hash === to.hash ? " active" : ""
        } cursor-pointer`}
        activeClass="active"
        to={to.hash.slice(1)}
        spy={true}
        smooth={true}
        offset={offset}
        duration={500}
        onSetActive={() => history.replace(to)}
      >
        {label}
      </Link>
    </li>
  )
}

export default Navigation
