import React from "react"
import { withFormik } from "formik"
import FormInput from "../../components/FormInput"
import { SubmitButton } from "../../components/Button"
import Divider from "../../components/Divider"
import MaxAmountAddon from "../MaxAmountAddon"
import { normalizeAmount, formatAmount } from "../../forms/form.utils"
import { KEEP } from "../../utils/token.utils"
import List from "../List"
import * as Icons from "../Icons"
import Chip from "../Chip"
import TokenAmount from "../TokenAmount"

const DepositForm = ({
  tokenAmount,
  onSubmit,
  estimatedRewards = [
    { apy: 10, label: "Weekly", reward: "1000000000000000000000" },
    { apy: 20, label: "Monthly", reward: "1000000000000000000000" },
    { apy: 30, label: "Yearly", reward: "1000000000000000000000" },
  ],
}) => {
  return (
    <form className="deposit-form">
      <FormInput
        name="tokenAmount"
        type="text"
        label="Amount"
        placeholder="0"
        normalize={normalizeAmount}
        format={formatAmount}
        inputAddon={
          <MaxAmountAddon
            onClick={() => console.log("on clikc addon")}
            text="Max Stake"
          />
        }
        additionalInfoText={`KEEP Balance ${KEEP.displayAmountWithSymbol(
          tokenAmount
        )}`}
      />
      <List>
        <List.Title className="mb-2">Estimated Rewards</List.Title>
        <List.Content>{estimatedRewards.map(renderListItem)}</List.Content>
      </List>
      <Divider className="divider divider--tile-full-width" />

      <p>
        Risk warning:&nbsp;
        <a
          href="https://example.com"
          rel="noopener noreferrer"
          target="_blank"
          className="text-black"
        >
          Read the documentation
        </a>
      </p>
      <SubmitButton
        className="btn btn-lg btn-primary w-100"
        onSubmitAction={onSubmit}
      >
        deposit
      </SubmitButton>
    </form>
  )
}

const renderListItem = (item) => (
  <EstimatedAPYListItem key={item.label} {...item} />
)

const EstimatedAPYListItem = ({ apy, reward, label }) => {
  return (
    <List.Item className="mb-1">
      <div className="flex row center">
        <Icons.Time
          className="time-icon time-icon--grey-50"
          width={16}
          height={16}
        />
        &nbsp;
        <span className="text-grey-50">{label}</span>
        &nbsp;
        <Chip text={`${apy}% APY`} size="small" />
        <TokenAmount
          wrapperClassName="ml-a"
          amount={reward}
          amountClassName=""
          symbolClassName=""
        />
      </div>
    </List.Item>
  )
}

export default withFormik({
  validateOnChange: false,
  validateOnBlur: false,
  mapPropsToValues: () => ({
    tokenAmount: "0",
  }),
  validate: (values, props) => {
    return {}
  },
  displayName: "CovPoolsDepositForm",
})(DepositForm)
