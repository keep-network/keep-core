import React, { useCallback } from "react"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { DataTable, Column } from "./DataTable"
import Tile from "./Tile"
import ViewAddressInBlockExplorer from "./ViewAddressInBlockExplorer"
import { displayAmount } from "../utils/token.utils"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"

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
          headerStyle={{ width: "55%" }}
          header="operator contract details"
          field=""
          renderContent={({ contracts, operatorAddress }) => (
            <Contracts
              contracts={contracts}
              operatorAddress={operatorAddress}
              onAuthorizeBtn={onAuthorizeBtn}
            />
          )}
        />
      </DataTable>
    </Tile>
  )
}

const Contracts = ({ contracts, operatorAddress, onAuthorizeBtn }) => {
  return (
    <ul className="line-separator">
      {contracts.map((contract) => (
        <AuthorizeContractItem
          key={contract.contractName}
          {...contract}
          operatorAddress={operatorAddress}
          onAuthorizeBtn={onAuthorizeBtn}
        />
      ))}
    </ul>
  )
}

const AuthorizeContractItem = ({
  contractName,
  operatorAddress,
  isAuthorized,
  operatorContractAddress,
  onAuthorizeBtn,
}) => {
  const onAuthorize = useCallback(
    async (transactionHashCallback) => {
      await onAuthorizeBtn(
        { operatorAddress, contractName },
        transactionHashCallback
      )
    },
    [contractName, operatorAddress, onAuthorizeBtn]
  )
  return (
    <li className="pb-1 mt-1">
      <div className="flex row wrap space-between center">
        <div>
          <div className="text-big">{contractName}</div>
          <ViewAddressInBlockExplorer address={operatorContractAddress} />
        </div>
        {isAuthorized ? (
          <StatusBadge
            className="self-start"
            status={BADGE_STATUS.COMPLETE}
            text="authorized"
          />
        ) : (
          <SubmitButton
            onSubmitAction={onAuthorize}
            className="btn btn-secondary btn-sm"
            style={{ marginLeft: "auto" }}
          >
            authorize
          </SubmitButton>
        )}
      </div>
    </li>
  )
}

export default AuthorizeContracts
