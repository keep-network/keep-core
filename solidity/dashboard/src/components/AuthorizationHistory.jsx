import React from "react"
import AddressShortcut from "./AddressShortcut"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { ETHERSCAN_DEFAULT_URL } from "../constants/constants"
import { DataTable, Column } from "./DataTable"

const AuthorizationHistory = ({ contracts }) => {
  return (
    <section className="tile">
      <h3 className="text-grey-60">Authorization History</h3>
      <DataTable data={contracts || []} itemFieldId="contractAddress">
        <Column
          header="contract address"
          field="contractAddress"
          renderContent={({ contractAddress }) => (
            <AddressShortcut address={contractAddress} />
          )}
        />
        <Column
          header="status"
          field="status"
          renderContent={({ status }) => (
            <StatusBadge
              className="self-start"
              status={BADGE_STATUS.COMPLETE}
              text="authorized"
            />
          )}
        />
        <Column
          header="contract details"
          field="details"
          renderContent={({ contractAddress }) => (
            <a
              href={ETHERSCAN_DEFAULT_URL + contractAddress}
              rel="noopener noreferrer"
              target="_blank"
            >
              View in Block Explorer
            </a>
          )}
        />
      </DataTable>
    </section>
  )
}

export default AuthorizationHistory
