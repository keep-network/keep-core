import React from "react"
import { ViewInBlockExplorer } from "./ViewInBlockExplorer"

const TransactionIsPendingMsgContent = ({ title, txHash }) => {
  return (
    <div className="flex flex-1 row wrap">
      <span className="text-small text-grey-70">{title}</span>
      <span className="text-caption ml-1">Transaction hash: {txHash}</span>
      <ViewInBlockExplorer
        type="tx"
        className="arrow-link grey text-grey-70"
        style={{ marginLeft: "auto" }}
        id={txHash}
      />
    </div>
  )
}

export default React.memo(TransactionIsPendingMsgContent)
