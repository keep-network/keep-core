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
        <NavLink exact to="/tokens" label="tokens" icon={<Icons.KeepToken />} />
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

const NavLink = ({ label, to, exact, icon }) => {
  const match = useRouteMatch({
    path: to,
    exact,
  })

  return (
    <Link to={to}>
      <li className={`text-label ${match ? "active-page-link" : ""}`}>
        {icon}
        <span className="ml-1">{label}</span>
        <Icons.ArrowRight />
      </li>
    </Link>
  )
}
