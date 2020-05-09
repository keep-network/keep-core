import React from "react"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { DataTable, Column } from "./DataTable"
import Tile from "./Tile"
import ViewAddressInBlockExplorer from "./ViewAddressInBlockExplorer"
import { displayAmount } from "../utils/token.utils"

const AuthorizeContracts = ({
  data,
  onAuthorizeBtn,
  onAuthorizeSuccessCallback,
}) => {
  return (
    <Tile
      title="Authorize Contracts"
      withTooltip
      tooltipProps={{
        text:
          "By authorizing a contract, you are approving a set of terms for the governance of an operator, e.g. the rules for slashing tokens.",
      }}
    >
      <DataTable data={data} itemFieldId="operatorAddress">
        <Column
          header="opeartor address"
          field="operatorAddress"
          renderContent={({ operatorAddress }) => (
            <AddressShortcut address={operatorAddress} />
          )}
        />
        <Column
          header="stake"
          field="stakeAmount"
          renderContent={({ stakeAmount }) =>
            `${displayAmount(stakeAmount)} KEEP`
          }
        />
        <Column
          header="operator contract details"
          field=""
          renderContent={({ contracts }) => <Contracts contracts={contracts} />}
        />
      </DataTable>
    </Tile>
  )
}

const Contracts = ({ contracts }) => {
  return (
    <ul className="line-separator">
      {contracts.map((contract) => (
        <AuthorizeContractItem key={contract.contractName} {...contract} />
      ))}
    </ul>
  )
}

const AuthorizeContractItem = ({ contractName, operatorContractAddress }) => {
  return (
    <li className="pb-1 mt-1">
      <div className="flex row wrap space-between center">
        <div>
          <div className="text-big">{contractName}</div>
          <ViewAddressInBlockExplorer address={operatorContractAddress} />
        </div>
        <SubmitButton
          className="btn btn-secondary btn-sm"
          style={{ marginLeft: "auto" }}
        >
          authorize
        </SubmitButton>
      </div>
    </li>
  )
}

export default AuthorizeContracts
