import React, { useState, useContext } from 'react'
import { Link, useRouteMatch } from 'react-router-dom'
import { Web3Status } from './Web3Status'

export const SideMenuContext = React.createContext({})

export const SiedMenuProvider = (props) => {
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
    return (
        <nav className={`${isOpen ? 'active ' : '' }side-menu`}>
            <ul>
                <NavLink exact to="/overview" label='Overview'/>
                <NavLink exact to="/stake" label='Stake'/>
                <NavLink exact to="/token-grants" label='Token Grants'/>
                <NavLink exact to="/create-token-grants" label='Create Token Grant'/>
                <Web3Status />
                <div className='account-address'>
                    <strong>Account address: </strong>
                    <p className="txt-primary">0xcCFe2E36B3F10152D19dD7d14d651F213c9af4b0</p>
                </div>
            </ul>
        </nav>
    )
}

const NavLink = ({ label, to, exact }) => {
    const match = useRouteMatch({
        path: to,
        exact
    })
    return (
        <li className={ match ? 'active-page-link' : ''}>
            <Link to={to}>{label}</Link>
        </li>
    )
}