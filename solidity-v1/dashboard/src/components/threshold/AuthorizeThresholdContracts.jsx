import React, { useCallback } from "react"
import AddressShortcut from "./../AddressShortcut"
import { SubmitButton } from "../Button"
import { DataTable, Column } from "../DataTable"
import { ViewAddressInBlockExplorer } from "../ViewInBlockExplorer"
import { KEEP } from "../../utils/token.utils"
import { shortenAddress } from "../../utils/general.utils"
import resourceTooltipProps from "../../constants/tooltips"

const AuthorizeThresholdContracts = ({
  data,
  onAuthorizeBtn,
  onStakeBtn,
  onSelectOperator,
  selectedOperator,
  filterDropdownOptions,
  onSuccessCallback,
}) => {
  return (
    <section className="tile">
      <DataTable
        data={data}
        itemFieldId="operatorAddress"
        title="Authorize Contracts"
        subtitle="Below are the available operator contracts to authorize."
        withTooltip
        tooltipProps={resourceTooltipProps.authorize}
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
        centered
      >
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
          header="contract"
          field=""
          renderContent={({ contracts, operatorAddress }) => (
            <AuthorizeContractItem
              key={contracts[0].contractName}
              {...contracts[0]}
              operatorAddress={operatorAddress}
            />
          )}
        />
        <Column
          headerStyle={{ width: "20%", textAlign: "right" }}
          header="actions"
          tdStyles={{ textAlign: "right" }}
          field=""
          renderContent={({ contracts, operatorAddress, isStakedToT }) => (
            <AuthorizeActions
              key={contracts[0].contractName}
              {...contracts[0]}
              isStakedToT={isStakedToT}
              operatorAddress={operatorAddress}
              onAuthorizeBtn={onAuthorizeBtn}
              onStakeBtn={onStakeBtn}
              onSuccessCallback={onSuccessCallback}
            />
          )}
        />
      </DataTable>
    </section>
  )
}

// const styles = {
//   tooltipContentWrapper: { textAlign: "left", minWidth: "15rem" },
// }

const AuthorizeContractItem = ({ contractName, operatorContractAddress }) => {
  return (
    <div className="flex row wrap space-between center">
      <div>
        <div className="text-big">{contractName}</div>
        <ViewAddressInBlockExplorer address={operatorContractAddress} />
      </div>
    </div>
  )
}

const AuthorizeActions = ({
  contractName,
  operatorAddress,
  isAuthorized,
  onAuthorizeBtn,
  onStakeBtn,
  onSuccessCallback,
}) => {
  const onAuthorize = useCallback(
    async (awaitingPromise) => {
      await onAuthorizeBtn({ operatorAddress, contractName }, awaitingPromise)
    },
    [contractName, operatorAddress, onAuthorizeBtn]
  )

  const onStake = useCallback(
    async (awaitingPromise) => {
      await onStakeBtn({ operatorAddress, contractName }, awaitingPromise)
    },
    [contractName, operatorAddress, onStakeBtn]
  )

  const onSuccess = useCallback(
    (isAuthorized = true) => {
      onSuccessCallback(contractName, operatorAddress, isAuthorized)
    },
    [onSuccessCallback, contractName, operatorAddress]
  )

  return isAuthorized ? (
    <SubmitButton
      onSubmitAction={onStake}
      className="btn btn-secondary btn-semi-sm"
      style={{ marginLeft: "auto" }}
      successCallback={onSuccess}
    >
      stake
    </SubmitButton>
  ) : (
    <SubmitButton
      onSubmitAction={onAuthorize}
      className="btn btn-secondary btn-semi-sm"
      style={{ marginLeft: "auto" }}
      successCallback={onSuccess}
    >
      authorize and stake
    </SubmitButton>
  )
}

const OperatorDropdownItem = React.memo(({ operatorAddress }) => (
  <span key={operatorAddress} title={operatorAddress}>
    {shortenAddress(operatorAddress)}
  </span>
))

export default AuthorizeThresholdContracts
