import React, { useEffect, useMemo } from "react"
import * as Icons from "../../components/Icons"
import NavLink from "../../components/NavLink"
import TokenAmount from "../../components/TokenAmount"
import { KEEP } from "../../utils/token.utils"
import AllocationProgressBar from "../../components/threshold/AllocationProgressBar"
import UpgradeTokensTile from "../../components/threshold/UpgradeTokensTile"
import resourceTooltipProps from "../../constants/tooltips"
import useKeepBalanceInfo from "../../hooks/useKeepBalanceInfo"
import useGrantedBalanceInfo from "../../hooks/useGrantedBalanceInfo"
import { add, lte } from "../../utils/arithmetics.utils"
import {
  useWeb3Address,
  useWeb3Context,
} from "../../components/WithWeb3Context"
import { useDispatch, useSelector } from "react-redux"
import { useModal } from "../../hooks/useModal"
import { LINK, MODAL_TYPES } from "../../constants/constants"
import TokenAmountSkeleton from "../../components/skeletons/TokenAmountSkeleton"
import ResourceTooltip from "../../components/ResourceTooltip"
import useUpdatePendingUndelegations from "../../hooks/useUpdatePendingUndelegations"
import useUpdateInitializedDelegations from "../../hooks/useUpdateInitializedDelegations"
import HowItWorksPage from "./HowItWorksPage"

const ThresholdUpgradePage = () => {
  const { isConnected } = useWeb3Context()
  const address = useWeb3Address()
  const dispatch = useDispatch()
  const { openModal } = useModal()

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

  const { delegations, undelegations, isDelegationDataFetching } = useSelector(
    (state) => state.staking
  )

  useUpdateInitializedDelegations(delegations)
  useUpdatePendingUndelegations(undelegations)

  const { grants, isFetching: isGrantDataFetching } = useSelector(
    (state) => state.tokenGrants
  )

  const { value: keepBalance, isFetching: isKeepBalanceFetching } = useSelector(
    (state) => state.keepTokenBalance
  )

  const isDataFetching = useMemo(() => {
    return (
      isDelegationDataFetching || isGrantDataFetching || isKeepBalanceFetching
    )
  }, [isDelegationDataFetching, isGrantDataFetching, isKeepBalanceFetching])

  const { totalOwnedStakedBalance } = useKeepBalanceInfo()

  const { totalGrantedStakedBalance } = useGrantedBalanceInfo()

  const totalGrantedReadyToReleaseTokens = useMemo(() => {
    if (isDataFetching || grants.length === 0) return "0"
    return grants
      .map(({ readyToRelease }) => readyToRelease)
      .reduce(add, "0")
      .toString()
  }, [grants, isDataFetching])

  const notStakedTotalAmount = useMemo(() => {
    if (isDataFetching) return "0"
    return add(keepBalance, totalGrantedReadyToReleaseTokens).toString()
  }, [keepBalance, totalGrantedReadyToReleaseTokens, isDataFetching])

  const stakedTotalAmount = useMemo(() => {
    return add(totalOwnedStakedBalance, totalGrantedStakedBalance).toString()
  }, [totalOwnedStakedBalance, totalGrantedStakedBalance])

  const {
    totalStakedPendingKeep,
    totalStakedAvailableKeep,
    totalUndelegatedAvailableKeep,
  } = useMemo(() => {
    const totalStakedPendingKeep = [...undelegations]
      .filter((delegation) => !delegation.canRecoverStake)
      .map(({ amount }) => amount)
      .reduce(add, "0")
      .toString()

    const totalStakedAvailableKeep = [...delegations]
      .map(({ amount }) => amount)
      .reduce(add, "0")
      .toString()

    const totalUndelegatedAvailableKeep = [...undelegations]
      .filter((delegation) => delegation.canRecoverStake)
      .map(({ amount }) => amount)
      .reduce(add, "0")
      .toString()

    return {
      totalStakedPendingKeep,
      totalStakedAvailableKeep,
      totalUndelegatedAvailableKeep,
    }
  }, [delegations, undelegations])

  const onWithdrawFromGrant = () => {
    openModal(MODAL_TYPES.WithdrawGrantedTokens)
  }

  return (
    <section className="threshold-upgrade-page">
      <section className="tile threshold-upgrade-page__explanation">
        <Icons.KeepTUpgrade className="threshold-upgrade-page__explanation__icon" />
        <header>
          <h2 className="text-grey-70 threshold-how-it-works-page__explanation-title">
            Upgrade Your KEEP to T
          </h2>
          <h3 className="text-grey-50 threshold-how-it-works-page__explanation-description">
            Threshold Network is the network merger between Keep and NuCypher.
            Upgrade your KEEP to T below.
          </h3>
        </header>
        <NavLink
          to={HowItWorksPage.route.path}
          className="btn btn-secondary btn-md explanation__upgrade-btn"
        >
          learn more
        </NavLink>
      </section>

      <section className="tile not-staked">
        <h3 className="mb-1">Not staked</h3>
        {isDataFetching ? (
          <TokenAmountSkeleton
            wrapperClassName="not-staked__token-amount staked__token-amount--loading"
            textStyles={{ width: "40%", marginLeft: ".3rem" }}
            textClassName="h2"
            icon={Icons.KeepOutline}
            iconHeight={32}
            iconWidth={32}
            iconClassName="token-amount__icon keep-outline keep-outline--mint-80"
          />
        ) : (
          <TokenAmount
            wrapperClassName={"not-staked__token-amount mb-2"}
            amount={notStakedTotalAmount}
            token={KEEP}
            withIcon
          />
        )}
        <AllocationProgressBar
          title={"wallet"}
          currentValue={keepBalance}
          currentValueLabel={"unstaked wallet balance"}
          totalValue={notStakedTotalAmount}
          className={"mb-1"}
          isDataFetching={isDataFetching}
        />
        <AllocationProgressBar
          title={"available grant allocation"}
          currentValue={totalGrantedReadyToReleaseTokens}
          currentValueLabel={"unstaked grant balance"}
          totalValue={notStakedTotalAmount}
          className={"mb-2"}
          isDataFetching={isDataFetching}
        />
        <div className="upgrade-not-staked">
          <h4 className={"mb-1"}>Upgrade Not Staked Tokens</h4>
          <UpgradeTokensTile
            title={"Wallet"}
            btnText={"upgrade to t"}
            className={"mb-1"}
            isLink
            buttonDisabled={isDataFetching}
            renderButton={() => (
              <UpgradeTokensTile.Link
                to={LINK.thresholdDapp}
                linkText={"upgrade to t"}
              />
            )}
          >
            <UpgradeTokensTile.Row
              label={"Liquid KEEP"}
              amount={keepBalance}
              isDataFetching={isDataFetching}
            />
          </UpgradeTokensTile>
          <UpgradeTokensTile
            title={"Grant Allocation"}
            tooltipComponent={ResourceTooltip}
            titleTooltipProps={
              resourceTooltipProps.thresholdPageGrantAllocation
            }
            renderButton={() => (
              <UpgradeTokensTile.Button
                btnText={"withdraw from grant"}
                buttonDisabled={
                  lte(totalGrantedReadyToReleaseTokens, 0) || isDataFetching
                }
                onBtnClick={onWithdrawFromGrant}
              />
            )}
          >
            <UpgradeTokensTile.Row
              label={"Available KEEP"}
              amount={totalGrantedReadyToReleaseTokens}
              isDataFetching={isDataFetching}
            />
          </UpgradeTokensTile>
        </div>
      </section>

      <section className="tile staked">
        <div className="staked__title-container">
          <h3 className="staked__title">Staked</h3>
          {/* <div className="staked__additional-info">*/}
          {/*  <span className="staked__additional-info-row mr-2">*/}
          {/*    <Icons.Success*/}
          {/*      width={16}*/}
          {/*      height={16}*/}
          {/*      className="staked__additional-info-icon staked__additional-info-icon--color-green"*/}
          {/*    />{" "}*/}
          {/*    ECDSA*/}
          {/*  </span>*/}
          {/*  <span className="staked__additional-info-row">*/}
          {/*    <Icons.Success*/}
          {/*      width={16}*/}
          {/*      height={16}*/}
          {/*      className="staked__additional-info-icon staked__additional-info-icon--color-green"*/}
          {/*    />{" "}*/}
          {/*    Random Beacon*/}
          {/*  </span>*/}
          {/* </div>*/}
        </div>
        {isDataFetching ? (
          <TokenAmountSkeleton
            wrapperClassName="staked__token-amount staked__token-amount--loading"
            textStyles={{ width: "40%", marginLeft: ".3rem" }}
            textClassName="h2"
            icon={Icons.KeepOutline}
            iconHeight={32}
            iconWidth={32}
            iconClassName="token-amount__icon keep-outline keep-outline--mint-80"
          />
        ) : (
          <TokenAmount
            wrapperClassName={"staked__token-amount mb-2"}
            amount={stakedTotalAmount}
            token={KEEP}
            withIcon
          />
        )}
        <AllocationProgressBar
          title={"staked"}
          currentValue={totalStakedAvailableKeep}
          totalValue={stakedTotalAmount}
          className={"mb-1"}
          secondaryValue={totalStakedPendingKeep}
          withLegend
          currentValueLabel={"Staked"}
          secondaryValueLabel={"Pending Undelegation"}
          isDataFetching={isDataFetching}
        />
        <AllocationProgressBar
          title={"undelegated"}
          currentValue={totalUndelegatedAvailableKeep}
          currentValueLabel={"Undelegated"}
          totalValue={stakedTotalAmount}
          className={"mb-3"}
          isDataFetching={isDataFetching}
        />
        <div className="upgrade-staked">
          <h4 className={"mb-1"}>Upgrade Staked Tokens</h4>
          <UpgradeTokensTile
            title={"Staked"}
            buttonDisabled={isDataFetching}
            className={"mb-1"}
            renderButton={() => (
              <UpgradeTokensTile.NavLink
                linkText={"undelegate"}
                to={"/overview"}
              />
            )}
          >
            <UpgradeTokensTile.Row
              label={"Total Pending KEEP"}
              amount={totalStakedPendingKeep}
              isDataFetching={isDataFetching}
            />
            <UpgradeTokensTile.Row
              label={"Total Available KEEP"}
              amount={totalStakedAvailableKeep}
              isDataFetching={isDataFetching}
            />
          </UpgradeTokensTile>
          <UpgradeTokensTile
            title={"Undelegated"}
            buttonDisabled={
              lte(totalUndelegatedAvailableKeep, 0) || isDataFetching
            }
            renderButton={() => (
              <UpgradeTokensTile.NavLink
                linkText={"claim tokens"}
                to={"/overview"}
              />
            )}
          >
            <UpgradeTokensTile.Row
              label={"Total Available KEEP"}
              amount={totalUndelegatedAvailableKeep}
              isDataFetching={isDataFetching}
            />
          </UpgradeTokensTile>
        </div>
      </section>
    </section>
  )
}

ThresholdUpgradePage.route = {
  title: "Upgrade",
  path: "/threshold/upgrade",
  exact: true,
}

export default ThresholdUpgradePage
