import React from "react"
import { useDispatch } from "react-redux"
import { ModalBody, ModalFooter } from "../Modal"
import TokenAmount from "../../TokenAmount"
import OnlyIf from "../../OnlyIf"
import Button from "../../Button"
import { FormCheckboxBase } from "../../FormCheckbox"
import { SubmitButton } from "../../Button"
import List from "../../List"
import { CoveragePoolV1ExchangeRate } from "../../coverage-pools/ExchangeRate"
import { withTimeline } from "./withTimeline"
import { useAcceptTermToConfirmFormik } from "../../../hooks/useAcceptTermToConfirmFormik"
import { withdrawAssetPool } from "../../../actions/coverage-pool"
import { covKEEP, KEEP } from "../../../utils/token.utils"
import { COV_POOL_TIMELINE_STEPS, LINK } from "../../../constants/constants"
import { Keep } from "../../../contracts"

const InitiateWithdrawComponent = ({
  amount, // amount of covKEEP that user wants to withdraw
  covBalanceOf,
  totalValueLocked,
  covTotalSupply,
  onClose,
  isReinitialization = false,
  transactionHash = null,
}) => {
  const formik = useAcceptTermToConfirmFormik()
  const dispatch = useDispatch()

  return (
    <>
      <ModalBody>
        <h3 className="mb-1">
          {transactionHash ? "Almost there..." : "You are about to withdraw:"}
        </h3>
        <TokenAmount amount={amount} token={covKEEP} />
        <TokenAmount
          amount={Keep.coveragePoolV1.estimatedBalanceFor(
            amount,
            covTotalSupply,
            totalValueLocked
          )}
          token={KEEP}
          amountClassName="text-grey-60"
          symbolClassName="text-grey-60"
        />
        <p className="mt-1 text-grey-70">
          {transactionHash
            ? "After the 21 day cooldown you can claim your tokens in the dashboard."
            : "The withrawal initiation requires two transactions – an approval and a confirmation."}
        </p>
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
              <span>{covKEEP.displayAmountWithSymbol(covBalanceOf)}</span>
            </List.Item>
            {/* TODO: Display estimated gas cost */}
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
              dispatch(
                withdrawAssetPool(
                  isReinitialization ? 0 : amount,
                  awaitingPromise
                )
              )
            }}
            disabled={!(formik.isValid && formik.dirty)}
          >
            withdraw
          </SubmitButton>
        </OnlyIf>
        <Button
          className={`btn btn-${
            transactionHash ? "secondary btn-lg" : "unstyled text-link"
          }`}
          onClick={onClose}
        >
          {transactionHash ? "Close" : "Cancel"}
        </Button>
      </ModalFooter>
    </>
  )
}

export const InitiateWithdraw = withTimeline({
  title: "Withdraw",
  step: COV_POOL_TIMELINE_STEPS.WITHDRAW_DEPOSIT,
  withDescription: true,
})(InitiateWithdrawComponent)

export const WithdrawInitialized = withTimeline({
  title: "Withdraw",
  step: COV_POOL_TIMELINE_STEPS.COOLDOWN,
})(InitiateWithdrawComponent)
