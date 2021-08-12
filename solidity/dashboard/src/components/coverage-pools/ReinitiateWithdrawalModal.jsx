import React, { useEffect, useState } from "react"
import OnlyIf from "../OnlyIf";
import InitiateCovPoolsWithdrawModal from "./InitiateCovPoolsWithdrawModal";
import TokenAmount from "../TokenAmount";
import {covKEEP, KEEP} from "../../utils/token.utils";
import AddAmountToWithdrawalForm from "./AddAmountToWithdrawalForm";
import {gt, eq, lt} from "../../utils/arithmetics.utils";
import { useDispatch } from "react-redux";
import { EVENTS } from "../../constants/events";
import IncreaseWithdrawalModal from "./IncreaseWithdrawalModal";

const step1Title = "You are about to re-initiate this withdrawal:"
const step2Title = "You are about to re-withdraw:"

const ReinitiateWithdrawalModal = ({
 pendingWithdrawalBalance,
 covTokensAvailableToWithdraw,
 submitBtnText,
 onBtnClick,
 onCancel,
 className = "",
 initialAmountValue = "0", // in smallest unit
 transactionFinished = false,
}) => {
  const [step, setStep] = useState(1)
  const [amount, setAmount] = useState(initialAmountValue.toString())
  const dispatch = useDispatch()

  const onSubmit = (values) => {
    if (step === 1) {
      setStep((prevStep) => prevStep + 1)
      setAmount(KEEP.fromTokenUnit(values.tokenAmount).toString())
    } else if (step === 2) {
    }
  }

  useEffect(() => {
    if (step === 2) {
      if (eq(amount, 0)) {
        dispatch({
          type: "modal/set_emitted_event",
          payload: EVENTS.COVERAGE_POOLS.RE_WITHDRAWAL_INITIATED
        })
      } else if (gt(amount, 0)) {
        dispatch({
          type: "modal/set_emitted_event",
          payload: EVENTS.COVERAGE_POOLS.ADD_BALANCE_TO_WITHDRAWAL
        })
      }
    }
  }, [step, amount, amount])

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
          initialAmountValue={initialAmountValue}
        />
      </OnlyIf>
      <OnlyIf condition={step === 2}>
        <OnlyIf condition={eq(amount, 0)}>
          <InitiateCovPoolsWithdrawModal
            amount={pendingWithdrawalBalance}
            containerTitle={step2Title}
            submitBtnText={"withdraw"}
            onBtnClick={onBtnClick}
            onCancel={onCancel}
            className={"reinitiate-withdrawal-modal__main-container"}
            transactionFinished={false}
          />
        </OnlyIf>
        <OnlyIf condition={lt(amount, 0)}>
          <InitiateCovPoolsWithdrawModal
            amount={pendingWithdrawalBalance}
            containerTitle={step2Title}
            submitBtnText={"withdraw"}
            onBtnClick={onBtnClick}
            onCancel={onCancel}
            className={"reinitiate-withdrawal-modal__main-container"}
            transactionFinished={false}
          />
        </OnlyIf>
        <OnlyIf condition={gt(amount, 0)}>
          <IncreaseWithdrawalModal
            pendingWithdrawalBalance={pendingWithdrawalBalance}
            addedAmount={amount}
            submitBtnText={"withdraw"}
            onBtnClick={onBtnClick}
            onCancel={onCancel}
            className={"increase-withdrawal-modal__main-container"}
            transactionFinished={false}
          />
        </OnlyIf>
      </OnlyIf>
    </>
  )
}

const ReinitiateWithdrawalModalStep1 = ({
  containerTitle,
  pendingWithdrawalBalance,
  covTokensAvailableToWithdraw,
  initialAmountValue,
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
      <AddAmountToWithdrawalForm
        initialValue={initialAmountValue}
        onSubmit={onBtnClick}
        tokenAmount={covTokensAvailableToWithdraw} />
    </div>
  )
}

export default ReinitiateWithdrawalModal