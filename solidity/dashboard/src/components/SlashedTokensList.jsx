import React from "react"
import { formatDate } from "../utils/general.utils"
import { displayAmount } from "../utils/token.utils"
import AddressShortcut from "./AddressShortcut"
import { DataTable, Column } from "./DataTable"
import resourceTooltipProps from "../constants/tooltips"

const SlashedTokensList = ({ slashedTokens }) => {
  return (
    <DataTable
      data={slashedTokens}
      itemFieldId="id"
      title="Slashed Tokens"
      withTooltip
      tooltipProps={resourceTooltipProps.slashing}
      noDataMessage="No slashed tokens."
    >
      <Column
        header="amount"
        headerStyle={{ width: "30%" }}
        field="amount"
        renderContent={({ amount }) => (
          <>
            <span className="text-error">
              {amount && `-${displayAmount(amount)} `}
            </span>
            <span className="text-grey-40">KEEP</span>
          </>
        )}
      />
      <Column
        header="details"
        field="event"
        renderContent={({ event, groupPublicKey, date }) => (
          <>
            <div className="text-big text-grey-70">
              Group&nbsp;
              <AddressShortcut
                address={groupPublicKey}
                classNames="text-big text-grey-70"
              />
              &nbsp;
              {event === "UnauthorizedSigningReported"
                ? "key was leaked. Private key was published outside of the members of the signing group."
                : "was selected to do work and not enough members participated."}
            </div>
            <div className="text-small text-grey-50">
              {formatDate(date, "MMM DD, YYYY")}
            </div>
          </>
        )}
      />
    </DataTable>
  )
}

export default React.memo(SlashedTokensList)
