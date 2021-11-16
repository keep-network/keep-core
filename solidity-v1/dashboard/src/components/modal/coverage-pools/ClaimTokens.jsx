import React from "react"
import { useDispatch } from "react-redux"
import { ModalBody, ModalFooter } from "../Modal"
import TokenAmount from "../../TokenAmount"
import List from "../../List"
import OnlyIf from "../../OnlyIf"
import { ViewInBlockExplorer } from "../../ViewInBlockExplorer"
import { covKEEP, KEEP } from "../../../utils/token.utils"
import { shortenAddress } from "../../../utils/general.utils"
import { withTimeline } from "./withTimeline"
import { COV_POOL_TIMELINE_STEPS } from "../../../constants/constants"
import { SubmitButton } from "../../Button"
import Button from "../../Button"
import { claimTokensFromWithdrawal } from "../../../actions/coverage-pool"

const ClaimTokensComponent = ({
  covAmount,
  collateralTokenAmount,
  address,
  onClose,
  transactionHash = null,
}) => {
  const dispatch = useDispatch()
  return (
    <>
      <ModalBody>
        <h3 className="mb-1">
          {transactionHash ? "Success!" : "You are about to claim:"}
        </h3>
        <TokenAmount
          amount={transactionHash ? collateralTokenAmount : covAmount}
          token={transactionHash ? KEEP : covKEEP}
          withIcon={!!transactionHash}
        />
        <OnlyIf condition={transactionHash}>
          <p className="text-grey-70 mt-1">
            View your transaction&nbsp;
            <ViewInBlockExplorer
              type="tx"
              className="text-grey-70"
              id={transactionHash}
              text="here"
            />
            .
          </p>
        </OnlyIf>
        <List className="mt-2">
          <List.Content className="text-grey-50">
            <List.Item className="flex row center">
              <span className="mr-a">Initial Withdrawal</span>
              <span>{KEEP.displayAmountWithSymbol(collateralTokenAmount)}</span>
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a">Wallet</span>
              <span>{shortenAddress(address)}</span>
            </List.Item>
          </List.Content>
        </List>
      </ModalBody>
      <ModalFooter>
        <OnlyIf condition={!transactionHash}>
          <SubmitButton
            className="btn btn-primary btn-lg mr-2"
            type="submit"
            onSubmitAction={(awaitingPromise) => {
              dispatch(claimTokensFromWithdrawal(awaitingPromise))
            }}
          >
            claim
          </SubmitButton>
        </OnlyIf>
        <Button
          className={`btn btn-${
            transactionHash ? "secondary btn-lg" : "unstyled text-link"
          }`}
          onClick={onClose}
        >
          {transactionHash ? "close" : "Cancel"}
        </Button>
      </ModalFooter>
    </>
  )
}

export const ClaimTokens = withTimeline({
  title: "Claim",
  step: COV_POOL_TIMELINE_STEPS.CLAIM_TOKENS,
  withDescription: true,
})(ClaimTokensComponent)
