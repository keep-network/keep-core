import React, { useEffect, useState } from "react"
import OnlyIf from "../OnlyIf";
import * as Icons from "./../Icons"
import {colors} from "../../constants/colors";
import ProgressBar from "../ProgressBar";
import moment from "moment";
import {useSelector} from "react-redux";
import List from "../List";
import Divider from "../Divider";
import Button from "../Button";
import {KEEP} from "../../utils/token.utils";
import {
  getSamePercentageValue,
  shortenAddress
} from "../../utils/general.utils";
import WithdrawalInfo from "./WithdrawalInfo";
import {add} from "../../utils/arithmetics.utils";
import {useWeb3Address} from "../WithWeb3Context";

const infoBannerTitle = "The cooldown period is 21 days"

const getItems = (keepAmount) => {
  return [
    {
      label: `Add ${KEEP.displayAmount(keepAmount)} KEEP to your existing expired withdrawal`
    },
    {
      label: `Reset the 21 day cooldown period.`
    }
  ]
}

const IncreaseWithdrawalModal = ({
  pendingWithdrawalBalance,
  addedAmount,
  submitBtnText,
  onBtnClick,
  onCancel,
  className = "",
  transactionFinished = false,
}) => {
  const yourAddress = useWeb3Address()
  const [step, setStep] = useState(transactionFinished ? 2 : 1)
  const {
    totalValueLocked,
    covTotalSupply,
  } = useSelector((state) => state.coveragePool)

  const onSubmit = (values) => {
    if (step === 1) {
      setStep((prevStep) => prevStep + 1)
    } else if (step === 2) {
      onBtnClick()
    }
  }

  return (
    <ModalWithOverview
      className={`${className} withdraw-modal__main-container`}
      pendingWithdrawalBalance={pendingWithdrawalBalance}
      addedAmount={addedAmount}
    >
      <OnlyIf condition={step === 1}>
        <IncreaseWithdrawalModalStep1 addedAmount={addedAmount} onSubmit={onSubmit}/>
      </OnlyIf>
      <OnlyIf condition={step === 2}>
        <WithdrawalInfo
          transactionFinished={transactionFinished}
          containerTitle={"Your new withdrawal amount"}
          submitBtnText={"withdraw"}
          onBtnClick={onSubmit}
          onCancel={onCancel}
          amount={add(pendingWithdrawalBalance, addedAmount)}
          infoBannerTitle={infoBannerTitle}
        >
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Expired withdrawal &nbsp;</h4>
            <h4 className={"withdraw-modal__data__value text-grey-70"}>
              {KEEP.displayAmount(getSamePercentageValue(pendingWithdrawalBalance, covTotalSupply, totalValueLocked))} KEEP
            </h4>
          </div>
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Increase amount &nbsp;</h4>
            <h4 className={"withdraw-modal__data__value text-grey-70"}>
              {KEEP.displayAmount(getSamePercentageValue(addedAmount, covTotalSupply, totalValueLocked))} KEEP
            </h4>
          </div>
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Wallet &nbsp;</h4>
            <h4 className={"withdraw-modal__data__value text-grey-70"}>
              {shortenAddress(yourAddress)}
            </h4>
          </div>
        </WithdrawalInfo>
      </OnlyIf>

    </ModalWithOverview>
  )
}

const IncreaseWithdrawalModalStep1 = ({addedAmount, onSubmit}) => {
  const {
    totalValueLocked,
    covTotalSupply,
  } = useSelector((state) => state.coveragePool)

  const items = getItems(getSamePercentageValue(addedAmount, covTotalSupply, totalValueLocked))
  return (
    <div>
      <h3 className={"mb-1"}>Take note!</h3>
      <h4 className={"color-grey-70"}>Your expired withdrawal needs to be re-initiated. This withdrawal will:</h4>
      <List items={items} className="increase-withdrawal-modal-step1__list mt-1 mb-1">
        <List.Content className="bullets text-grey-70" />
      </List>
      <h4 className={"color-grey-70 mb-3"}>Do you want to continue?</h4>
      <Divider className="divider divider--tile-fluid" />
      <div className={"flex row center"}>
        <Button
          className="btn btn-lg btn-primary"
          onClick={onSubmit}
        >
          continue
        </Button>
        <span className="ml-2 text-link text-grey-70">
          Cancel
        </span>
      </div>
    </div>
  )
}

const ModalWithOverview = ({
  children,
  className = "",
  pendingWithdrawalBalance,
  addedAmount,
 }) => {
  const {
    totalValueLocked,
    covTotalSupply,
  } = useSelector((state) => state.coveragePool)

  return (
    <div className={`modal-with-overview__content-container ${className}`}>
      <div className={"modal-with-overview-modal__info"}>
        {children}
      </div>
      <div className={"modal-with-overview__overview-container"}>
        <h4 className={"mb-1"}>Overview</h4>
        <IncreaseWithdrawalModal.Tile title={"expired withdrawal"} amount={pendingWithdrawalBalance} expired/>
        <h4 className={"modal-with-overview__added-amount color-grey-70"}>
          <Icons.ArrowDown />
          <Icons.Add />
          {KEEP.displayAmount(getSamePercentageValue(addedAmount, covTotalSupply, totalValueLocked))} KEEP
        </h4>
        <IncreaseWithdrawalModal.Tile title={"new withdrawal"} amount={add(pendingWithdrawalBalance, addedAmount)}/>
      </div>
    </div>
  )
}

const IncreaseWithdrawalModalTile = ({title, amount, expired = false}) => {
  const {
    totalValueLocked,
    covTotalSupply,
    withdrawalDelay,
  } = useSelector((state) => state.coveragePool)

  const endOfWithdrawalDate = moment().add(withdrawalDelay, "days")
  return (
    <div className={"modal-with-overview__tile"}>
      <h5 className={"modal-with-overview__tile-title"}>{title}</h5>
      <div className={"modal-with-overview__withdrawal-info"}>
        <h4 className={"modal-with-overview__amount text-grey-70"}>
          {KEEP.displayAmount(getSamePercentageValue(amount, covTotalSupply, totalValueLocked))} KEEP
        </h4>
        <OnlyIf condition={!expired}>
          <div className={"modal-with-overview__delay text-grey-50"}>21 days: {endOfWithdrawalDate.format("MM/DD")}</div>
        </OnlyIf>
        <OnlyIf condition={expired}>
          <div className={"modal-with-overview__delay text-error"}>expired</div>
        </OnlyIf>

      </div>
      <ProgressBar
        value={expired ? 100 : 0}
        total={100}
        color={colors.error}
        bgColor={colors.yellow30}
      >
        <ProgressBar.Inline
          height={10}
          className={"modal-with-overview__progress-bar"}
        />
      </ProgressBar>
    </div>
  )
}

IncreaseWithdrawalModal.Tile = IncreaseWithdrawalModalTile


export default IncreaseWithdrawalModal