import React from "react"
import { useDispatch } from "react-redux"
import { ModalBody, ModalFooter } from "../Modal"
import { withWithdrawalOverview } from "./withWithdrawalOverview"
import TokenAmount from "../../TokenAmount"
import Button from "../../Button"
import Banner from "../../Banner"
import * as Icons from "../../Icons"
import { FormCheckboxBase } from "../../FormCheckbox"
import { SubmitButton } from "../../Button"
import List from "../../List"
import { CoveragePoolV1ExchangeRate } from "../../coverage-pools/ExchangeRate"
import { covKEEP, KEEP } from "../../../utils/token.utils"
import { shortenAddress } from "../../../utils/general.utils"
import { add } from "../../../utils/arithmetics.utils"
import { LINK, PENDING_WITHDRAWAL_STATUS } from "../../../constants/constants"
import { Keep } from "../../../contracts"
import { withdrawAssetPool } from "../../../actions/coverage-pool"
import { useAcceptTermToConfirmFormik } from "../../../hooks/useAcceptTermToConfirmFormik"
import { colors } from "../../../constants/colors"

const BannerIcon = () => (
  <Icons.Tooltip color={colors.black} backgroundColor={colors.grey10} />
)
const IncreaseWithdrawalComponent = ({
  existingWithdrawalCovAmount,
  covAmountToAdd,
  address,
  withdrawalStatus,
  covBalanceOf,
  totalValueLocked,
  covTotalSupply,
  onClose,
}) => {
  const formik = useAcceptTermToConfirmFormik()
  const dispatch = useDispatch()
  const amount = add(existingWithdrawalCovAmount, covAmountToAdd)

  return (
    <>
      <ModalBody>
        <h3 className="mb-1">
          {withdrawalStatus === PENDING_WITHDRAWAL_STATUS.EXPIRED
            ? "Your new withdrawal amount:"
            : "You are about to withdraw:"}
        </h3>
        <TokenAmount amount={amount} token={covKEEP} />
        <TokenAmount
          amount={Keep.coveragePoolV1.estimatedBalanceFor(
            covBalanceOf,
            covTotalSupply,
            totalValueLocked
          )}
          token={KEEP}
          amountClassName="text-grey-60"
          symbolClassName="text-grey-60"
        />
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
              <span className="mr-a">
                {withdrawalStatus === PENDING_WITHDRAWAL_STATUS.EXPIRED
                  ? "Expired Withdrawal"
                  : "Existing Withdrawal"}
              </span>
              <span>
                {covKEEP.displayAmountWithSymbol(existingWithdrawalCovAmount)}
              </span>
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a">Increase Amount</span>
              <span>{covKEEP.displayAmountWithSymbol(covAmountToAdd)}</span>
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a">Wallet</span>
              <span>{shortenAddress(address)}</span>
            </List.Item>
          </List.Content>
        </List>

        <Banner
          style={{
            padding: "0.375rem 0.75rem",
            backgroundColor: colors.grey10,
            borderRadius: "100px",
            minWidth: "100%",
            boxShadow: "none",
          }}
          className="flex row center mt-2"
        >
          <Banner.Icon icon={BannerIcon} />
          <Banner.Description
            className="text-black text-bold"
            style={{ margin: "0", marginLeft: "0.5rem", fontSize: "0.875rem" }}
          >
            The cooldown period is 21 days.
          </Banner.Description>
        </Banner>
      </ModalBody>
      <ModalFooter>
        <form>
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
            dispatch(withdrawAssetPool(covAmountToAdd, awaitingPromise))
          }}
          disabled={!(formik.isValid && formik.dirty)}
        >
          withdraw
        </SubmitButton>
        <Button className="btn btn-unstyled text-link" onClick={onClose}>
          Cancel
        </Button>
      </ModalFooter>
    </>
  )
}

export const IncreaseWithdrawal = withWithdrawalOverview({
  title: "Increase Withdrawal",
})(IncreaseWithdrawalComponent)
