import React, { useCallback } from "react"
import { formatDate } from "../utils/general.utils"
import { displayAmount } from "../utils/token.utils"
import AddressShortcut from "./AddressShortcut"
import UndelegateStakeButton from "./UndelegateStakeButton"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { PENDING_STATUS, COMPLETE_STATUS } from "../constants/constants"
import { DataTable, Column } from "./DataTable"
import Tile from "./Tile"
import { SubmitButton } from "./Button"
import { useShowMessage, messageType } from "./Message"
import { useWeb3Context } from "./WithWeb3Context"
import { tokensPageService } from "../services/tokens-page.service"
import { useModal } from "../hooks/useModal"
import AddTopUpModal from "./AddTopUpModal"

const DelegatedTokensTable = ({
  delegatedTokens,
  cancelStakeSuccessCallback,
  availableToStake,
}) => {
  const showMessage = useShowMessage()
  const web3Context = useWeb3Context()
  const { openConfirmationModal } = useModal()

  const onTopUpBtn = useCallback(
    async (delegationData, transactionHashCallback) => {
      try {
        const { amount } = await openConfirmationModal(
          {
            modalOptions: { title: "Add KEEP" },
            submitBtnText: "add keep",
            availableAmount: availableToStake,
            currentAmount: delegationData.amount,
            ...delegationData,
          },
          AddTopUpModal
        )
        delegationData.beneficiaryAddress = delegationData.beneficiary
        delegationData.stakeTokens = amount
        delegationData.selectedGrant = {
          id: delegationData.grantId,
          isManagedGrant: delegationData.isManagedGrant,
          managedGrantContractInstance:
            delegationData.managedGrantContractInstance,
        }
        delegationData.context = delegationData.isFromGrant
          ? "granted"
          : "owned"
        await tokensPageService.delegateStake(
          web3Context,
          delegationData,
          transactionHashCallback
        )
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Top up committed successfully",
        })
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Commit action has failed ",
          content: error.message,
        })
      }
    },
    [showMessage, openConfirmationModal, availableToStake, web3Context]
  )
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
          renderContent={({ amount }) => `${displayAmount(amount)} KEEP`}
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
          header=""
          field=""
          renderContent={(delegation) => (
            <>
              <div>
                <UndelegateStakeButton
                  isInInitializationPeriod={delegation.isInInitializationPeriod}
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
                />
              </div>
              <div className="mt-1">
                <SubmitButton
                  className="btn btn-secondary btn-sm"
                  onSubmitAction={(transactionHashCallback) =>
                    onTopUpBtn(delegation, transactionHashCallback)
                  }
                >
                  add keep
                </SubmitButton>
              </div>
            </>
          )}
        />
      </DataTable>
    </Tile>
  )
}

DelegatedTokensTable.defaultProps = {
  title: "Delegations",
}

export default DelegatedTokensTable
