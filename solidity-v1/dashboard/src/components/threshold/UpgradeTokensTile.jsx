import React, { useMemo } from "react"
import TokenAmount from "../TokenAmount"
import Button from "../Button"
import OnlyIf from "../OnlyIf"
import ResourceTooltip from "../ResourceTooltip"
import NavLink from "../NavLink"

const UpgradeTokensTile = ({
  title,
  className = "",
  renderButton = () => <UpgradeTokensTile.Button btnText={"button"} />,
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
      {renderButton()}
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

UpgradeTokensTile.Button = ({
  btnText,
  onBtnClick,
  buttonDisabled = false,
  className = "",
}) => {
  return (
    <Button
      className={`btn btn-primary btn-md upgrade-tokens-tile__button ${className}`}
      onClick={onBtnClick}
      disabled={buttonDisabled}
    >
      {btnText}
    </Button>
  )
}

UpgradeTokensTile.Link = ({ linkText, to = "", className = "" }) => {
  return (
    <a
      href={to}
      rel="noopener noreferrer"
      target="_blank"
      className={`btn btn-primary btn-md upgrade-tokens-tile__button ${className}`}
    >
      {linkText} ↗
    </a>
  )
}

UpgradeTokensTile.NavLink = ({ linkText, to = "", className = "" }) => {
  return (
    <NavLink
      to={to}
      className={`btn btn-primary btn-md upgrade-tokens-tile__button ${className}`}
    >
      {linkText}
    </NavLink>
  )
}

UpgradeTokensTile.Row = UpgradeTokenFileRow

export default UpgradeTokensTile
