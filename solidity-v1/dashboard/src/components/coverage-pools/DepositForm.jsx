import React from "react"
import { withFormik } from "formik"
import FormInput from "../../components/FormInput"
import Button from "../../components/Button"
import Divider from "../../components/Divider"
import MaxAmountAddon from "../MaxAmountAddon"
import {
  formatFloatingAmount,
  normalizeFloatingAmount,
} from "../../forms/form.utils"
import { KEEP } from "../../utils/token.utils"
import List from "../List"
import * as Icons from "../Icons"
import Chip from "../Chip"
import TokenAmount from "../TokenAmount"
import {
  validateAmountInRange,
  getErrorsObj,
} from "../../forms/common-validators"
import { lte } from "../../utils/arithmetics.utils"
import useSetMaxAmountToken from "../../hooks/useSetMaxAmountToken"
import { displayPercentageValue } from "../../utils/general.utils"
import OnlyIf from "../OnlyIf"
import { LINK } from "../../constants/constants"

const DepositForm = ({ tokenAmount, apy, ...formikProps }) => {
  const onAddonClick = useSetMaxAmountToken(
    "tokenAmount",
    tokenAmount,
    KEEP,
    KEEP.decimals
  )

  const getEstimatedReward = () => {
    if (!formikProps.values.tokenAmount) {
      return null
    } else if (!isFinite(apy) || apy > 999) {
      return Infinity
    } else {
      return KEEP.fromTokenUnit(formikProps.values.tokenAmount)
        .multipliedBy(apy.toString())
        .toFixed()
        .toString()
    }
  }

  return (
    <form className="deposit-form">
      <div className="deposit-form__token-amount-wrapper">
        <FormInput
          name="tokenAmount"
          type="text"
          label="Amount"
          placeholder="0"
          normalize={normalizeFloatingAmount}
          format={formatFloatingAmount}
          leftIconComponent={
            <Icons.KeepOutline
              className="keep-outline--grey-60"
              width={20}
              height={20}
              style={{ margin: "0 1rem" }}
            />
          }
          inputAddon={
            <MaxAmountAddon onClick={onAddonClick} text="Max Amount" />
          }
          additionalInfoText={
            <>
              <span>Keep Balance</span>&nbsp;
              <TokenAmount
                amount={tokenAmount}
                wrapperClassName={"deposit-form__keep-balance-amount"}
                amountClassName={"text-success"}
                symbolClassName={"text-success"}
              />
            </>
          }
        />
      </div>
      <List>
        <List.Title className="mb-2">Estimated Rewards</List.Title>
        <List.Content>
          <EstimatedAPYListItem
            apy={apy}
            reward={getEstimatedReward()}
            label="Yearly"
          />
        </List.Content>
      </List>
      <Divider className="divider divider--tile-fluid" />
      <Button
        className="btn btn-lg btn-primary w-100"
        type="submit"
        onClick={formikProps.handleSubmit}
        disabled={!(formikProps.isValid && formikProps.dirty)}
      >
        deposit
      </Button>
      <p className="text-center text-secondary mt-1 mb-0">
        Risk warning:&nbsp;
        <a
          href={LINK.coveragePools.docs}
          rel="noopener noreferrer"
          target="_blank"
        >
          Read the documentation
        </a>
      </p>
    </form>
  )
}

const EstimatedAPYListItem = ({ apy, reward, label }) => {
  return (
    <List.Item className="mb-1">
      <div className="flex row center">
        <Icons.Time
          className="time-icon time-icon--grey-70"
          width={16}
          height={16}
        />
        &nbsp;
        <span className="text-grey-70">{label}</span>
        &nbsp;
        <Chip
          text={`${displayPercentageValue(apy * 100, false)}`}
          size="small"
          color="primary"
        />
        <OnlyIf condition={!reward}>
          <span className="text-grey-50 ml-a">Enter amount above</span>
        </OnlyIf>
        <OnlyIf condition={reward === Infinity}>
          <span className="text-grey-50 ml-a">∞</span>
        </OnlyIf>
        <OnlyIf condition={reward && reward !== Infinity}>
          <TokenAmount
            wrapperClassName="ml-a"
            amount={reward}
            amountClassName=""
            symbolClassName=""
          />
        </OnlyIf>
      </div>
    </List.Item>
  )
}

export default withFormik({
  validateOnChange: true,
  validateOnBlur: true,
  mapPropsToValues: () => ({
    tokenAmount: "0",
  }),
  validate: (values, props) => {
    const { tokenAmount } = values
    const errors = {}

    if (lte(props.tokenAmount || 0, 0)) {
      errors.tokenAmount = "Insufficient funds"
    } else {
      errors.tokenAmount = validateAmountInRange(
        tokenAmount,
        props.tokenAmount,
        KEEP.fromTokenUnit(1)
      )
    }

    return getErrorsObj(errors)
  },
  handleSubmit: (values, { props, resetForm }) => {
    props.onSubmit(values)
    resetForm({ tokenAmount: "0" })
  },
  displayName: "CovPoolsDepositForm",
})(DepositForm)
