import React from "react"
import AddressShortcut from "./AddressShortcut"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { DataTable, Column } from "./DataTable"
import Tile from "./Tile"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { KEEP } from "../utils/token.utils"

const AuthorizationHistory = ({ contracts }) => {
  return (
    <Tile>
      <DataTable
        data={contracts || []}
        title="Authorizations History"
        itemFieldId="contractAddress"
        noDataMessage="No authorization history."
      >
        <Column
          header="details"
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
          header="operator"
          field="operatorAddress"
          renderContent={({ operatorAddress }) => (
            <AddressShortcut address={operatorAddress} />
          )}
        />
        <Column
          header="stake"
          field="stakeAmount"
          renderContent={({ stakeAmount }) =>
            `${KEEP.displayAmountWithSymbol(stakeAmount)}`
          }
        />
        <Column
          header="operator contract details"
          field="details"
          renderContent={({ contractName, operatorContractAddress }) => (
            <div>
              <div className="text-big">{contractName}</div>
              <ViewAddressInBlockExplorer address={operatorContractAddress} />
            </div>
          )}
        />
      </DataTable>
    </Tile>
  )
}

export default AuthorizationHistory
