import React from "react"
import Banner from "../Banner"
import * as Icons from "../Icons"
import NavLink from "../NavLink"

const LearnMoreBanner = ({ onClose }) => {
  return (
    <Banner className="banner banner--coverage-pools">
      <Banner.CloseIcon onClick={onClose} />
      <div className="banner__content-wrapper">
        <Banner.Icon icon={Icons.CoveragePool} />
        <Banner.Title className="h3 text-white banner__title--font-weight-600">
          <h3 className="mb-0">Deposit KEEP in the coverage pool to</h3>
          <h3 className="mb-0">secure the network and earn rewards.</h3>
        </Banner.Title>
        <NavLink
          to="/coverage-pools/how-it-works"
          className="btn btn-tertiary btn-lg"
        >
          learn more
        </NavLink>
      </div>
    </Banner>
  )
}

export default LearnMoreBanner
