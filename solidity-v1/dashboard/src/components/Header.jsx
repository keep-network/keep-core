import React from "react"
import { isEmptyArray } from "../utils/array.utils"
import { Web3Status } from "./Web3Status"
import Chip from "./Chip"
import NavLink from "./NavLink"
import OnlyIf from "./OnlyIf"

const Header = ({ title, subLinks, className = "", newPage = false }) => {
  return (
    <header className={`header ${className}`}>
      <div className="header__content">
        <h1 className="header__title">
          {title}{" "}
          {newPage && <Chip text="NEW" className={"header__chip ml-1"} />}
        </h1>
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

const SubNavItem = ({ title, path, withNewLabel }) => {
  return (
    <li className="sub-nav__item-wrapper">
      <NavLink
        to={path}
        className="sub-nav__item"
        activeClassName="sub-nav__item--active"
        exact={true}
      >
        {title}{" "}
        <OnlyIf condition={withNewLabel}>
          <span
            style={{
              height: "10px",
              width: "10px",
              backgroundColor: "#48DBB4",
              borderRadius: "50%",
              display: "inline-block",
            }}
          />
        </OnlyIf>
      </NavLink>
    </li>
  )
}

const renderSubNavItem = (item, index) => (
  <SubNavItem key={`${index}-${item.path}`} {...item} />
)

export default Header
