import React from "react"
import * as Icons from "../../Icons"
import TokenAmount from "../../TokenAmount"
import { ViewInBlockExplorer } from "../../ViewInBlockExplorer"
import List from "../../List"
import AddressShortcut from "../../AddressShortcut"
import OnlyIf from "../../OnlyIf"
import DoubleIcon from "../../DoubleIcon"
import Button from "../../Button"
import { TBTC } from "../../../utils/token.utils"
import { gt } from "../../../utils/arithmetics.utils"
import { LINK, TBTC_TOKEN_VERSION } from "../../../constants/constants"
import commonStyles from "../../tbtc-migration/styles"
import { ModalHeader } from "../Modal"
import { ModalBody } from "../Modal"
import { ModalFooter } from "../Modal"
import { withBaseModal } from "../withBaseModal"

const styles = {
  successIcon: {
    width: 20,
    height: 20,
  },
}

export const MigrationCompleted = withBaseModal(
  ({ from, to, amount, txHash, address, onClose, fee = 0 }) => {
    return (
      <>
        <ModalHeader>
          {to === TBTC_TOKEN_VERSION.v2 ? "Upgrade" : "Downgrade"}
        </ModalHeader>
        <ModalBody>
          <h3 className="flex row center">
            <Icons.Success
              style={styles.successIcon}
              className="success-icon success-icon--green"
            />
            &nbsp; Success! Tokens&nbsp;
            {to === TBTC_TOKEN_VERSION.v2 ? "upgraded" : "downgraded"}.
          </h3>
          <h4 className="text-grey-70 mb-1">
            View your transaction&nbsp;
            <ViewInBlockExplorer
              text="here"
              type="tx"
              id={txHash}
              className="text-grey-70"
            />
          </h4>
          <TokenAmount
            token={TBTC}
            symbol="tBTCv2"
            iconProps={{ className: "tbtc-icon tbtc-icon--black" }}
            amount={amount}
            amountClassName="h2 text-black"
            symbolClassName="h3 text-black"
            withIcon
          />
          <div className="text-center mt-1 mb-1">
            <span className="text-smaller" style={commonStyles[to]}>
              {`tBTC${to}`}
            </span>
          </div>

          <List>
            <List.Content>
              <List.Item className="flex row">
                <div className="text-grey-50">
                  {to === TBTC_TOKEN_VERSION.v2 ? "Upgraded" : "Downgraded"}
                  &nbsp;Tokens
                </div>
                <div className="ml-a">
                  <TokenAmount
                    token={TBTC}
                    symbol={`tBTC${to}`}
                    amount={amount}
                    amountClassName="text-grey-70"
                    symbolClassName="text-grey-70"
                  />
                </div>
              </List.Item>
              <OnlyIf condition={gt(fee, 0) && from === TBTC_TOKEN_VERSION.v2}>
                <List.Item className="flex row">
                  <div className="text-grey-50">Fee</div>
                  <div className="ml-a text-grey-70">
                    <TokenAmount
                      token={TBTC}
                      symbol="tBTCv2"
                      amount={fee}
                      amountClassName="text-grey-70"
                      symbolClassName="text-grey-70"
                    />
                  </div>
                </List.Item>
              </OnlyIf>

              <List.Item className="flex row">
                <div className="text-grey-50">Wallet</div>
                <div className="ml-a text-grey-70">
                  <AddressShortcut address={address} />
                </div>
              </List.Item>
            </List.Content>
          </List>
          <div
            className="mt-1 bg-mint-10 flex row center"
            style={commonStyles.boxWrapper}
          >
            <DoubleIcon
              MainIcon={Icons.TBTC}
              SecondaryIcon={
                to === TBTC_TOKEN_VERSION.v2
                  ? Icons.SaddleWhite
                  : Icons.KeepBlackGreen
              }
            />
            <div>
              &nbsp;
              {to === TBTC_TOKEN_VERSION.v2 ? "tBTCv2/Saddle" : "tBTCv1/KEEP"}
            </div>
            <a
              href={
                to === TBTC_TOKEN_VERSION.v2
                  ? LINK.pools.saddle.tbtcV2
                  : LINK.pools.uniswap.tbtcETH
              }
              rel="noopener noreferrer"
              target="_blank"
              className="btn btn-primary btn-md ml-a"
            >
              go to pool ↗
            </a>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button className="btn btn-secondary btn-lg" onClick={onClose}>
            Close
          </Button>
        </ModalFooter>
      </>
    )
  }
)
