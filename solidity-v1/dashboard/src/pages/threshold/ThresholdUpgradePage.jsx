import React, { useEffect } from "react"
import * as Icons from "../../components/Icons"
import NavLink from "../../components/NavLink"
import TokenAmount from "../../components/TokenAmount"
import { KEEP } from "../../utils/token.utils"
import AllocationProgressBar from "../../components/threshold/AllocationProgressBar"
import UpgradeTokensTile from "../../components/threshold/UpgradeTokensTile"
import resourceTooltipProps from "../../constants/tooltips"
import useKeepBalanceInfo from "../../hooks/useKeepBalanceInfo"
import useGrantedBalanceInfo from "../../hooks/useGrantedBalanceInfo"
import { lte } from "../../utils/arithmetics.utils"
import {
  useWeb3Address,
  useWeb3Context,
} from "../../components/WithWeb3Context"
import { useDispatch } from "react-redux"

const ThresholdUpgradePage = () => {
  const { isConnected } = useWeb3Context()
  const address = useWeb3Address()
  const dispatch = useDispatch()

  useEffect(() => {
    if (isConnected) {
      dispatch({
        type: "staking/fetch_delegations_request",
        payload: { address },
      })
      dispatch({
        type: "token-grant/fetch_grants_request",
        payload: { address },
      })
    }
  }, [dispatch, isConnected, address])

  const { totalOwnedUnstakedBalance, totalKeepTokenBalance } =
    useKeepBalanceInfo()

  const { totalGrantedUnstakedBalance, totalGrantedTokenBalance } =
    useGrantedBalanceInfo()

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
          currentValue={totalOwnedUnstakedBalance}
          totalValue={totalKeepTokenBalance}
          className={"mb-1"}
        />
        <AllocationProgressBar
          title={"available grant allocation"}
          currentValue={totalGrantedUnstakedBalance}
          totalValue={totalGrantedTokenBalance}
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
              amount={totalOwnedUnstakedBalance}
            />
          </UpgradeTokensTile>
          <UpgradeTokensTile
            title={"Grant Allocation"}
            btnText={"withdraw from grant"}
            onBtnClick={onWithdrawFromGrant}
            buttonDisabled={lte(totalGrantedUnstakedBalance, 0)}
            titleTooltipProps={
              resourceTooltipProps.thresholdPageGrantAllocation
            }
          >
            <UpgradeTokensTile.Row
              label={"Available KEEP"}
              amount={totalGrantedUnstakedBalance}
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
