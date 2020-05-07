import React, { useState, useContext } from "react"
import { Link, useRouteMatch } from "react-router-dom"
import { Web3Status } from "./Web3Status"
import { Web3Context } from "./WithWeb3Context"
import { ContractsDataContext } from "./ContractsDataContextProvider"
import AddressShortcut from "./AddressShortcut"
import { NetworkStatus } from "./NetworkStatus"
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
  const { yourAddress } = useContext(Web3Context)
  const { isKeepTokenContractDeployer } = useContext(ContractsDataContext)

  return (
    <nav className={`${isOpen ? "active " : ""}side-menu`}>
      <ul>
        <NavLink
          to="/tokens"
          label="tokens"
          icon={<Icons.KeepToken />}
          sublinks={[
            { to: "/tokens/delegate", exact: true, label: "Delegate Tokens" },
            { to: "/tokens/grants", exact: true, label: "Token Grants" },
          ]}
        />
        <NavLink
          exact
          to="/operations"
          label="operations"
          icon={<Icons.Operations />}
        />
        <NavLink exact to="/rewards" label="rewards" icon={<Icons.Rewards />} />
        <NavLink
          exact
          to="/authorizer"
          label="authorizer"
          icon={<Icons.Authorizer />}
        />
        {isKeepTokenContractDeployer && (
          <NavLink exact to="/create-token-grants" label="token grants" />
        )}
        <Web3Status />
        <div className="account-address">
          <h5 className="text-grey-50">
            <span>ADDRESS&nbsp;</span>
            <AddressShortcut classNames="text-small" address={yourAddress} />
          </h5>
          <NetworkStatus />
        </div>
      </ul>
    </nav>
  )
}

const NavLink = ({
  label,
  to,
  exact,
  icon,
  sublinks,
  wrapperClassName,
  activeClassName,
  withArrowRight,
}) => {
  const match = useRouteMatch({
    path: to,
    exact,
  })

  return (
    <li className={`${wrapperClassName} ${match ? activeClassName : ""}`}>
      <Link to={to}>
        {icon}
        <span className="ml-1">{label}</span>
        {withArrowRight && <Icons.ArrowRight />}
      </Link>
      <SubNavLinks sublinks={sublinks} />
    </li>
  )
}

NavLink.defaultProps = {
  wrapperClassName: "text-label",
  activeClassName: "active-page-link",
  withArrowRight: true,
}

const SubNavLinks = ({ sublinks }) => {
  if (!sublinks) return null

  return (
    <ul className="sublinks">
      {sublinks.map((sublink) => (
        <NavLink
          key={sublink.label}
          {...sublink}
          wrapperClassName="sublink"
          withArrowRight={false}
        />
      ))}
    </ul>
  )
}
