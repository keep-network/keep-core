import React from "react"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import TokenAmount from "./TokenAmount"
import * as Icons from "./Icons"
import { ViewInBlockExplorer } from "./ViewInBlockExplorer"

const TBTCRewardsDataTable = ({ rewards }) => {
  return (
    <DataTable data={rewards} itemFieldId="depositTokenId">
      <Column
        header="amount"
        field="amount"
        renderContent={({ amount }) => (
          <TokenAmount
            currencyIcon={Icons.TBTC}
            currencyIconProps={{ width: 15, height: 15 }}
            amount={amount}
            amountClassName="text-big text-grey-70"
            withMetricSuffix={false}
          />
        )}
      />
      <Column
        header="transactionHash"
        field="transactionHash"
        renderContent={({ transactionHash }) => (
          <ViewInBlockExplorer type="tx" id={transactionHash} />
        )}
      />
      <Column
        header="deposit token id"
        field="depositTokenId"
        renderContent={({ depositTokenId }) => (
          <AddressShortcut address={depositTokenId} />
        )}
      />
      <Column
        header="operator"
        field="operator"
        renderContent={({ depositTokenId }) => <span>Load address</span>}
      />
    </DataTable>
  )
}

export default React.memo(TBTCRewardsDataTable)
