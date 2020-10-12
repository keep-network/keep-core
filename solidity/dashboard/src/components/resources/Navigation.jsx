import React from "react"
import { useLocation, useHistory } from "react-router-dom"
import { Link } from "react-scroll"

const Navigation = () => {
  return (
    <nav className="tile resources-nav">
      <h5 className="text-grey-50">Contents</h5>
      <ul>
        <ResourcesLink to="documentation">Documentation</ResourcesLink>
        <ResourcesLink to="diagram">Diagram</ResourcesLink>
        <ResourcesLink to="quick-terminology">Quick Terminology</ResourcesLink>
      </ul>
    </nav>
  )
}

const ResourcesLink = ({ to, children, offset = -100 }) => {
  const { hash } = useLocation()
  const history = useHistory()

  return (
    <li>
      <Link
        className={`text-small${hash === to.hash ? " active" : ""}`}
        activeClass="active"
        to={to}
        spy={true}
        smooth={true}
        offset={offset}
        duration={500}
        onSetActive={() =>
          history.replace({ pathname: "/resources", hash: to })
        }
      >
        {children}
      </Link>
    </li>
  )
}

export default Navigation
