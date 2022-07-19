import React, { useMemo } from "react"
import TokenAmount from "./TokenAmount"
import * as Icons from "./Icons"
import { colors } from "../constants/colors"
import { add, percentageOf } from "../utils/arithmetics.utils"
import Divider from "./Divider"
import ProgressBar from "./ProgressBar"
import Chip from "./Chip"
import ResourceTooltip from "./ResourceTooltip"
import { KEEP } from "../utils/token.utils"
import { Skeleton } from "./skeletons"
import { useWeb3Context } from "./WithWeb3Context"
import resourceTooltipProps from "../constants/tooltips"
import DelegationPage from "../pages/delegation"
import { formatValue } from "../utils/general.utils"
import NavLink from "./NavLink"

const TokenOverview = ({
  totalKeepTokenBalance,
  totalGrantedTokenBalance,
  totalGrantedStakedBalance,
  totalOwnedStakedBalance,
  isFetching,
}) => {
  return (
    <div className="balances-layout">
      <TotalKeepBalance
        walletKeepBalance={totalKeepTokenBalance}
        grantedKeepBalance={totalGrantedTokenBalance}
        isFetching={isFetching}
      />
      <TokenBalance
        type="granted"
        totalBalance={totalGrantedTokenBalance}
        staked={totalGrantedStakedBalance}
        icon={Icons.Grant}
        isFetching={isFetching}
        tooltipProps={resourceTooltipProps.tokenGrant}
      />
      <TokenBalance
        type="wallet"
        totalBalance={totalKeepTokenBalance}
        staked={totalOwnedStakedBalance}
        isFetching={isFetching}
      />
    </div>
  )
}

const TotalKeepBalance = ({
  walletKeepBalance,
  grantedKeepBalance,
  isFetching,
}) => {
  const { isConnected } = useWeb3Context()

  const totalKeep = useMemo(() => {
    return add(walletKeepBalance, grantedKeepBalance)
  }, [walletKeepBalance, grantedKeepBalance])

  return (
    <section className="balance__overview">
      <h3 className="mb-1">Total KEEP Balance</h3>
      <TokenAmount amount={totalKeep} withIcon withMetricSuffix />
      <div className="balance__overview__granted mt-3">
        <h4>Granted Tokens</h4>
        {isFetching ? (
          <Skeleton
            shining
            tag="h4"
            color="grey-20"
            className="ml-a"
            width="25%"
          />
        ) : (
          <h4 className="ml-a">
            {isConnected
              ? `${KEEP.displayAmountWithSymbol(grantedKeepBalance)}`
              : "No data to display"}
          </h4>
        )}
      </div>
      <Divider className="balance__overview__divider" />
      <div className="balance__overview__wallet">
        <h4>Wallet Tokens</h4>
        {isFetching ? (
          <Skeleton
            shining
            tag="h4"
            color="grey-20"
            className="ml-a"
            width="30%"
          />
        ) : (
          <h4 className="ml-a">
            {isConnected
              ? `${KEEP.displayAmountWithSymbol(walletKeepBalance)}`
              : "No data to display"}
          </h4>
        )}
      </div>
    </section>
  )
}

const TokenBalance = ({
  tooltipProps,
  totalBalance = 0,
  staked = 0,
  type = "wallet",
  icon: IconComponent = Icons.Wallet,
  isFetching = false,
}) => {
  const { isConnected } = useWeb3Context()

  const inPercentage = useMemo(() => {
    return formatValue(percentageOf(staked, totalBalance))
  }, [staked, totalBalance])

  const renderAmount = () => {
    if (!isConnected) {
      return "No Data"
    } else if (isConnected && isFetching) {
      return <Skeleton shining tag="h4" color="grey-20" width="80%" />
    } else if (isConnected && !isFetching) {
      return (
        <h4 className="text-grey-40">
          {totalBalance ? (
            <TokenAmount
              amount={totalBalance}
              withMetricSuffix
              amountClassName="text-mint-100"
              symbolClassName="text-mint-100"
            />
          ) : (
            "No data"
          )}
        </h4>
      )
    }
  }

  return (
    <section className={`balance__${type}`}>
      <header className="flex row center mb-2">
        <Chip
          icon={<IconComponent />}
          color="disabled"
          className={`balance__${type}__chip`}
        />
        <div className="flex-1">
          <h3 className="text-grey-70 flex row center">
            {type}&nbsp;
            {tooltipProps && <ResourceTooltip {...tooltipProps} />}
          </h3>
          {renderAmount()}
        </div>
      </header>
      <ProgressBar total={totalBalance} bgColor={colors.grey20}>
        <ProgressBar.Inline
          height={8}
          className={`balance__${type}__progressbar`}
        >
          <ProgressBar.InlineItem value={staked} color={colors.primary} />
        </ProgressBar.Inline>
      </ProgressBar>
      <span className="text-small text-grey-40">
        {isConnected ? inPercentage : "-"}% Staked
      </span>
      <NavLink
        to={`${DelegationPage.route.path}/${type}`}
        className="btn btn-secondary btn-lg mt-2"
        style={{ width: "100%" }}
      >
        stake
      </NavLink>
    </section>
  )
}

export default React.memo(TokenOverview)
