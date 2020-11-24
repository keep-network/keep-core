import React from "react"
import { Link, useRouteMatch } from "react-router-dom"
import * as Icons from "./Icons"

export const SideMenu = (props) => {
  return (
    <nav className="side-menu--active">
      <ul className="side-menu__list">
        <NavLink to="/overview" label="overview" icon={<Icons.Home />} />
        <NavLink
          to="/delegation"
          label="delegation"
          icon={<Icons.FeesVector />}
        />
        <NavLink
          exact
          to="/operations"
          label="operations"
          icon={<Icons.Beacon />}
        />
        <NavLink to="/rewards" label="rewards" icon={<Icons.Rewards />} />
        <NavLink
          to="/applications"
          label="applications"
          icon={<Icons.Authorize />}
        />
        <NavLink
          exact
          to="/resources"
          label="resources"
          icon={<Icons.Question />}
        />
        {/* TODO: display this link if a user is a keep token contract deployer. This is only used in development mode*/}
        {/* {isKeepTokenContractDeployer && (
          <NavLink exact to="/create-token-grants" label="token grants" />
        )} */}
      </ul>
    </nav>
  )
}

const NavLink = ({ label, to, exact, icon }) => {
  const match = useRouteMatch({
    path: to,
    exact,
  })

  return (
    <li className="side-menu__item-wrapper">
      <Link to={to} className={`side-menu__item${match ? "--active" : ""}`}>
        <div className="side-menu__item__icon">{icon}</div>
        <span className="side-menu__item__title">{label}</span>
      </Link>
    </li>
  )
}
