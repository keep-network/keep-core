import React, { useContext } from "react"
import { Link, useRouteMatch } from "react-router-dom"
import { ContractsDataContext } from "./ContractsDataContextProvider"
import Footer from "./Footer"
import * as Icons from "./Icons"

export const SideMenu = (props) => {
  const { isKeepTokenContractDeployer } = useContext(ContractsDataContext)

  return (
    <nav className="side-menu--active">
      <Link to="/">
        <Icons.KeepDashboardLogo className="side-menu__logo" />
      </Link>
      <ul className="side-menu__list">
        <NavLink to="/overview" label="overview" icon={<Icons.Home />} />
        <NavLink to="/tokens" label="tokens" icon={<Icons.FeesVector />} />
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
          to="/glossary"
          label="glossary"
          icon={<Icons.Question />}
        />
        {isKeepTokenContractDeployer && (
          <NavLink exact to="/create-token-grants" label="token grants" />
        )}
      </ul>
      <Footer />
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
