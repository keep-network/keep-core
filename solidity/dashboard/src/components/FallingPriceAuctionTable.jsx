import React from "react"
import { displayAmount } from "../utils/token.utils"
import moment from "moment"
import { DataTable, Column } from "./DataTable"
import Tile from "./Tile"


const FallingPriceAuctionTable = ({
  auctionScheduleData,
}) => {
  return (
    <Tile>
      <DataTable
        title="Falling-Price Auction Schedule"
        data={auctionScheduleData}
        itemFieldId="releasedInTimestamp"
        noDataMessage="No schedule available."
        cellStyle={{padding: "1rem 0.6rem"}}
      >
        <Column
          header="Deposit on Offer (%)"
          field="depositPctOnOffer"
          renderContent={({ depositPctOnOffer }) => `${(depositPctOnOffer)}`}
        />
        <Column
          header="Amount (ETH)"
          field="amountOnOffer"
          renderContent={({ amountOnOffer: wei }) => `${displayAmount(wei, false)} Ξ`}
        />
        <Column
          header="Released at approx."
          field="releasedInTimestamp"
          renderContent={({ releasedInTimestamp }) => `${moment.unix(releasedInTimestamp)}` }
        />
          {/* (
          <div>
            {`${moment.unix(releasedInTimestamp)}`}<br/>
            {`${moment.unix(releasedInTimestamp).fromNow()}`}
          </div>
          ) */}
      </DataTable>
    </Tile>
  )
}

FallingPriceAuctionTable.defaultProps = {
  title: "Delegations",
}

export default FallingPriceAuctionTable
