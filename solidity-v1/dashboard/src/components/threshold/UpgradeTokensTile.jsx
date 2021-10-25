import React from "react"
import TokenAmount from "../TokenAmount"
import Button from "../Button"
import OnlyIf from "../OnlyIf"
import { LINK } from "../../constants/constants"

const UpgradeTokensTile = ({
  title,
  className = "",
  btnText,
  onBtnClick,
  isLink = false,
  children,
}) => {
  return (
    <div className={`upgrade-tokens-tile ${className}`}>
      <div className="upgrade-tokens-tile__title">{title}</div>
      <div>{children}</div>
      <OnlyIf condition={!isLink}>
        <Button
          className="btn btn-primary btn-md upgrade-tokens-tile__button"
          onClick={onBtnClick}
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

const UpgradeTokenFileRow = ({ label, amount }) => {
  return (
    <div className={"upgrade-tokens-tile-row"}>
      <span className="text-grey-50">{label}</span>
      <TokenAmount
        amount={amount}
        symbolClassName="upgrade-token-tile-row__token-amount-symbol"
        amountClassName="upgrade-token-tile-row__token-amount-amount"
      />
    </div>
  )
}

UpgradeTokensTile.Row = UpgradeTokenFileRow

export default UpgradeTokensTile
