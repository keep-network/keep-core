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
        <Banner.Title className="h3 text-white">
          <p className="mb-0">Deposit KEEP in the coverage pool to</p>
          <p className="mb-0">secure the network and earn rewards</p>
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
