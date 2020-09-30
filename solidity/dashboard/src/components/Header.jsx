import React from "react"
import { NavLink } from "react-router-dom"

const Header = (props) => {
  return (
    <header className="header">
      <div className="header__content">
        <h1 className="header__title">Earnings</h1>
        {/* TODO add a web3 status component */}
      </div>
      <nav className="header__sub-nav">
        <ul className="sub-nav__list">
          <li className="sub-nav__item-wrapper">
            <NavLink
              to="/random-beacon"
              className="sub-nav__item"
              activeClassName="sub-nav__item--active"
            >
              Random Beacon
            </NavLink>
          </li>
          <li className="sub-nav__item-wrapper">
            <NavLink
              to="/random-beacon"
              className="sub-nav__item"
              activeClassName="sub-nav__item--active"
            >
              tBtc
            </NavLink>
          </li>
        </ul>
      </nav>
    </header>
  )
}

export default Header
