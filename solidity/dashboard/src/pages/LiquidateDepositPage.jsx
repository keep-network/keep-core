import React, { useMemo, useCallback, useState, useEffect } from "react"
import { useParams } from "react-router-dom"
import FallingPriceAuctionTable from "../components/FallingPriceAuctionTable"
import StatusBadge, { BADGE_STATUS } from "../components/StatusBadge"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import {
  ArbitrageurTokenDetails, DepositAuctionOverview
} from "../components/DepositLiquidationOverview"
import { TokenGrantSkeletonOverview } from "../components/skeletons/TokenOverviewSkeleton"
import AddressShortcut from "../components/AddressShortcut"
import { isSameEthAddress } from "../utils/general.utils"
import { add } from "../utils/arithmetics.utils"
import moment from "moment"
import { LoadingOverlay } from "../components/Loadable"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"
import PageWrapper from "../components/PageWrapper"
import { ViewAddressInBlockExplorer } from "../components/ViewInBlockExplorer"
import { SubmitButton } from "../components/Button"
import { useShowMessage, messageType } from "../components/Message"
import { useWeb3Context } from "../components/WithWeb3Context"
import { liquidationService } from "../services/tbtc-liquidation.service"
import { useModal } from "../hooks/useModal"
import { useFetchData } from "../hooks/useFetchData"
import {
  satsToTBtcViaWeitoshi,
  displayAmount
} from "../utils/token.utils"
import { colors } from "../constants/colors"

const filterByOwned = (delegation) => !delegation.grantId
const filterBySelectedGrant = (selectedGrant) => (delegation) =>
  selectedGrant.id && delegation.grantId === selectedGrant.id

const LiquidateDepositPage = (props) => {
  const {
    undelegations,
    delegations,
    tokensContext,
    selectedGrant,
    isFetching,
    grantsAreFetching,
    keepTokenBalance,
    availableTopUps,
    grants,
  } = useTokensPageContext()


  // const { depositAddress } = props
  const web3Context = useWeb3Context()
  const { depositAddress } = useParams()
  const { openConfirmationModal } = useModal()

  const [lastRefreshedMoment, setLastRefreshedMoment] = useState(moment())

  useEffect(
    () => {
      setTimeout(
        () => {
          refreshData()
          setLastRefreshedMoment(moment())
          console.log(`I refreshed at ${moment().toString()}`)
        },
        75000
        // 15000
      )
    },
    [lastRefreshedMoment]
  )
  // startRefreshDataTimer()

  const refreshData = () => {
    refreshCurrentAuctionValue()
    refreshUserTBtcBalance()
    refreshDepositState()
  }

  const [stCurrentAuctionValue, , refreshCurrentAuctionValue] = useFetchData(
    liquidationService.getDepositCurrentAuctionValue,
    {},
    depositAddress
  )
  const {
    isFetching: currentAuctionValueIsFetching,
    data: auctionValueBN,
  } = stCurrentAuctionValue

  const [stUserTBtcBalance, , refreshUserTBtcBalance] = useFetchData(
    liquidationService.getTBtcBalanceOf,
    {})
  const {
    isFetching: userTBtcBalanceIsFetching,
    data: tBtcBalance,
  } = stUserTBtcBalance

  const [stDepositState, , refreshDepositState] = useFetchData(
    liquidationService.getDepositState,
    {},
    depositAddress
    )
  const {
    isFetching: depositStateIsFetching,
    data: depositStateObj,
  } = stDepositState


  const [stAuctionSchedule] = useFetchData(
    liquidationService.getDepositAuctionOfferingSchedule,
    {},
    depositAddress)
  const {
    isFetching: auctionScheduleIsFetching,
    data: auctionOfferingSchedule,
  } = stAuctionSchedule

  const [stLastStartedLiquidationEvent] = useFetchData(
    liquidationService.getLastStartedLiquidationEvent,
    {},
    depositAddress
  )
  const {
    isFetching: lastStartedLiquidationEventIsFetching,
    data: startedLiquidationEvent,
  } = stLastStartedLiquidationEvent

  const [stDepositBondAmount] = useFetchData(
    liquidationService.getDepositEthBalance,
    {},
    depositAddress
  )
  const {
    isFetching: bondAmountIsFetching,
    data: bondAmountWei,
  } = stDepositBondAmount

  const [stDepositSizeSatoshis] = useFetchData(
    liquidationService.getDepositSizeSatoshis,
    {},
    depositAddress
  )
  const {
    isFetching: depositSizeSatoshisIsFetching,
    data: depositSizeSatoshis,
  } = stDepositSizeSatoshis

  const confirmationModalOptions = useCallback(() => {
    if (bondAmountIsFetching || depositSizeSatoshisIsFetching)
      return {}
    else {
      return {
        modalOptions: { title: "Purchase ETH Bond" },
        title: "You’re about to purchase ETH with tBTC.",
        subtitle:
          `This transaction will spend ${satsToTBtcViaWeitoshi(depositSizeSatoshis).toString()} tBTC to obtain ${displayAmount(bondAmountWei, false)} ETH (or more, depending on the block it goes through).
           It can fail if deposit state changes before this transaction gets accepted. 
           Transaction can also fail if you don’t have enough tBTC`,
        btnText: "Purchase ETH",
        confirmationText: "Y",
      }
    }
  }, [satsToTBtcViaWeitoshi, displayAmount, bondAmountIsFetching, depositSizeSatoshisIsFetching]) 

  const getPercentageOnOffer = useCallback(() => {
    const utils = web3Context.web3.utils
    let pct = utils.toBN(0)
    if (!currentAuctionValueIsFetching && !bondAmountIsFetching) {
      pct = auctionValueBN.mul(utils.toBN(100))
        .div(utils.toBN(bondAmountWei))
    }
    // console.log(`getPercentageOnOffer: ${pct}`)
    return pct
  }, [web3Context, bondAmountIsFetching, bondAmountWei, currentAuctionValueIsFetching, auctionValueBN])

  const showMessage = useShowMessage()

  const onLiquidateFromSummaryBtn = async () => {
    try {
      // const availableAmount = delegationData.isFromGrant
      //   ? getAvailableToStakeFromGrant(delegationData.grantId)
      //   : keepTokenBalance
      // const { amount } = 19
      // delegationData.beneficiaryAddress = delegationData.beneficiary
      // delegationData.stakeTokens = amount
      // delegationData.selectedGrant = {
      //   id: delegationData.grantId,
      //   isManagedGrant: delegationData.isManagedGrant,
      //   managedGrantContractInstance:
      //     delegationData.managedGrantContractInstance,
      // }
      // delegationData.context = delegationData.isFromGrant ? "granted" : "owned"
      await liquidationService.depositNotifySignatureTimeout(
        web3Context,
        depositAddress
      )
      showMessage({
        type: messageType.SUCCESS,
        title: "Success",
        content: "Top up committed successfully",
      })
    } catch (error) {
      showMessage({
        type: messageType.ERROR,
        title: "Commit action has failed ",
        content: error.message,
      })
      throw error
    }
  }

    const handleSubmit = async (onTransactionHashCallback) => {
    try {
      await openConfirmationModal(confirmationModalOptions())
      // const depositSizeWeitoshi = fromTokenUnit(depositSizeSatoshis, 10)
      await liquidationService.purchaseDepositAtAuction(
        web3Context,
        depositAddress,
        onTransactionHashCallback
      )
      showMessage({
        type: messageType.SUCCESS,
        title: "Success",
        content: "Staking delegate transaction has been successfully completed",
      })
    } catch (error) {
      console.error(error)
      showMessage({
        type: messageType.ERROR,
        title: "Staking delegate action has failed ",
        content: error.message,
      })
      throw error
    }
  }


  //TODO: Fix first section. Conditional should be at a higher level and check if it is a deposit.
  // If it's not in liquidation (has ever been) it will fail startedLiquidationEvent.returnValues
  return (
    // <PageWrapper title="Liquidations">
    <section>
      <div className="flex wrap self-center mb-2">
        <h2 className="text-grey-70">
          {satsToTBtcViaWeitoshi(depositSizeSatoshis).toString()}{` tBTC Deposit`}
        </h2>
        
        {lastStartedLiquidationEventIsFetching === false && depositStateIsFetching === false && (
          <>
            <span className="flex self-center ml-2">
              <ViewAddressInBlockExplorer address={depositAddress} urlSuffix={""} />
            </span>
            <span className="flex self-center ml-2">
              <StatusBadge
                className="self-center"
                status={BADGE_STATUS.DISABLED}
                text={depositStateObj.name}
              />
              {/* <span className="h4 text-grey-50 ml-1">
                {!lastStartedLiquidationEventIsFetching &&
                  moment.unix(startedLiquidationEvent.returnValues._timestamp).toString()}
              </span> */}
            </span>
          </>
        )}
      </div>
      <>
        <DepositAuctionOverview
          auctionOfferSummaryIsFetching={currentAuctionValueIsFetching || bondAmountIsFetching || userTBtcBalanceIsFetching || depositSizeSatoshisIsFetching || depositStateIsFetching}
          depositState={depositStateObj}
          tBtcBalance={tBtcBalance}
          auctionValueBN={auctionValueBN}
          bondAmountWei={bondAmountWei}
          depositSizeSatoshis={depositSizeSatoshis}
          refreshData={refreshData}
          getPercentageOnOffer={getPercentageOnOffer}
          onLiquidateFromSummaryBtn={handleSubmit}
        >
        </DepositAuctionOverview>
      </>
      {/* <LoadingOverlay
        isFetching={grantsAreFetching || currentAuctionValueIsFetching || bondAmountIsFetching || userTBtcBalanceIsFetching || depositSizeSatoshisIsFetching || lastStartedLiquidationEventIsFetching}
        skeletonComponent={<TokenGrantSkeletonOverview />}
      >
        <section
          key={"tokenGrant.id"}
          className="tile deposit-liquidation-overview"
          style={{ marginBottom: "1.2rem" }}
        >
          <div className="grant-amount">
            <ArbitrageurTokenDetails title="Your tBTC Balance" selectedGrant={tokenGrant} tokenValue={tBtcBalance} />
          </div>
          <div className="unlocking-details">
            <>
              <div className="flex-1 self-center">
                <CircularProgressBars
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
                      onLiquidateFromSummaryBtn(null, transactionHashCallback)
                    }
                  >
                    Purchase ETH
                </SubmitButton>
                </div>
              </div>
            </>
          </div>
        </section>
      </LoadingOverlay> */}

      <LoadingOverlay
        isFetching={
          (tokensContext === "granted" ? grantsAreFetching : isFetching) || auctionScheduleIsFetching
        }
        skeletonComponent={<DataTableSkeleton />}
      >
        <FallingPriceAuctionTable
          auctionScheduleData={auctionOfferingSchedule}
          // cancelStakeSuccessCallback={cancelStakeSuccessCallback}
        />
        {/* <DelegatedTokensTable
          delegatedTokens={getDelegations()}
          cancelStakeSuccessCallback={cancelStakeSuccessCallback}
          keepTokenBalance={keepTokenBalance}
          grants={grants}
        /> */}
      </LoadingOverlay>
    </section>
    // </PageWrapper>
  )
}


export default LiquidateDepositPage
