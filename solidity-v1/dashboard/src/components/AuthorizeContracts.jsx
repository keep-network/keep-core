import React, { useCallback } from "react"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { DataTable, Column } from "./DataTable"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { KEEP } from "../utils/token.utils"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { shortenAddress } from "../utils/general.utils"
import resourceTooltipProps from "../constants/tooltips"
import Tooltip, { TOOLTIP_DIRECTION } from "./Tooltip"
import { AUTH_CONTRACTS_LABEL } from "../constants/constants"

const AuthorizeContracts = ({
  data,
  onAuthorizeBtn,
  onDeauthorizeBtn,
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
            `${KEEP.displayAmountWithSymbol(stakeAmount)}`
          }
        />
        <Column
          headerStyle={{ width: "55%" }}
          header="operator contract details"
          field=""
          renderContent={({ contracts, operatorAddress }) => (
            <ContractsToAuthorizeCell
              contracts={contracts}
              renderAuthContract={(contract) => (
                <AuthorizeContractItem
                  key={contract.contractName}
                  {...contract}
                  operatorAddress={operatorAddress}
                  onAuthorizeBtn={onAuthorizeBtn}
                  onDeauthorizeBtn={onDeauthorizeBtn}
                  onSuccessCallback={onSuccessCallback}
                />
              )}
            />
          )}
        />
      </DataTable>
    </section>
  )
}

const styles = {
  tooltipContentWrapper: { textAlign: "left", minWidth: "15rem" },
}

const ContractsToAuthorizeCell = ({ renderAuthContract, contracts }) => {
  return <ul className="line-separator">{contracts.map(renderAuthContract)}</ul>
}

const AuthorizeContractItem = ({
  contractName,
  operatorAddress,
  isAuthorized,
  operatorContractAddress,
  onAuthorizeBtn,
  onDeauthorizeBtn,
  onSuccessCallback,
}) => {
  const onAuthorize = useCallback(
    async (awaitingPromise) => {
      await onAuthorizeBtn({ operatorAddress, contractName }, awaitingPromise)
    },
    [contractName, operatorAddress, onAuthorizeBtn]
  )

  const onDeauthorize = useCallback(
    async (awaitingPromise) => {
      await onDeauthorizeBtn({ operatorAddress, contractName }, awaitingPromise)
    },
    [contractName, operatorAddress, onDeauthorizeBtn]
  )

  const onSuccess = useCallback(
    (isAuthorized = true) => {
      onSuccessCallback(contractName, operatorAddress, isAuthorized)
    },
    [onSuccessCallback, contractName, operatorAddress]
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
            {contractName === AUTH_CONTRACTS_LABEL.TBTC_SYSTEM && (
              <SubmitButton
                onSubmitAction={onDeauthorize}
                className="btn btn-secondary btn-sm ml-1"
                successCallback={() => onSuccess(false)}
              >
                deauthorize
              </SubmitButton>
            )}
          </div>
        ) : (
          <Tooltip
            shouldShowTooltip={
              contractName === AUTH_CONTRACTS_LABEL.RANDOM_BEACON
            }
            simple
            delay={10}
            direction={TOOLTIP_DIRECTION.TOP}
            contentWrapperStyles={styles.tooltipContentWrapper}
            triggerComponent={() => (
              <SubmitButton
                onSubmitAction={onAuthorize}
                className="btn btn-secondary btn-sm"
                style={{ marginLeft: "auto" }}
                successCallback={onSuccess}
                disabled={contractName === AUTH_CONTRACTS_LABEL.RANDOM_BEACON}
              >
                authorize
              </SubmitButton>
            )}
          >
            Keep Random Beacon Operator contract has been disabled due to&nbsp;
            <a
              href="https://docs.keep.network/status-reports/2020-11-11-retro-geth-hardfork.html"
              rel="noopener noreferrer"
              target="_blank"
              className="text-white text-link"
            >
              the impact of the geth hardfork that occurred on 11 November 2020
            </a>
            .
          </Tooltip>
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
