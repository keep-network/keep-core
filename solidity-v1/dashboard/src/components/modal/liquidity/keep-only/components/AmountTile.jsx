import React from "react"
import MetricsTile from "../../../../MetricsTile"
import TokenAmount from "../../../../TokenAmount"

const styles = {
  amountTileWrapper: {
    justifyContent: "flex-start",
    flexGrow: "1",
    padding: "0.5rem",
    height: "auto",
  },
}
export const AmountTile = ({ amount, title, icon }) => {
  return (
    <MetricsTile
      className="bg-grey-10 self-start"
      style={styles.amountTileWrapper}
    >
      <h5 className="text-grey-40 text-left mb-1">{title}</h5>
      <TokenAmount
        wrapperClassName="mb-1"
        amount={amount}
        icon={icon}
        withIcon
        withMetricSuffix
      />
    </MetricsTile>
  )
}
