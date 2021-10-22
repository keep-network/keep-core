import React from "react"
import TokenAmount from "../TokenAmount"
import Button from "../Button"

const UpgradeTokensTile = ({
  title,
  className = "",
  btnText,
  onBtnClick,
  children,
}) => {
  return (
    <div className={`upgrade-tokens-tile ${className}`}>
      <div className="upgrade-tokens-tile__title">{title}</div>
      <div>{children}</div>
      <Button
        className="btn btn-primary btn-lg upgrade-tokens-tile__button"
        onClick={onBtnClick}
      >
        {btnText}
      </Button>
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
