import React from "react"
import { useDispatch } from "react-redux"
import List from "../../List"
import TokenAmount from "../../TokenAmount"
import * as Icons from "../../Icons"
import Button from "../../Button"
import { TBTC } from "../../../utils/token.utils"
import { TBTC_TOKEN_VERSION } from "../../../constants/constants"
import commonStyles from "../../tbtc-migration/styles"
import { ModalHeader, ModalBody, ModalFooter } from "../Modal"
import { SubmitButton } from "../../Button"
import { tbtcV2Migration } from "../../../actions"
import { withBaseModal } from "../withBaseModal"

const remember = [
  { label: "the upgrade is reversible" },
  { label: "the upgrade and downgrade cost zero fee on Keep Network" },
  {
    label:
      "for liquidations or redemptions for tBTC bridge, only tBTC v1 is accepted",
  },
]

const styles = {
  chevronRighIcon: {
    marginRight: "0.5rem",
  },
}

export const ConfirmMigration = withBaseModal(
  ({
    from = TBTC_TOKEN_VERSION.v1,
    to = TBTC_TOKEN_VERSION.v2,
    amount,
    onClose,
  }) => {
    const dispatch = useDispatch()

    const onSubmit = async (awaitingPromise) => {
      if (to === TBTC_TOKEN_VERSION.v2) {
        dispatch(tbtcV2Migration.mint(amount, awaitingPromise))
      } else {
        dispatch(tbtcV2Migration.unmint(amount, awaitingPromise))
      }
    }

    return (
      <>
        <ModalHeader>
          {to === TBTC_TOKEN_VERSION.v2 ? "Upgrade" : "Downgrade"}
        </ModalHeader>
        <ModalBody>
          <h3 className="mb-1">{`You are about to ${
            to === TBTC_TOKEN_VERSION.v2 ? "upgrade" : "downgrade"
          }`}</h3>
          <TokenAmount
            token={TBTC}
            symbol={`tBTC${from}`}
            iconProps={{ className: "tbtc-icon tbtc-icon--black" }}
            amount={amount}
            amountClassName="h2 text-black"
            symbolClassName="h3 text-black"
            withIcon
          />
          <div className="flex row full-center mt-2 mb-1">
            <div className="text-smaller" style={commonStyles[from]}>
              {`tBTC${from}`}
            </div>
            <div className="bg-violet-10 flex row" style={commonStyles.swapBox}>
              <Icons.ChevronRight
                className="chevron-right-icon chevron-right-icon--secondary"
                style={styles.chevronRighIcon}
              />
              <Icons.ChevronRight
                className="chevron-right-icon chevron-right-icon--secondary"
                style={styles.chevronRighIcon}
              />
              <Icons.ChevronRight
                className="chevron-right-icon chevron-right-icon--secondary"
                style={styles.chevronRighIcon}
              />
              <Icons.ChevronRight
                className="chevron-right-icon chevron-right-icon--secondary"
                style={styles.chevronRighIcon}
              />
              <Icons.ChevronRight className="chevron-right-icon chevron-right-icon--secondary" />
            </div>
            <div className="text-smaller" style={commonStyles[to]}>
              {`tBTC${to}`}
            </div>
          </div>
          <div
            className="bg-grey-10 w-100 mb-1"
            style={commonStyles.boxWrapper}
          >
            <List items={remember}>
              <List.Title className="mb-1">Always remember:</List.Title>
              <List.Content className="bullets text-smaller" />
            </List>
          </div>
        </ModalBody>
        <ModalFooter>
          <SubmitButton
            className="btn btn-lg btn-primary mr-1"
            type="submit"
            onSubmitAction={onSubmit}
          >
            {to === TBTC_TOKEN_VERSION.v2 ? "upgrade" : "downgrade"}
          </SubmitButton>
          <Button className="btn btn-unstyled" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
