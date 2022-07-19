import React, { useCallback } from "react"
import { formatDate } from "../utils/general.utils"
import { KEEP } from "../utils/token.utils"
import AddressShortcut from "./AddressShortcut"
import UndelegateStakeButton from "./UndelegateStakeButton"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { PENDING_STATUS, COMPLETE_STATUS } from "../constants/constants"
import { DataTable, Column } from "./DataTable"
import Tile from "./Tile"
import { SubmitButton } from "./Button"
import { connect } from "react-redux"
import web3Utils from "web3-utils"
import useUpdateInitializedDelegations from "../hooks/useUpdateInitializedDelegations"
import * as Icons from "./Icons"

const DelegatedTokensTable = ({
  delegations,
  cancelStakeSuccessCallback,
  keepTokenBalance,
  grants,
  addKeep,
  undelegationPeriod,
}) => {
  useUpdateInitializedDelegations(delegations)
  const getAvailableToStakeFromGrant = useCallback(
    (grantId) => {
      const grant = grants.find(({ id }) => id === grantId)

      return grant ? grant.availableToStake : 0
    },
    [grants]
  )

  const isAddKeepBtnDisabled = (delegationData) => {
    const availableAmount = delegationData.isFromGrant
      ? getAvailableToStakeFromGrant(delegationData.grantId)
      : keepTokenBalance

    return web3Utils.toBN(availableAmount).lten(0)
  }

  const onTopUpBtn = async (delegationData, awaitingPromise) => {
    const availableAmount = delegationData.isFromGrant
      ? getAvailableToStakeFromGrant(delegationData.grantId)
      : keepTokenBalance
    addKeep(
      {
        ...delegationData,
        beneficiaryAddress: delegationData.beneficiary,
        currentAmount: delegationData.amount,
        availableAmount,
      },
      awaitingPromise
    )
  }

  return (
    <Tile>
      <DataTable
        title="Delegations"
        data={delegations}
        itemFieldId="operatorAddress"
        noDataMessage="No delegated tokens."
        centered
      >
        <Column
          header="amount"
          field="amount"
          renderContent={({ amount, isFromGrant }) => {
            return (
              <>
                <div>{KEEP.displayAmountWithSymbol(amount)}</div>
                <div className={"text-grey-50"} style={{ fontSize: "14px" }}>
                  {isFromGrant ? "Grant Tokens" : "Wallet Tokens"}
                </div>
              </>
            )
          }}
        />
        <Column
          header="status"
          field="delegationStatus"
          renderContent={(delegation) => {
            const delegationStatus = delegation.isInInitializationPeriod
              ? PENDING_STATUS
              : COMPLETE_STATUS
            const statusBadgeText =
              delegationStatus === PENDING_STATUS
                ? `${delegationStatus.toLowerCase()}, ${delegation.initializationOverAt.fromNow(
                    true
                  )}`
                : formatDate(delegation.initializationOverAt)

            return (
              <>
                <StatusBadge
                  status={BADGE_STATUS[delegationStatus]}
                  className="self-start"
                  text={statusBadgeText}
                  onlyIcon={delegationStatus === COMPLETE_STATUS}
                />
                <div className={"text-grey-50"} style={{ fontSize: "14px" }}>
                  {delegation.initializationOverAt.format("HH:mm:ss")}
                </div>
              </>
            )
          }}
        />
        <Column
          header="beneficiary"
          field="beneficiary"
          renderContent={({ beneficiary }) => (
            <AddressShortcut address={beneficiary} />
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
          header="authorizer"
          field="authorizerAddress"
          renderContent={({ authorizerAddress }) => (
            <AddressShortcut address={authorizerAddress} />
          )}
        />
        <Column
          headerStyle={{ width: "25%", textAlign: "right" }}
          header="actions"
          field=""
          renderContent={(delegation) =>
            delegation.isCopiedStake ? (
              <StatusBadge
                status={BADGE_STATUS.COMPLETE}
                className="self-start"
                text="stake copied"
              />
            ) : (
              <div className="flex row center justify-right">
                <div>
                  <UndelegateStakeButton
                    isInInitializationPeriod={
                      delegation.isInInitializationPeriod
                    }
                    isFromGrant={delegation.isFromGrant}
                    btnClassName="btn btn-semi-sm btn-secondary"
                    operator={delegation.operatorAddress}
                    amount={delegation.amount}
                    authorizer={delegation.authorizerAddress}
                    beneficiary={delegation.beneficiary}
                    undelegationPeriod={undelegationPeriod}
                    successCallback={
                      delegation.isInInitializationPeriod
                        ? cancelStakeSuccessCallback
                        : () => {}
                    }
                    disabled={
                      delegation.isTStakingContractAuthorized &&
                      delegation.isStakedToT
                    }
                  />
                </div>
                <div className={"ml-2"}>
                  <SubmitButton
                    className="btn btn-secondary btn-semi-sm"
                    onSubmitAction={(awaitingPromise) =>
                      onTopUpBtn(delegation, awaitingPromise)
                    }
                    disabled={isAddKeepBtnDisabled(delegation)}
                  >
                    <span className={"flex row center"}>
                      <Icons.KeepBlackGreen width={12} height={12} />
                      &nbsp;add keep
                    </span>
                  </SubmitButton>
                </div>
              </div>
            )
          }
        />
      </DataTable>
    </Tile>
  )
}

DelegatedTokensTable.defaultProps = {
  title: "Delegations",
}

const mapDispatchToProps = (dispatch) => ({
  addKeep: (values, meta) =>
    dispatch({
      type: "staking/top-up",
      payload: values,
      meta,
    }),
})

export default connect(null, mapDispatchToProps)(DelegatedTokensTable)
