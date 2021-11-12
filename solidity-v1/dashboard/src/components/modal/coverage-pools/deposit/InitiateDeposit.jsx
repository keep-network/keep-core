import React from "react"
import TokenAmount from "../../../TokenAmount"

import OnlyIf from "../../../OnlyIf"
import { covKEEP, KEEP } from "../../../../utils/token.utils"
import Button from "../../../Button"
import { ModalBody, ModalFooter } from "../../Modal"
import { ViewInBlockExplorer } from "../../../ViewInBlockExplorer"
import { useAcceptTermToConfirmFormik } from "../../../../hooks/useAcceptTermToConfirmFormik"
import { FormCheckboxBase } from "../../../FormCheckbox"
import { SubmitButton } from "../../../Button"
import { COV_POOL_TIMELINE_STEPS, LINK } from "../../../../constants/constants"
import { withTimeline } from "../withTimeline"
import List from "../../../List"
import { CoveragePoolV1ExchangeRate } from "../../../coverage-pools/ExchangeRate"
import { useDispatch } from "react-redux"
import { depositAssetPool } from "../../../../actions/coverage-pool"

const InitiateDepositComponent = ({
  amount, // amount of KEEP that user wants to deposit (in KEEP)
  covKEEPReceived, // total balance of the user after the deposit is done (in covKEEP)
  estimatedBalanceAmountInKeep, // estimated total balance of user in KEEP
  onClose,
  totalValueLocked,
  covTotalSupply,
  transactionHash = null,
}) => {
  const formik = useAcceptTermToConfirmFormik()
  const dispatch = useDispatch()

  return (
    <>
      <ModalBody>
        <h3 className="mb-1">
          {transactionHash ? "Success!" : "You are about to deposit:"}
        </h3>
        <TokenAmount amount={amount} token={KEEP} withIcon />
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
              <span className="mr-a">Exchange Rate</span>
              <CoveragePoolV1ExchangeRate
                covToken={covKEEP}
                collateralToken={KEEP}
                covTotalSupply={covTotalSupply}
                totalValueLocked={totalValueLocked}
              />
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a">Your Pool Balance</span>
              <span>
                {KEEP.displayAmountWithSymbol(estimatedBalanceAmountInKeep)}
              </span>
            </List.Item>
            <OnlyIf condition={transactionHash}>
              <List.Item className="flex row center">
                <span className="mr-a">CovKEEP Received</span>
                {covKEEP.displayAmountWithSymbol(covKEEPReceived)}
              </List.Item>
            </OnlyIf>
          </List.Content>
        </List>
      </ModalBody>
      <ModalFooter>
        <OnlyIf condition={!transactionHash}>
          <form className="mb-1">
            <FormCheckboxBase
              name="checked"
              type="checkbox"
              onChange={formik.handleChange}
              checked={formik.values.checked}
            >
              I confirm that I have read the{" "}
              <a
                href={LINK.coveragePools.docs}
                className="text-link text-black"
                rel="noopener noreferrer"
                target="_blank"
              >
                coverage pool documentation
              </a>{" "}
              and understand the risks.
            </FormCheckboxBase>
          </form>
          <SubmitButton
            className="btn btn-primary btn-lg mr-2"
            type="submit"
            onSubmitAction={(awaitingPromise) => {
              dispatch(depositAssetPool(amount, awaitingPromise))
            }}
            disabled={!(formik.isValid && formik.dirty)}
          >
            deposit
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

export const InitiateDeposit = withTimeline({
  title: "Deposit",
  step: COV_POOL_TIMELINE_STEPS.DEPOSITED_TOKENS,
  withDescription: true,
})(InitiateDepositComponent)
