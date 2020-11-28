import React from "react"
import * as Icons from "./Icons"
import TokenAmount from "./TokenAmount"
import { useState, useEffect } from "react"
import { TokenGrantSkeletonOverview } from "../components/skeletons/TokenOverviewSkeleton"
import moment from "moment"
import { LoadingOverlay } from "./Loadable"
import { SubmitButton } from "./Button"
import { useWeb3Context } from "./WithWeb3Context"
import { CircularProgressBarEth } from "./CircularProgressBar"
import {
  satsToTBtcViaWeitoshi,
  displayAmount,
  displayAmountHigherOrderFn
} from "../utils/token.utils"
import { colors } from "../constants/colors"

export const ArbitrageurTokenDetails = ({
  title = "Your Balance",
  tokenValue
}) => {

  return (
    <section className="token-grant-details">
      <div className="flex wrap center">
        <h3 className="text-grey-70">
          {title}
        </h3>
      </div>
      <TokenAmount
        currencyIcon={Icons.TBTC}
        amountClassName="h1 text-grey-70"
        suffixClassName="h2"
        displayWithMetricSuffix={false}
        // Hardcoded "4" since the smaller relevant decimal is 3 based on current lot sizes
        displayAmountFunction={displayAmountHigherOrderFn(false,4)}
        amount={tokenValue} />
    </section>
  )
}

export const DepositAuctionOverview = (props) => {
  // variables
  const {
    auctionOfferSummaryIsFetching,
    depositState,
    tBtcBalance,
    auctionValueBN,
    bondAmountWei,
    depositSizeSatoshis
  } = props
  // functions
  const {
    refreshData,
    getPercentageOnOffer,
    onLiquidateFromSummaryBtn
  } = props
  const web3Context = useWeb3Context()

  const [lastRefreshedMoment, setLastRefreshedMoment] = useState(moment())

  useEffect(
    () => {
      setTimeout(
        () => {
          refreshData()
          setLastRefreshedMoment(moment())
        },
        75000
        // 15000
      )
    },
    [lastRefreshedMoment, refreshData]
  )

  return (
    <section>
      <LoadingOverlay
        isFetching={!(auctionOfferSummaryIsFetching === false)}
        skeletonComponent={<TokenGrantSkeletonOverview />}
      >
        {depositState.name && depositState.name.includes("LIQUIDATION") ? 
        (<section
          key={"tokenGrant.id"}
          className="tile deposit-liquidation-overview"
          style={{ marginBottom: "1.2rem" }}
        >
          <div className="grant-amount">
            <ArbitrageurTokenDetails title="Your tBTC Balance" tokenValue={tBtcBalance} />
          </div>

          <div className="unlocking-details">
            <>
              <div className="flex-1 self-center">
                <CircularProgressBarEth
                  total={web3Context.web3.utils.toBN(100)
                    .mul(web3Context.web3.utils.toBN(10).pow(web3Context.web3.utils.toBN(18)))}
                  items={[
                    {
                      value: getPercentageOnOffer()
                        .mul(web3Context.web3.utils.toBN(10).pow(web3Context.web3.utils.toBN(18))), //TODO: Make sense. Looks like inner function for rendering legend need this as "wei"
                      color: colors.grey70,
                      backgroundStroke: colors.grey10,
                      label: "% On Offer",
                    },
                  ]}
                  withLegend
                />
              </div>
              <div className="ml-2 mt-1 self-start flex-1">
                <h5 className="text-grey-70">On Offer</h5>
                <h4 className="text-grey-70">{displayAmount(auctionValueBN, false)}</h4>
                <div className="text-smaller text-grey-40">
                  of {displayAmount(bondAmountWei, false)} Total
              </div>
              </div>
            </>

          </div>
          <div className="staked-details pl-0">
            <>
              <div className="flex-1 self-center">
                <h4>Liquidate this deposit</h4>
                <span>Purchase {displayAmount(auctionValueBN, false)} ETH in exchange of {satsToTBtcViaWeitoshi(depositSizeSatoshis).toString()} tBTC</span>

                <div>
                  <SubmitButton
                    className="btn btn-primary btn-sm"
                    onSubmitAction={(transactionHashCallback) =>
                      onLiquidateFromSummaryBtn(transactionHashCallback)
                    }
                  >
                    Purchase ETH
                </SubmitButton>
                </div>
              </div>
            </>
          </div>
        </section>)
         : (
          <section
          key={"tokenGrant.id"}
          className="tile deposit-liquidation-overview"
          style={{ marginBottom: "1.2rem" }}
        >
          <div className="grant-amount ">
            <ArbitrageurTokenDetails title="Your tBTC Balance" tokenValue={tBtcBalance} />
          </div>

          <div className="unlocking-details">
              <h4 className="flex-1 self-center">
                This deposit is not in liquidation
              </h4>
          </div>
        </section>
         ) }
        
      </LoadingOverlay>
    </section>
  )
}

// export default DepositLiquidationOverview