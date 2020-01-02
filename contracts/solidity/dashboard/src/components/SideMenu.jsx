import React, { useState, useContext } from 'react'
import { Link, useRouteMatch } from 'react-router-dom'
import { Web3Status } from './Web3Status'
import { Web3Context } from './WithWeb3Context'
import { ContractsDataContext } from './ContractsDataContextProvider'

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
  const { isTokenHolder, isBeneficiary } = useContext(ContractsDataContext)

  return (
    <nav className={`${ isOpen ? 'active ' : '' }side-menu`}>
      <ul>
        <NavLink exact to="/overview" label='Overview'/>
        { isBeneficiary && <NavLink exact to="/rewards" label='Rewards'/> }
        { isTokenHolder &&
            <>
                <NavLink exact to="/stake" label='Stake'/>
                <NavLink exact to="/token-grants" label='Token Grants'/>
                <NavLink exact to="/create-token-grants" label='Create Token Grant'/>
            </>
        }
        <Web3Status />
        <div className='account-address'>
          <strong>Account address: </strong>
          <p className="txt-primary">{yourAddress}</p>
        </div>
      </ul>
    </nav>
  )
}

const NavLink = ({ label, to, exact }) => {
  const match = useRouteMatch({
    path: to,
    exact,
  })

  return (
    <Link to={to}>
      <li className={ match ? 'active-page-link' : '' }>
        {label}
      </li>
    </Link>
  )
}
