import React, { useCallback } from "react"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { DataTable, Column } from "./DataTable"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { displayAmount } from "../utils/token.utils"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { shortenAddress } from "../utils/general.utils"

const AuthorizeContracts = ({
  data,
  onAuthorizeBtn,
  onDeauthorizeBtn,
  onSelectOperator,
  selectedOperator,
  filterDropdownOptions,
}) => {
  return (
    <section className="tile">
      <DataTable
        data={data}
        itemFieldId="operatorAddress"
        title="Authorize Contracts"
        subtitle="Below are the available operator contracts to authorize."
        withTooltip
        tooltipProps={{
          text:
            "By authorizing a contract, you are approving a set of terms for the governance of an operator, e.g. the rules for slashing tokens.",
        }}
        noDataMessage="No contracts to authorize."
        withFilterDropdown
        filterDropdownProps={{
          options: filterDropdownOptions,
          onSelect: onSelectOperator,
          valuePropertyName: "operatorAddress",
          labelPropertyName: "operatorAddress",
          selectedItem: selectedOperator,
          noItemSelectedText: "All operators",
          renderOptionComponent: ({ operatorAddress }) => (
            <OperatorDropdownItem operatorAddress={operatorAddress} />
          ),
          selectedItemComponent: (
            <OperatorDropdownItem
              operatorAddress={selectedOperator.operatorAddress}
            />
          ),
          allItemsFilterText: "All Operators",
        }}
      >
        <Column
          header="operator address"
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
              onDeauthorizeBtn={onDeauthorizeBtn}
            />
          )}
        />
      </DataTable>
    </section>
  )
}

const Contracts = ({
  contracts,
  operatorAddress,
  onAuthorizeBtn,
  onDeauthorizeBtn,
}) => {
  return (
    <ul className="line-separator">
      {contracts.map((contract) => (
        <AuthorizeContractItem
          key={contract.contractName}
          {...contract}
          operatorAddress={operatorAddress}
          onAuthorizeBtn={onAuthorizeBtn}
          onDeauthorizeBtn={onDeauthorizeBtn}
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
  onDeauthorizeBtn,
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

  const onDeauthorize = useCallback(
    async (transactionHashCallback) => {
      await onDeauthorizeBtn(
        { operatorAddress, contractName },
        transactionHashCallback
      )
    },
    [contractName, operatorAddress, onDeauthorizeBtn]
  )

  return (
    <li className="pb-1 mt-1">
      <div className="flex row wrap space-between center">
        <div>
          <div className="text-big">{contractName}</div>
          <ViewAddressInBlockExplorer address={operatorContractAddress} />
        </div>
        {isAuthorized ? (
          <div>
            <StatusBadge status={BADGE_STATUS.COMPLETE} text="authorized" />
            {contractName === "TBTCSystem" && (
              <SubmitButton
                onSubmitAction={onDeauthorize}
                className="btn btn-secondary btn-sm ml-1"
              >
                deauthorize
              </SubmitButton>
            )}
          </div>
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

const OperatorDropdownItem = React.memo(({ operatorAddress }) => (
  <span key={operatorAddress} title={operatorAddress}>
    {shortenAddress(operatorAddress)}
  </span>
))

export default AuthorizeContracts
