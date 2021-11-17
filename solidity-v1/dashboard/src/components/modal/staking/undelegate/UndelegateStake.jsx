import React from "react"
import moment from "moment"
import { ModalHeader, ModalBody, ModalFooter } from "../../Modal"
import { InfoList } from "./components"
import List from "../../../List"
import Button, { SubmitButton } from "../../../Button"
import { KEEP } from "../../../../utils/token.utils"
import { colors } from "../../../../constants/colors"
import { undelegateStake } from "../../../../actions/web3"
import { useDispatch } from "react-redux"
import { withBaseModal } from "../../withBaseModal"

const styles = {
  box: {
    marginTop: "2rem",
    border: `1px solid ${colors.grey30}`,
    borderRadius: "0.5rem",
    padding: "1rem",
  },
}

const UndelegateStakeComponent = ({
  undelegationPeriod,
  amount,
  authorizer,
  operator,
  beneficiary,
  onClose,
}) => {
  const dispatch = useDispatch()

  return (
    <>
      <ModalHeader>Undelegate</ModalHeader>
      <ModalBody>
        <h3 className="mb-1">Undelegate Stake</h3>
        <InfoList
          undelegationPeriod={undelegationPeriod}
          undelegatedAt={moment().unix()}
        />
        <List style={styles.box}>
          <List.Title className="mb-1 text-mint-100">
            {KEEP.displayAmountWithSymbol(amount)}
          </List.Title>
          <List.Content>
            <DelegationDataItem label="authorizeer" text={authorizer} />
            <DelegationDataItem label="operator" text={operator} />
            <DelegationDataItem label="beneficiary" text={beneficiary} />
          </List.Content>
        </List>
      </ModalBody>
      <ModalFooter>
        <SubmitButton
          className="btn btn-primary btn-lg mr-2"
          onSubmitAction={(awaitingPromise) =>
            dispatch(undelegateStake(operator, awaitingPromise))
          }
        >
          undelegate
        </SubmitButton>
        <Button className="btn btn-unstyled" onClick={onClose}>
          Cancel
        </Button>
      </ModalFooter>
    </>
  )
}

const DelegationDataItem = ({ label, text }) => {
  return (
    <List.Item className="flex row center">
      <div className="text-label text-grey-70 mr-a">{label}</div>
      <div className="text-small text-grey-60">{text}</div>
    </List.Item>
  )
}

export const UndelegateStake = withBaseModal(UndelegateStakeComponent)
