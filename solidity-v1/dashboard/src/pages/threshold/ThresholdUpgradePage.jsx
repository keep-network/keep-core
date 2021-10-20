import React from "react"
import * as Icons from "../../components/Icons"
import NavLink from "../../components/NavLink"

const ThresholdUpgradePage = () => {
  return (
    <section className="threshold-upgrade-page">
      <section className="tile threshold-upgrade-page__explanation">
        <Icons.CoveragePool className="threshold-upgrade-page__explanation__icon" />
        <header>
          <h2 className="text-grey-70">Upgrade Your KEEP to T</h2>
          <h3 className="text-grey-50">
            Threshold Network is the network merger between Keep and NuCypher.
            Upgrade your KEEP to T below.
          </h3>
        </header>
        <NavLink
          to="/threshold/how-it-works"
          className="btn btn-secondary btn-md explanation__upgrade-btn"
        >
          learn more
        </NavLink>
      </section>

      <section className="tile threshold-upgrade-page__not-staked">
        <div>Not staked</div>
      </section>

      <section className="tile threshold-upgrade-page__staked">
        <div>Staked</div>
      </section>
    </section>
  )
}

ThresholdUpgradePage.route = {
  title: "Threshold Upgrade",
  path: "/threshold/upgrade",
  exact: true,
}

export default ThresholdUpgradePage
