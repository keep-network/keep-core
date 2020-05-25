import React from "react"
import { Link, useRouteMatch } from "react-router-dom"

const Navigation = () => {
  return (
    <nav className="tile glossary-nav">
      <h5>Content</h5>
      <ul>
        <GlossaryLink to={"/glossary#documentation"} label="Documentation" />
        <GlossaryLink
          to={"/glossary#documentation"}
          label="Quick Terminology"
        />
        <GlossaryLink
          to={"/glossary#documentation"}
          label="Delegation Diagram"
        />
      </ul>
    </nav>
  )
}

const GlossaryLink = ({ label, to, exact }) => {
  const match = useRouteMatch({
    path: to,
    exact,
  })

  return (
    <li className="">
      <Link to={to} exact className="text-small">
        {label}
      </Link>
    </li>
  )
}

export default Navigation
