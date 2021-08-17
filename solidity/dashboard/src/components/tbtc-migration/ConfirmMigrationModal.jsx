import React from "react"
import List from "../List"
import TokenAmount from "../TokenAmount"
import * as Icons from "../Icons"
import Divider from "../Divider"
import Button from "../Button"
import { TBTC } from "../../utils/token.utils"

const remember = [
  { label: "the upgrade is reversible" },
  { label: "the upgrade and downgrade cost zero fee on Keep Network" },
  {
    label:
      "for liquidations or redemptions for tBTC bridge, only tBTC v1 is accepted",
  },
]

const swapBoxStyle = {
  padding: "0.5rem 0.75rem",
  borderRadius: "0.5rem",
  height: "40px",
  marginRight: "0.5rem",
}
const styles = {
  swapBox: swapBoxStyle,
  v1: {
    color: "white",
    backgroundColor: "black",
    textAlign: "center",
    ...swapBoxStyle,
  },
  v2: {
    color: "black",
    backgroundColor: "white",
    border: "1px solid black",
    textAlign: "center",

    ...swapBoxStyle,
  },
}

const ConfirmMigrationModal = ({
  from = "v1",
  to = "v2",
  amount,
  onBtnClick,
  onCancel,
}) => {
  return (
    <section>
      <h3 className="mb-1">{`You are about to ${
        to === "v2" ? "upgrade" : "downgrade"
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
        <div className="text-smaller" style={styles[from]}>
          {`tBTC${from}`}
        </div>
        <div className="bg-violet-10 flex row" style={styles.swapBox}>
          <Icons.ChevronRight
            className="chevron-right-icon chevron-right-icon--secondary"
            style={{ marginRight: "0.5rem" }}
          />
          <Icons.ChevronRight
            className="chevron-right-icon chevron-right-icon--secondary"
            style={{ marginRight: "0.5rem" }}
          />
          <Icons.ChevronRight
            className="chevron-right-icon chevron-right-icon--secondary"
            style={{ marginRight: "0.5rem" }}
          />
          <Icons.ChevronRight
            className="chevron-right-icon chevron-right-icon--secondary"
            style={{ marginRight: "0.5rem" }}
          />
          <Icons.ChevronRight className="chevron-right-icon chevron-right-icon--secondary" />
        </div>
        <div className="text-smaller" style={styles[to]}>
          {`tBTC${to}`}
        </div>
      </div>
      <div
        className="bg-grey-10 w-100 mb-1"
        style={{ borderRadius: "8px", padding: "1rem" }}
      >
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
        {to === "v2" ? "upgrade" : "downgrade"}
      </Button>
      <span onClick={onCancel} className="ml-1 text-link">
        Cancel
      </span>
    </section>
  )
}

export default ConfirmMigrationModal
