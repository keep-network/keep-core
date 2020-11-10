import React from "react"
import { Link, NavLink, useRouteMatch } from "react-router-dom"
import * as Icons from "./Icons"
import { isEmptyArray } from "../utils/array.utils"
import Divider from "./Divider"

const styles = {
  overviewDivider: { margin: "1rem 1.5rem" },
}

export const SideMenu = (props) => {
  return (
    <nav className="side-menu--active">
      <ul className="side-menu__list">
        <li className="side-menu__route-wrapper">
          <NavLink
            to="/overview"
            className="side-menu__route"
            activeClassName="side-menu__route--active"
          >
            Overview
          </NavLink>
          <Divider style={styles.overviewDivider} />
        </li>
        <NavLinkSection
          label="stake"
          icon={<Icons.StakeDrop />}
          subroutes={[
            {
              label: "Delegation",
              path: "/delegation",
              exact: false,
            },
            {
              label: "Token Grants",
              path: "/token-grants",
              exact: true,
            },
          ]}
        />
        <NavLinkSection
          label="work"
          icon={<Icons.SwordOperations />}
          subroutes={[
            {
              label: "Applications",
              path: "/applications",
              exact: false,
            },
            { label: "Operations", path: "/operations", exact: true },
          ]}
        />
        <NavLinkSection
          label="earn"
          icon={<Icons.FeesVector />}
          subroutes={[
            { label: "Earnings", path: "/earnings", exact: false },
            { label: "Reweards", path: "rewards", exact: false },
          ]}
        />
        <NavLinkSection
          label="help"
          icon={<Icons.Question />}
          subroutes={[
            { label: "FAQ", path: "/faq", exact: "true" },
            { label: "Resources", path: "/resources", exact: "false" },
          ]}
        />

        {/* TODO: display this link if a user is a keep token contract deployer. This is only used in development mode*/}
        {/* {isKeepTokenContractDeployer && (
          <NavLink exact to="/create-token-grants" label="token grants" />
        )} */}
      </ul>
    </nav>
  )
}

const NavLinkSection = ({ label, icon, subroutes = [] }) => {
  return (
    <li className="side-menu__section">
      <div className="side-menu__section__content">
        <div className="side-menu__section__content__icon">{icon}</div>
        <span className="side-menu__section__content__title">{label}</span>
      </div>
      {!isEmptyArray(subroutes) && (
        <ul className="side-menu__section__routes">
          {subroutes.map(renderRoute)}
        </ul>
      )}
    </li>
  )
}

const renderRoute = (route) => (
  <NavLinkSectionRoute key={route.path} {...route} />
)

const NavLinkSectionRoute = ({ label, path, exact }) => {
  const match = useRouteMatch({
    path,
    exact,
  })

  return (
    <li className="side-menu__route-wrapper">
      <Link to={path} className={`side-menu__route${match ? "--active" : ""}`}>
        {label}
      </Link>
    </li>
  )
}
