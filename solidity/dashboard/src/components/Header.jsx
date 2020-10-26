import React from "react"
import { NavLink } from "react-router-dom"
import { isEmptyArray } from "../utils/array.utils"
import { Web3Status } from "./Web3Status"

const Header = ({ title, subLinks }) => {
  return (
    <header className="header">
      <div className="header__content">
        <h1 className="header__title">{title}</h1>
        <Web3Status />
      </div>
      {!isEmptyArray(subLinks) && (
        <nav className="header__sub-nav">
          <ul className="sub-nav__list">{subLinks.map(renderSubNavItem)}</ul>
        </nav>
      )}
    </header>
  )
}

const SubNavItem = ({ title, path }) => {
  return (
    <li className="sub-nav__item-wrapper">
      <NavLink
        to={path}
        className="sub-nav__item"
        activeClassName="sub-nav__item--active"
        exact={true}
      >
        {title}
      </NavLink>
    </li>
  )
}

const renderSubNavItem = (item, index) => (
  <SubNavItem key={`${index}-${item.path}`} {...item} />
)

export default Header
