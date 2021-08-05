import React, {useState} from "react"
import OnlyIf from "../OnlyIf";
import InitiateCovPoolsWithdrawModal from "./InitiateCovPoolsWithdrawModal";
import TokenAmount from "../TokenAmount";
import {covKEEP, KEEP} from "../../utils/token.utils";
import AddAmountToWithdrawalForm from "./AddAmountToWithdrawalForm";
import {add} from "../../utils/arithmetics.utils";

const step1Title = "You are about to re-initiate this withdrawal:"
const step2Title = "You are about to re-withdraw:"

const ReinitiateWithdrawalModal = ({
 pendingWithdrawalBalance,
 covTokensAvailableToWithdraw,
 submitBtnText,
 onBtnClick,
 onCancel,
 className = "",
 transactionFinished = false,
}) => {
  const [step, setStep] = useState(1)
  const [amount, setAmount] = useState("0")

  const onSubmit = (values) => {
    if (step === 1) {
      setStep((preveStep) => preveStep + 1)
      setAmount(KEEP.fromTokenUnit(values.tokenAmount).toString())
    } else if (step === 2) {
    }
  }

  return (
    <>
      <OnlyIf condition={step === 1}>
        <ReinitiateWithdrawalModalStep1
          containerTitle={step1Title}
          pendingWithdrawalBalance={pendingWithdrawalBalance}
          covTokensAvailableToWithdraw={covTokensAvailableToWithdraw}
          submitBtnText={submitBtnText}
          onBtnClick={onSubmit}
          onCancel={onCancel}
          transactionFinished={transactionFinished}
        />
      </OnlyIf>
      <OnlyIf condition={step === 2}>
        <InitiateCovPoolsWithdrawModal
          amount={add(pendingWithdrawalBalance, amount)}
          containerTitle={step2Title}
          submitBtnText={"withdraw"}
          onBtnClick={onBtnClick}
          onCancel={onCancel}
          className={"reinitiate-withdrawal-modal__main-container"}
          transactionFinished={false}
        />
      </OnlyIf>
    </>
  )
}

const ReinitiateWithdrawalModalStep1 = ({
  containerTitle,
  pendingWithdrawalBalance,
  covTokensAvailableToWithdraw,
  submitBtnText,
  onBtnClick,
  onCancel,
  transactionFinished,
}) => {
  return (
    <div className={"reinitiate-withdrawal-modal"}>
      <h3 className={"reinitiate-withdrawal-modal__container-title"}>{containerTitle}</h3>
      <div className={"reinitiate-withdrawal-modal__data"}>
        <TokenAmount
          amount={pendingWithdrawalBalance}
          wrapperClassName={"reinitiate-withdrawal-modal__token-amount"}
          token={KEEP}
          withIcon
        />
        <TokenAmount
          wrapperClassName={"reinitiate-withdrawal-modal__cov-token-amount"}
          amount={pendingWithdrawalBalance}
          amountClassName={"h4 text-grey-60"}
          symbolClassName={"h4 text-grey-60"}
          token={covKEEP}
        />
      </div>
      <AddAmountToWithdrawalForm onSubmit={onBtnClick} tokenAmount={covTokensAvailableToWithdraw} />
    </div>
  )
}

export default ReinitiateWithdrawalModal