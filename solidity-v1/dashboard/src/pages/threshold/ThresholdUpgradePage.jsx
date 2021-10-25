import React from "react"
import * as Icons from "../../components/Icons"
import NavLink from "../../components/NavLink"
import TokenAmount from "../../components/TokenAmount"
import { KEEP } from "../../utils/token.utils"
import AllocationProgressBar from "../../components/threshold/AllocationProgressBar"
import UpgradeTokensTile from "../../components/threshold/UpgradeTokensTile"
import resourceTooltipProps from "../../constants/tooltips"

const ThresholdUpgradePage = () => {
  const onWithdrawFromGrant = () => {
    console.log("withdraw from grant clicked!")
  }

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

      <section className="tile not-staked">
        <h3 className="mb-1">Not staked</h3>
        <TokenAmount
          wrapperClassName={"not-staked__token-amount mb-2"}
          amount={"500000320000000000000000"}
          token={KEEP}
          withIcon
        />
        <AllocationProgressBar
          title={"wallet"}
          currentValue={80}
          totalValue={100}
          className={"mb-1"}
        />
        <AllocationProgressBar
          title={"available grant allocation"}
          currentValue={20.33453453}
          totalValue={103.342324}
          className={"mb-2"}
        />
        <div className="upgrade-not-staked">
          <h4 className={"mb-1"}>Upgrade not staked tokens</h4>
          <UpgradeTokensTile
            title={"Wallet"}
            btnText={"upgrade to t"}
            className={"mb-1"}
            isLink
          >
            <UpgradeTokensTile.Row
              label={"Liquid KEEP"}
              amount={1000000000000000000000000}
            />
          </UpgradeTokensTile>
          <UpgradeTokensTile
            title={"Grant Allocation"}
            btnText={"withdraw from grant"}
            onBtnClick={onWithdrawFromGrant}
            titleTooltipProps={
              resourceTooltipProps.thresholdPageGrantAllocation
            }
          >
            <UpgradeTokensTile.Row
              label={"Available KEEP"}
              amount={1000000000000000000000000}
            />
          </UpgradeTokensTile>
        </div>
      </section>

      <section className="tile staked">
        <h3>Staked</h3>
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
