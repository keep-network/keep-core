import React from "react"
import List from "../List"
import TokenAmount from "../TokenAmount"
import * as Icons from "../Icons"
import Divider from "../Divider"
import Button from "../Button"
import { TBTC } from "../../utils/token.utils"
import { TBTC_TOKEN_VERSION } from "../../constants/constants"
import commonStyles from "./styles"

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

const ConfirmMigrationModal = ({
  from = TBTC_TOKEN_VERSION.v1,
  to = TBTC_TOKEN_VERSION.v2,
  amount,
  onBtnClick,
  onCancel,
}) => {
  return (
    <section>
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
      <div className="bg-grey-10 w-100 mb-1" style={commonStyles.boxWrapper}>
        <List items={remember}>
          <List.Title className="mb-1">Always remember:</List.Title>
          <List.Content className="bullets text-smaller" />
        </List>
      </div>

      <Divider className="divider divider--tile-fluid" />
      <Button
        className="btn btn-lg btn-primary"
        type="submit"
        onClick={onBtnClick}
      >
        {to === TBTC_TOKEN_VERSION.v2 ? "upgrade" : "downgrade"}
      </Button>
      <span onClick={onCancel} className="ml-1 text-link">
        Cancel
      </span>
    </section>
  )
}

export default ConfirmMigrationModal
