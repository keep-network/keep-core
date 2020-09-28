import React, { useState, useContext } from "react"
import { Link, useRouteMatch } from "react-router-dom"
import { ContractsDataContext } from "./ContractsDataContextProvider"
import Footer from "./Footer"
import * as Icons from "./Icons"

export const SideMenuContext = React.createContext({})

export const SideMenuProvider = (props) => {
  const [isOpen, setIsOpen] = useState(false)

  const toggle = () => {
    setIsOpen(!isOpen)
  }

  return (
    <SideMenuContext.Provider value={{ isOpen, toggle }}>
      {props.children}
    </SideMenuContext.Provider>
  )
}

export const SideMenu = (props) => {
  const { isOpen } = useContext(SideMenuContext)
  const { isKeepTokenContractDeployer } = useContext(ContractsDataContext)

  return (
    <nav className={`side-menu${isOpen ? "--active " : ""}`}>
      <ul className="side-menu__list">
        <NavLink to="/tokens" label="tokens" icon={<Icons.KeepToken />} />
        <NavLink
          exact
          to="/operations"
          label="operations"
          icon={<Icons.Operations />}
        />
        <NavLink to="/rewards" label="rewards" icon={<Icons.Rewards />} />
        <NavLink
          to="/applications"
          label="applications"
          icon={<Icons.Authorizer />}
        />
        <NavLink
          exact
          to="/glossary"
          label="glossary"
          icon={<Icons.Glossary />}
          wrapperClassName="glossary-link text-label"
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
