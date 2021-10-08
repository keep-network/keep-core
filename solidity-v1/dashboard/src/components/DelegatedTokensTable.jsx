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
import { useModal } from "../hooks/useModal"
import AddTopUpModal, { TopUpInitiatedConfirmationModal } from "./AddTopUpModal"
import { connect } from "react-redux"
import web3Utils from "web3-utils"

const DelegatedTokensTable = ({
  delegatedTokens,
  cancelStakeSuccessCallback,
  keepTokenBalance,
  grants,
  addKeep,
  undelegationPeriod,
}) => {
  const { openConfirmationModal, openModal } = useModal()

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
    const { amount } = await openConfirmationModal(
      {
        modalOptions: { title: "Add KEEP" },
        submitBtnText: "add keep",
        availableAmount,
        currentAmount: delegationData.amount,
        minimumAmount: 1,
        ...delegationData,
      },
      AddTopUpModal
    )
    addKeep(
      {
        ...delegationData,
        amount,
        beneficiaryAddress: delegationData.beneficiary,
      },
      awaitingPromise
    )
    try {
      await awaitingPromise.promise
      openModal(
        <TopUpInitiatedConfirmationModal
          {...delegationData}
          addedAmount={KEEP.fromTokenUnit(amount).toString()}
          currentAmount={delegationData.amount}
        />,
        {
          title: "Add KEEP",
        }
      )
    } catch (error) {
      console.error("Unexpected error", error)
    }
  }

  return (
    <Tile>
      <DataTable
        title="Delegations"
        data={delegatedTokens}
        itemFieldId="operatorAddress"
        noDataMessage="No delegated tokens."
      >
        <Column
          header="amount"
          field="amount"
          renderContent={({ amount }) =>
            `${KEEP.displayAmountWithSymbol(amount)}`
          }
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
              <StatusBadge
                status={BADGE_STATUS[delegationStatus]}
                className="self-start"
                text={statusBadgeText}
                onlyIcon={delegationStatus === COMPLETE_STATUS}
              />
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
          headerStyle={{ width: "25%" }}
          header=""
          field=""
          renderContent={(delegation) =>
            delegation.isCopiedStake ? (
              <StatusBadge
                status={BADGE_STATUS.COMPLETE}
                className="self-start"
                text="stake copied"
              />
            ) : (
              <div className="flex row center space-evenly">
                <div>
                  <UndelegateStakeButton
                    isInInitializationPeriod={
                      delegation.isInInitializationPeriod
                    }
                    isFromGrant={delegation.isFromGrant}
                    btnClassName="btn btn-sm btn-secondary"
                    operator={delegation.operatorAddress}
                    isManagedGrant={delegation.isManagedGrant}
                    managedGrantContractInstance={
                      delegation.managedGrantContractInstance
                    }
                    successCallback={
                      delegation.isInInitializationPeriod
                        ? cancelStakeSuccessCallback
                        : () => {}
                    }
                    undelegationPeriod={undelegationPeriod}
                  />
                </div>
                <div>
                  <SubmitButton
                    className="btn btn-secondary btn-sm"
                    onSubmitAction={(awaitingPromise) =>
                      onTopUpBtn(delegation, awaitingPromise)
                    }
                    disabled={isAddKeepBtnDisabled(delegation)}
                  >
                    add keep
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
      type: "staking/delegate_request",
      payload: values,
      meta,
    }),
})

export default connect(null, mapDispatchToProps)(DelegatedTokensTable)
