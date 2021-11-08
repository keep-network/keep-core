import React, { useMemo } from "react"
import TokenAmount from "../TokenAmount"
import Button from "../Button"
import OnlyIf from "../OnlyIf"
import { LINK } from "../../constants/constants"
import ResourceTooltip from "../ResourceTooltip"

const UpgradeTokensTile = ({
  title,
  className = "",
  btnText,
  onBtnClick,
  buttonDisabled = false,
  isLink = false,
  titleTooltipProps = null,
  children,
}) => {
  return (
    <div className={`upgrade-tokens-tile ${className}`}>
      <div className="upgrade-tokens-tile__title">
        <span>{title}</span>
        <OnlyIf condition={titleTooltipProps}>
          <ResourceTooltip tooltipClassName="ml-1" {...titleTooltipProps} />
        </OnlyIf>
      </div>
      <div>{children}</div>
      <OnlyIf condition={!isLink}>
        <Button
          className="btn btn-primary btn-md upgrade-tokens-tile__button"
          onClick={onBtnClick}
          disabled={buttonDisabled}
        >
          {btnText}
        </Button>
      </OnlyIf>
      <OnlyIf condition={isLink}>
        <a
          href={LINK.tbtcDapp}
          rel="noopener noreferrer"
          target="_blank"
          className="btn btn-primary btn-md upgrade-tokens-tile__button"
        >
          {btnText} ↗
        </a>
      </OnlyIf>
    </div>
  )
}

const UpgradeTokenFileRow = ({ label, amount, isDataFetching = false }) => {
  const renderTokenAmount = useMemo(() => {
    if (isDataFetching) {
      return <span>-- KEEP</span>
    }

    return (
      <TokenAmount
        amount={amount}
        symbolClassName="upgrade-token-tile-row__token-amount-symbol"
        amountClassName="upgrade-token-tile-row__token-amount-amount"
      />
    )
  }, [isDataFetching, amount])

  return (
    <div className={"upgrade-tokens-tile-row"}>
      <span className="text-grey-50">{label}</span>
      {renderTokenAmount}
    </div>
  )
}

UpgradeTokensTile.Row = UpgradeTokenFileRow

export default UpgradeTokensTile
