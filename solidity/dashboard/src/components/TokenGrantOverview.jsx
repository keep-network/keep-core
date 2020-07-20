import React, { useContext } from "react"
import { formatDate } from "../utils/general.utils"
import { SubmitButton } from "./Button"
import { colors } from "../constants/colors"
import { CircularProgressBars } from "./CircularProgressBar"
import { Web3Context } from "./WithWeb3Context"
import { useShowMessage, useCloseMessage, messageType } from "./Message"
import moment from "moment"
import { gt } from "../utils/arithmetics.utils"
import { SpeechBubbleTooltip } from "./SpeechBubbleTooltip"
import TokenAmount from "./TokenAmount"
import {
  displayAmountWithMetricSuffix,
  displayAmount,
} from "../utils/token.utils"
import { isEmptyArray } from "../utils/array.utils"
import { ViewAddressInBlockExplorer } from "./ViewInBlockExplorer"
import { contracts } from "../contracts"
import TransactionIsPendingMsgContent from "./TransactionIsPendingMsgContent"

const TokenGrantOverview = ({ selectedGrant, selectedGrantStakedAmount }) => {
  return (
    <>
      <TokenGrantDetails selectedGrant={selectedGrant} />
      <hr />
      <div className="flex">
        <TokenGrantUnlockingdDetails selectedGrant={selectedGrant} />
      </div>
      <div className="flex mt-1">
        <TokenGrantStakedDetails
          selectedGrant={selectedGrant}
          stakedAmount={selectedGrantStakedAmount}
        />
      </div>
    </>
  )
}

export const TokenGrantDetails = ({
  title = "Grant Details",
  selectedGrant,
}) => {
  const cliffPeriod = moment
    .unix(selectedGrant.cliff)
    .from(moment.unix(selectedGrant.start), true)
  const fullyUnlockedDate = moment
    .unix(selectedGrant.start)
    .add(selectedGrant.duration, "seconds")

  return (
    <section className="token-grant-details">
      <div className="flex wrap center">
        <h4 className="mr-1 text-grey-70" style={{ marginRight: "auto" }}>
          {title}
        </h4>
        <SpeechBubbleTooltip
          text={
            <>
              A &nbsp;<span className="text-bold">cliff</span>&nbsp; is a set
              period of time before vesting begins.
            </>
          }
          iconColor={colors.grey60}
          iconBackgroundColor="transparent"
          title={`${cliffPeriod} cliff`}
        />
      </div>
      <TokenAmount amount={selectedGrant.amount} />
      <h4 className="text-grey-30 mb-1">Grant ID {selectedGrant.id}</h4>
      <h5 className="text-grey-50">
        Issued:{" "}
        {selectedGrant.id && formatDate(moment.unix(selectedGrant.start))}
      </h5>
      <h5 className="text-grey-50">
        Fully Unlocked: {selectedGrant.id && formatDate(fullyUnlockedDate)}
      </h5>
    </section>
  )
}

export const TokenGrantUnlockingdDetails = ({
  selectedGrant,
  hideReleaseTokensBtn = false,
}) => {
  const { yourAddress, grantContract, tokenStakingEscrow } = useContext(
    Web3Context
  )
  const showMessage = useShowMessage()
  const closeMessage = useCloseMessage()

  const releaseTokens = async (
    onTransactionHashCallback,
    openMessageInfo,
    setFetching,
    openConfirmationModal
  ) => {
    try {
      const {
        isManagedGrant,
        managedGrantContractInstance,
        escrowOperatorsToWithdraw,
        withdrawableAmountGrantOnly,
      } = selectedGrant
      if (!isEmptyArray(escrowOperatorsToWithdraw)) {
        await openConfirmationModal({
          title: "You’re about to release tokens.",
          subtitle: (
            <>
              <span>You have deposited tokens in the</span>&nbsp;
              <ViewAddressInBlockExplorer
                text="TokenStakingEscrow contract"
                address={contracts.tokenStakingEscrow.options.address}
              />
              <p>
                To withdraw all tokens it may be necessary to confirm more than
                one transaction.
              </p>
            </>
          ),
          btnText: "release",
          confirmationText: "RELEASE",
        })
      }

      if (gt(withdrawableAmountGrantOnly, 0)) {
        const contractMethod = isManagedGrant
          ? managedGrantContractInstance.methods.withdraw()
          : grantContract.methods.withdraw(selectedGrant.id)
        const acceptTxMsg = showMessage({
          type: messageType.INFO,
          sticky: true,
          title: "Waiting for the transaction confirmation...",
        })
        let pendingMessage = { id: null }
        await contractMethod
          .send({ from: yourAddress })
          .on("transactionHash", (hash) => {
            closeMessage(acceptTxMsg)
            pendingMessage = showMessage({
              type: messageType.PENDING_ACTION,
              sticky: true,
              title: "",
              content: (
                <TransactionIsPendingMsgContent
                  txHash={hash}
                  title={"Transaction is pending"}
                />
              ),
            })
          })
        closeMessage(pendingMessage)
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Tokens have been successfully released",
        })
      }

      const escrowMethodName = isManagedGrant
        ? "withdrawToManagedGrantee"
        : "withdraw"
      for (const operator of escrowOperatorsToWithdraw) {
        const infoMessage = showMessage({
          type: messageType.INFO,
          sticky: true,
          title: "Waiting for the transaction confirmation...",
        })
        let pendingMessage = { id: null }
        await tokenStakingEscrow.methods[escrowMethodName](operator)
          .send({ from: yourAddress })
          .on("transactionHash", (hash) => {
            closeMessage(infoMessage)
            pendingMessage = showMessage({
              type: messageType.PENDING_ACTION,
              sticky: true,
              title: "",
              content: (
                <TransactionIsPendingMsgContent
                  txHash={hash}
                  title={"Transaction is pending"}
                />
              ),
            })
          })
        closeMessage(pendingMessage)
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content:
            "Tokens have been successfully released from a TokenStakingEscrow deposit.",
        })
      }
    } catch (error) {
      showMessage({
        type: messageType.ERROR,
        title: "Error",
        content: error.message,
      })
      throw error
    }
  }

  return (
    <>
      <div className="flex-1">
        <CircularProgressBars
          total={selectedGrant.amount}
          items={[
            {
              value: selectedGrant.unlocked,
              backgroundStroke: colors.bgSuccess,
              color: colors.primary,
              label: "Unlocked",
            },
            {
              value: selectedGrant.released,
              color: colors.secondary,
              backgroundStroke: colors.bgSecondary,
              radius: 48,
              label: "Released",
            },
          ]}
          withLegend
        />
      </div>
      <div
        className={`ml-2 mt-1 flex-1${
          selectedGrant.readyToRelease === "0" ? " self-start" : ""
        }`}
      >
        <h5 className="text-grey-70">unlocked</h5>
        <h4 className="text-grey-70">
          {displayAmount(selectedGrant.unlocked)}
        </h4>
        <div className="text-smaller text-grey-40">
          of {displayAmountWithMetricSuffix(selectedGrant.amount)} Total
        </div>
        {gt(selectedGrant.readyToRelease || 0, 0) && (
          <div className="mt-2">
            <div className="text-secondary text-small flex wrap">
              <span className="mr-1">
                {`${displayAmountWithMetricSuffix(
                  selectedGrant.readyToRelease
                )} Available`}
              </span>
              <SpeechBubbleTooltip text="Releasing tokens allows them to be withdrawn from a grant." />
            </div>
            {!hideReleaseTokensBtn && (
              <SubmitButton
                className="btn btn-sm btn-secondary"
                onSubmitAction={releaseTokens}
                withMessageActionIsPending={false}
              >
                release tokens
              </SubmitButton>
            )}
          </div>
        )}
      </div>
    </>
  )
}

export const TokenGrantStakedDetails = ({ selectedGrant, stakedAmount }) => {
  return (
    <>
      <div className="flex-1 self-center">
        <CircularProgressBars
          total={selectedGrant.amount}
          items={[
            {
              value: stakedAmount,
              color: colors.grey70,
              backgroundStroke: colors.grey10,
              label: "Staked",
            },
          ]}
          withLegend
        />
      </div>
      <div className="ml-2 mt-1 self-start flex-1">
        <h5 className="text-grey-70">staked</h5>
        <h4 className="text-grey-70">{displayAmount(stakedAmount)}</h4>
        <div className="text-smaller text-grey-40">
          of {displayAmountWithMetricSuffix(selectedGrant.amount)} Total
        </div>
      </div>
    </>
  )
}

export default TokenGrantOverview
