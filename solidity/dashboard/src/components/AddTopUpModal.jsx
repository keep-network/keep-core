import React, { useState } from "react"
import AmountForm from "./AmountForm"
import AddressShortcut from "./AddressShortcut"
import Tag from "./Tag"
import * as Icons from "./Icons"
import { withConfirmationModal } from "./ConfirmationModal"
import Button from "./Button"
import Divider from "./Divider"
import TokenAmount from "./TokenAmount"
import { KEEP } from "../utils/token.utils"
import { add } from "../utils/arithmetics.utils"
import { useModal } from "../hooks/useModal"

const AddTopUpModal = ({
  authorizerAddress,
  beneficiary,
  operatorAddress,
  currentAmount,
  onBtnClick: submitConfirmationModal,
  onCancel,
  availableAmount,
}) => {
  const [step, setStep] = useState(1)
  const [amount, setAmount] = useState("0")

  const onSubmit = (values) => {
    if (step === 1) {
      setAmount(values.amount)
      setStep((prevStep) => prevStep + 1)
    } else if (step === 2) {
      submitConfirmationModal({ amount })
    }
  }

  const amountInSmallestUnit = KEEP.fromTokenUnit(amount)

  return (
    <>
      {step === 1 && (
        <Step1
          authorizerAddress={authorizerAddress}
          beneficiary={beneficiary}
          operatorAddress={operatorAddress}
          currentAmount={currentAmount}
          onSubmit={onSubmit}
          availableAmount={availableAmount}
          onCancel={onCancel}
        />
      )}
      {step === 2 && (
        <Step2
          title={`You are adding ${KEEP.toFormat(
            amount
          )} KEEP to this delegation. Note that this top up cannot be canceled once initiated.`}
          confirmationText="CONFIRM"
          getLabelText={(confirmationText) =>
            `Type ${confirmationText} to add KEEP.`
          }
          btnText="confirm"
          onBtnClick={onSubmit}
          onCancel={onCancel}
          authorizerAddress={authorizerAddress}
          beneficiary={beneficiary}
          operatorAddress={operatorAddress}
          newAmount={add(amountInSmallestUnit, currentAmount)}
          amountToAdd={amountInSmallestUnit}
          onSubmit={onSubmit}
        />
      )}
    </>
  )
}

const Step1 = ({
  authorizerAddress,
  beneficiary,
  operatorAddress,
  currentAmount,
  availableAmount,
  onSubmit,
  onCancel,
}) => {
  return (
    <>
      <h3 className="mb-1">
        Enter an amount of KEEP to add to this existing delegation.
      </h3>
      <h4 className="text-grey-70 mb-1">Delegation Details</h4>
      <DelegationDetails
        authorizerAddress={authorizerAddress}
        beneficiary={beneficiary}
        operatorAddress={operatorAddress}
      />
      <div className="flex row center mt-2 mb-2">
        <div className="flex-1">
          <Tag text="Current" IconComponent={Icons.KeepToken} />
        </div>
        <TokenAmount amount={currentAmount} />
      </div>
      <AmountForm
        onCancel={onCancel}
        submitBtnText="add keep"
        availableAmount={availableAmount}
        currentAmount={currentAmount}
        onBtnClick={onSubmit}
      />
    </>
  )
}

const Step2Content = ({
  authorizerAddress,
  beneficiary,
  operatorAddress,
  amountToAdd,
  newAmount,
}) => {
  return (
    <>
      <TokenAmount
        amount={amountToAdd}
        withIcon
        icon={Icons.Plus}
        iconProps={{ width: 24, height: 24, className: "plus-icon" }}
      />
      <h4 className="text-grey-70 mb-1">
        New delegation balance: {KEEP.displayAmountWithSymbol(newAmount)}
      </h4>
      <DelegationDetails
        authorizerAddress={authorizerAddress}
        beneficiary={beneficiary}
        operatorAddress={operatorAddress}
      />
    </>
  )
}

const Step2 = withConfirmationModal(Step2Content)

const DelegationAddress = ({ address, label }) => {
  return (
    <div className="flex row center text-grey-70">
      <h5 className="flex-1">{label}</h5>
      <AddressShortcut address={address} classNames="h6 text-grey-70" />
    </div>
  )
}

const DelegationDetails = ({
  authorizerAddress,
  beneficiary,
  operatorAddress,
}) => (
  <>
    <DelegationAddress address={authorizerAddress} label="authorizer" />
    <DelegationAddress address={operatorAddress} label="operator" />
    <DelegationAddress address={beneficiary} label="beneficiary" />
  </>
)

const TopUpInitiatedConfirmationModal = ({
  addedAmount,
  currentAmount,
  authorizerAddress,
  beneficiary,
  operatorAddress,
}) => {
  const { closeModal } = useModal()
  const newAmount = add(currentAmount, addedAmount)

  return (
    <>
      <section className="text-center">
        <h3>Almost there...</h3>
        <div className="text-big text-grey-70">
          Return to the token dashboard <strong>in 12 hours</strong> to commit a
          top-up!
        </div>
      </section>
      <TokenAmount
        amount={addedAmount}
        withIcon
        icon={Icons.Plus}
        iconProps={{ width: 24, height: 24, className: "plus-icon" }}
      />
      <h4 className="text-grey-70 mb-1">
        New delegation balance: {KEEP.displayAmountWithSymbol(newAmount)}
      </h4>
      <DelegationDetails
        authorizerAddress={authorizerAddress}
        beneficiary={beneficiary}
        operatorAddress={operatorAddress}
      />
      <Divider
        style={{
          margin: "2rem -2rem 0",
          padding: "2rem 2rem 0",
        }}
      />
      <Button className="btn btn-lg btn-secondary" onClick={closeModal}>
        close
      </Button>
    </>
  )
}

export default AddTopUpModal

export { TopUpInitiatedConfirmationModal }
