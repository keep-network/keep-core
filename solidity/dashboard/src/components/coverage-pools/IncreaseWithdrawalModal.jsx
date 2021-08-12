import React, { useEffect, useState } from "react"
import OnlyIf from "../OnlyIf";
import * as Icons from "./../Icons"
import {colors} from "../../constants/colors";
import ProgressBar from "../ProgressBar";
import moment from "moment";
import {useSelector} from "react-redux";
import List from "../List";
import Divider from "../Divider";
import Button, {SubmitButton} from "../Button";
import {AcceptTermConfirmationModal} from "../ConfirmationModal";
import InitiateCovPoolsWithdrawModal from "./InitiateCovPoolsWithdrawModal";
import TokenAmount from "../TokenAmount";
import {covKEEP, KEEP} from "../../utils/token.utils";
import {shortenAddress} from "../../utils/general.utils";
import Banner from "../Banner";
import WithdrawalInfo from "./WithdrawalInfo";
import {add} from "../../utils/arithmetics.utils";

const infoBannerTitle = "The cooldown period is 21 days"

const items = [
  {
    label: "Add 1000 KEEP to your existing expired withdrawal"
  },
  {
    label: "Reset the 21 day cooldown period."
  }
]

const IncreaseWithdrawalModal = ({
  pendingWithdrawalBalance,
  addedAmount,
  submitBtnText,
  onBtnClick,
  onCancel,
  className = "",
  transactionFinished = false,
}) => {
  const [step, setStep] = useState(transactionFinished ? 2 : 1)

  const onSubmit = (values) => {
    if (step === 1) {
      setStep((prevStep) => prevStep + 1)
    } else if (step === 2) {
    }
  }

  return (
    <ModalWithOverview className={`${className} withdraw-modal__main-container` }>
      <OnlyIf condition={step === 1}>
        <IncreaseWithdrawalModalStep1 onSubmit={onSubmit}/>
      </OnlyIf>
      <OnlyIf condition={step === 2}>
        <WithdrawalInfo
          transactionFinished={transactionFinished}
          containerTitle={"Your new withdrawal amount"}
          submitBtnText={"withdraw"}
          onBtnClick={onBtnClick}
          onCancel={onCancel}
          amount={add(pendingWithdrawalBalance, addedAmount)}
          infoBannerTitle={infoBannerTitle}
        >
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Expired withdrawal &nbsp;</h4>
            <h4 className={"withdraw-modal__data__value text-grey-70"}>
              1,000 KEEP
            </h4>
          </div>
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Increase amount &nbsp;</h4>
            <h4 className={"withdraw-modal__data__value text-grey-70"}>
              1,000 KEEP
            </h4>
          </div>
          <div className={"withdraw-modal__data-row"}>
            <h4 className={"text-grey-50"}>Wallet &nbsp;</h4>
            <h4 className={"withdraw-modal__data__value text-grey-70"}>
              {shortenAddress("0x254673e7c7d76e051e80d30FCc3EA6A9C2a22222")}
            </h4>
          </div>
        </WithdrawalInfo>
      </OnlyIf>

    </ModalWithOverview>
  )
}

const IncreaseWithdrawalModalStep1 = ({onSubmit}) => {
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
 }) => {
  return (
    <div className={`modal-with-overview__content-container ${className}`}>
      <div className={"modal-with-overview-modal__info"}>
        {children}
      </div>
      <div className={"modal-with-overview__overview-container"}>
        <h4 className={"mb-1"}>Overview</h4>
        <IncreaseWithdrawalModal.Tile title={"expired withdrawal"} amount={"99000"} expired/>
        <h4 className={"modal-with-overview__added-amount color-grey-70"}>
          <Icons.ArrowDown />
          <Icons.Add />
          1,000 KEEP
        </h4>
        <IncreaseWithdrawalModal.Tile title={"new withdrawal"} amount={"100000"}/>
      </div>
    </div>
  )
}

const IncreaseWithdrawalModalTile = ({title, amount, expired = false}) => {
  const {
    withdrawalDelay,
  } = useSelector((state) => state.coveragePool)
  const endOfWithdrawalDate = moment().add(withdrawalDelay, "days")
  return (
    <div className={"modal-with-overview__tile"}>
      <h5 className={"modal-with-overview__tile-title"}>{title}</h5>
      <div className={"modal-with-overview__withdrawal-info"}>
        <h4 className={"modal-with-overview__amount text-grey-70"}>{amount} KEEEP</h4>
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