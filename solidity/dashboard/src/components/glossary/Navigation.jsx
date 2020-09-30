import React from "react"
import { useLocation, useHistory } from "react-router-dom"
import { Link } from "react-scroll"

const Navigation = () => {
  return (
    <nav className="tile glossary-nav">
      <h5 className="text-grey-50">Contents</h5>
      <ul>
        <GlossaryLink
          to={{ pathname: "/glossary", hash: "#documentation" }}
          label="Documentation"
        />
        <GlossaryLink
          to={{ pathname: "/glossary", hash: "#quick-terminology" }}
          label="Quick Terminology"
        />
        <GlossaryLink
          to={{ pathname: "/glossary", hash: "#diagram" }}
          label="Delegation Diagram"
        />
      </ul>
    </nav>
  )
}

const GlossaryLink = ({ label, to }) => {
  const { hash } = useLocation()
  const history = useHistory()

  return (
    <li>
      <Link
        containerId="main-content-wrapper"
        className={`text-small${hash === to.hash ? " active" : ""}`}
        activeClass="active"
        to={to.hash.slice(1)}
        spy={true}
        smooth={true}
        offset={-200}
        duration={500}
        onSetActive={() => history.replace(to)}
      >
        {label}
      </Link>
    </li>
  )
}

export default Navigation
