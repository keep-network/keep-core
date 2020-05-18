import React, { useCallback } from "react"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { DataTable, Column } from "./DataTable"
import ViewAddressInBlockExplorer from "./ViewAddressInBlockExplorer"
import { displayAmount } from "../utils/token.utils"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import SpeechBubbleTooltip from "./SpeechBubbleTooltip"
import Dropdown from "./Dropdown"
import { shortenAddress } from "../utils/general.utils"

const AuthorizeContracts = ({
  data,
  onAuthorizeBtn,
  onSelectOperator,
  selectedOperator,
  filterDropdownOptions,
}) => {
  return (
    <section className="tile">
      <div className="flex row wrap center space-between">
        <header>
          <div className="flex row">
            <h4 className="mr-1 text-grey-70">Authorize Contracts</h4>
            <SpeechBubbleTooltip
              text={
                "By authorizing a contract, you are approving a set of terms for the governance of an operator, e.g. the rules for slashing tokens."
              }
            />
          </div>
          <div className="text-grey-40 text-small">
            Below are the available operator contracts to authorize.
          </div>
        </header>
        <div style={{ marginLeft: "auto" }}>
          <Dropdown
            withLabel={false}
            options={filterDropdownOptions}
            onSelect={(operator) => onSelectOperator(operator)}
            valuePropertyName="operatorAddress"
            labelPropertyName="operatorAddress"
            selectedItem={selectedOperator}
            noItemSelectedText="All operators"
            renderOptionComponent={({ operatorAddress }) => (
              <OperatorDropdownItem operatorAddress={operatorAddress} />
            )}
            selectedItemComponent={
              <OperatorDropdownItem
                operatorAddress={selectedOperator.operatorAddress}
              />
            }
            isFilterDropdow
            allItemsFilterText="All Operators"
          />
        </div>
      </div>
      <DataTable
        data={data}
        itemFieldId="operatorAddress"
        noDataMessage="No contracts to authorize."
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
            />
          )}
        />
      </DataTable>
    </section>
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
          <StatusBadge status={BADGE_STATUS.COMPLETE} text="authorized" />
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
