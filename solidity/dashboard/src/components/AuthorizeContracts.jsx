import React, { useCallback } from "react"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { DataTable, Column } from "./DataTable"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { displayAmount } from "../utils/token.utils"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { shortenAddress } from "../utils/general.utils"
import resourceTooltipProps from "../constants/tooltips"
import Tooltip from "./Tooltip"

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
            `${displayAmount(stakeAmount)} KEEP`
          }
        />
        <Column
          headerStyle={{ width: "55%" }}
          header="operator contract details"
          field=""
          renderContent={({ contracts, operatorAddress }) => (
            <ul className="line-separator">
              {contracts.map((contract) => (
                <AuthorizeContractItem
                  key={contract.contractName}
                  {...contract}
                  operatorAddress={operatorAddress}
                  onAuthorizeBtn={onAuthorizeBtn}
                  onDeauthorizeBtn={onDeauthorizeBtn}
                  onSuccessCallback={onSuccessCallback}
                />
              ))}
            </ul>
          )}
        />
      </DataTable>
    </section>
  )
}

const styles = {
  tooltipContentWrapper: { textAlign: "left", minWidth: "15rem" },
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
            {contractName === "TBTCSystem" && (
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
              contractName === "Keep Random Beacon Operator Contract"
            }
            simple
            delay={10}
            direction="top"
            contentWrapperStyles={styles.tooltipContentWrapper}
            triggerComponent={() => (
              <SubmitButton
                onSubmitAction={onAuthorize}
                className="btn btn-secondary btn-sm"
                style={{ marginLeft: "auto" }}
                successCallback={onSuccess}
                disabled={
                  contractName === "Keep Random Beacon Operator Contract"
                }
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
