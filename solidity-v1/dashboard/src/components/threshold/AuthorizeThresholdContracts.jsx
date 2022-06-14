import React, { useCallback } from "react"
import AddressShortcut from "./../AddressShortcut"
import Button from "../Button"
import { DataTable, Column } from "../DataTable"
import { ViewAddressInBlockExplorer } from "../ViewInBlockExplorer"
import { KEEP } from "../../utils/token.utils"
import { shortenAddress } from "../../utils/general.utils"
import resourceTooltipProps from "../../constants/tooltips"
import * as Icons from "../Icons"
import ReactTooltip from "react-tooltip"

const AuthorizeThresholdContracts = ({
  data,
  onAuthorizeBtn,
  onStakeBtn,
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
          renderContent={({ stakeAmount, isFromGrant }) => {
            return (
              <>
                <div>{KEEP.displayAmountWithSymbol(stakeAmount)}</div>
                <div className={"text-grey-50"} style={{ fontSize: "14px" }}>
                  {isFromGrant ? "Grant Tokens" : "Wallet Tokens"}
                </div>
              </>
            )
          }}
        />
        <Column
          header="contract"
          field=""
          renderContent={({ contract, operatorAddress }) => (
            <AuthorizeContractItem
              key={contract.contractName}
              {...contract}
              operatorAddress={operatorAddress}
            />
          )}
        />
        <Column
          headerStyle={{ width: "20%", textAlign: "right" }}
          header="actions"
          tdStyles={{ textAlign: "right" }}
          field=""
          renderContent={({
            contract,
            owner,
            operatorAddress,
            authorizerAddress,
            beneficiaryAddress,
            isStakedToT,
            stakeAmount,
            isFromGrant,
            canBeMovedToT,
            isInInitializationPeriod,
          }) => (
            <AuthorizeActions
              key={contract.contractName}
              {...contract}
              isStakedToT={isStakedToT}
              owner={owner}
              operatorAddress={operatorAddress}
              authorizerAddress={authorizerAddress}
              beneficiaryAddress={beneficiaryAddress}
              stakeAmount={stakeAmount}
              isFromGrant={isFromGrant}
              canBeMovedToT={canBeMovedToT}
              onAuthorizeBtn={onAuthorizeBtn}
              onStakeBtn={onStakeBtn}
              isInInitializationPeriod={isInInitializationPeriod}
            />
          )}
        />
      </DataTable>
    </section>
  )
}

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
  owner,
  operatorAddress,
  authorizerAddress,
  beneficiaryAddress,
  stakeAmount,
  isAuthorized,
  isFromGrant,
  canBeMovedToT,
  onAuthorizeBtn,
  onStakeBtn,
  isInInitializationPeriod,
}) => {
  const onAuthorize = useCallback(
    async (awaitingPromise) => {
      await onAuthorizeBtn(
        {
          owner,
          operatorAddress,
          authorizerAddress,
          beneficiaryAddress,
          stakeAmount,
          contractName,
          isFromGrant,
          canBeMovedToT,
        },
        awaitingPromise
      )
    },
    [
      contractName,
      owner,
      operatorAddress,
      authorizerAddress,
      beneficiaryAddress,
      stakeAmount,
      onAuthorizeBtn,
      isFromGrant,
      canBeMovedToT,
    ]
  )

  const onStake = useCallback(
    async (awaitingPromise) => {
      await onStakeBtn(
        {
          owner,
          operatorAddress,
          authorizerAddress,
          beneficiaryAddress,
          stakeAmount,
          contractName,
          isAuthorized,
          isFromGrant,
          canBeMovedToT,
        },
        awaitingPromise
      )
    },
    [
      contractName,
      owner,
      operatorAddress,
      authorizerAddress,
      beneficiaryAddress,
      stakeAmount,
      isAuthorized,
      onStakeBtn,
      isFromGrant,
      canBeMovedToT,
    ]
  )

  return isAuthorized ? (
    <Button
      onClick={onStake}
      className="btn btn-secondary btn-semi-sm"
      style={{ marginLeft: "auto" }}
    >
      <Icons.AlertFill
        data-tip
        data-for={`stake-tooltip-for-operator-${operatorAddress}`}
        className={"tooltip--button-corner"}
      />
      <ReactTooltip
        id={`stake-tooltip-for-operator-${operatorAddress}`}
        place="top"
        type="dark"
        effect={"solid"}
        className={"react-tooltip-base react-tooltip-base--arrow-right"}
        offset={{ left: "100%!important" }}
      >
        <span>
          The stake amount is not yet confirmed. Click “Stake” to confirm the
          stake amount. This stake is not staked on Threshold until it is
          confirmed.
        </span>
      </ReactTooltip>
      stake
    </Button>
  ) : isInInitializationPeriod ? (
    <>
      <ReactTooltip
        id={`stake-tooltip-for-operator-${operatorAddress}`}
        place="top"
        type="dark"
        effect={"solid"}
        className={"react-tooltip-base react-tooltip-base--arrow-right"}
        offset={{ left: "100%!important" }}
      >
        <span>
          This stake is still in initialization period. You will be able to move
          the stake to T when the initialization period ends.
        </span>
      </ReactTooltip>
      <Button
        onClick={onAuthorize}
        className="btn btn-secondary btn-semi-sm"
        style={{ marginLeft: "auto" }}
        disabled={true}
      >
        <Icons.QuestionFill
          data-tip
          data-for={`stake-tooltip-for-operator-${operatorAddress}`}
          className={"tooltip--button-corner"}
        />
        authorize and stake
      </Button>
    </>
  ) : (
    <Button
      onClick={onAuthorize}
      className="btn btn-secondary btn-semi-sm"
      style={{ marginLeft: "auto" }}
    >
      authorize and stake
    </Button>
  )
}

const OperatorDropdownItem = React.memo(({ operatorAddress }) => (
  <span key={operatorAddress} title={operatorAddress}>
    {shortenAddress(operatorAddress)}
  </span>
))

export default AuthorizeThresholdContracts
