import React, { useContext } from "react"
import { formatDate, displayAmount } from "../utils/general.utils"
import { SubmitButton } from "./Button"
import { colors } from "../constants/colors"
import { CircularProgressBars } from "./CircularProgressBar"
import { Web3Context } from "./WithWeb3Context"
import { useShowMessage, messageType } from "./Message"
import moment from "moment"
import { gt } from "../utils/arithmetics.utils"

const TokenGrantOverview = ({ selectedGrant }) => {
  const { yourAddress, grantContract } = useContext(Web3Context)
  const showMessage = useShowMessage()
  const cliffPeriod = moment
    .unix(selectedGrant.cliff)
    .from(moment.unix(selectedGrant.start), true)
  const fullyUnlockedDate = moment
    .unix(selectedGrant.start)
    .add(selectedGrant.duration, "seconds")

  const releaseTokens = async (onTransactionHashCallback) => {
    try {
      await grantContract.methods
        .withdraw(selectedGrant.id)
        .send({ from: yourAddress })
        .on("transactionHash", onTransactionHashCallback)
      showMessage({
        type: messageType.SUCCESS,
        title: "Success",
        content: "Tokens have been successfully released",
      })
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
    <div className="token-grant-overview">
      <h2 className="balance">
        {displayAmount(selectedGrant.amount)}&nbsp;KEEP
      </h2>
      <div className="text-small text-grey-40">
        Issued: {formatDate(moment.unix(selectedGrant.start))}
        <span className="text-smaller text-grey-30">
          &nbsp;Cliff: {cliffPeriod}
        </span>
        <br />
        Fully Unlocked: {formatDate(fullyUnlockedDate)}
      </div>
      <hr />
      <div className="flex">
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
          <div className="text-label">unlocked</div>
          <div className="text-label">
            {displayAmount(selectedGrant.unlocked)}
            <div className="text-smaller text-grey-40">
              of {displayAmount(selectedGrant.amount)} total
            </div>
          </div>
          {gt(selectedGrant.readyToRelease || 0, 0) && (
            <div className="mt-2">
              <div className="text-secondary text-small">
                {`${displayAmount(selectedGrant.readyToRelease)} Available`}
              </div>
              <SubmitButton
                className="btn btn-sm btn-secondary"
                onSubmitAction={releaseTokens}
              >
                release tokens
              </SubmitButton>
            </div>
          )}
        </div>
      </div>
      <div className="flex mt-1">
        <div className="flex-1 self-center">
          <CircularProgressBars
            total={selectedGrant.amount}
            items={[
              {
                value: selectedGrant.staked,
                backgroundStroke: "#F8E9D3",
                color: colors.brown,
                label: "Staked",
              },
            ]}
            withLegend
          />
        </div>
        <div className="ml-2 mt-1 self-start flex-1">
          <div className="text-label">staked</div>
          <div className="text-label">
            {displayAmount(selectedGrant.staked)}
            <div className="text-smaller text-grey-40">
              of {displayAmount(selectedGrant.amount)} total
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default TokenGrantOverview
