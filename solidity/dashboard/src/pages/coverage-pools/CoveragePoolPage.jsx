import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import {
  DepositForm,
  InitiateDepositModal,
} from "../../components/coverage-pools"
import TokenAmount from "../../components/TokenAmount"
import MetricsTile from "../../components/MetricsTile"
import { APY } from "../../components/liquidity"
import { Skeleton } from "../../components/skeletons"
import { useWeb3Address } from "../../components/WithWeb3Context"
import OnlyIf from "../../components/OnlyIf"
import {
  fetchTvlRequest,
  fetchCovPoolDataRequest,
  depositAssetPool,
  fetchAPYRequest,
} from "../../actions/coverage-pool"
import { useModal } from "../../hooks/useModal"
import { lte } from "../../utils/arithmetics.utils"
import { KEEP } from "../../utils/token.utils"
import { displayPercentageValue } from "../../utils/general.utils"

const CoveragePoolPage = ({ title, withNewLabel }) => {
  const { openConfirmationModal } = useModal()
  const dispatch = useDispatch()
  const {
    totalValueLocked,
    totalValueLockedInUSD,
    isTotalValueLockedFetching,
    shareOfPool,
    estimatedRewards,
    estimatedKeepBalance,
    apy,
    isApyFetching,
    totalAllocatedRewards,
  } = useSelector((state) => state.coveragePool)
  const keepTokenBalance = useSelector((state) => state.keepTokenBalance)

  const address = useWeb3Address()

  useEffect(() => {
    dispatch(fetchTvlRequest())
    dispatch(fetchAPYRequest())
  }, [dispatch])

  useEffect(() => {
    if (address) {
      dispatch(fetchCovPoolDataRequest(address))
    }
  }, [dispatch, address])

  const onSubmitDepositForm = async (values, awaitingPromise) => {
    const { tokenAmount } = values
    const amount = KEEP.fromTokenUnit(tokenAmount)
    await openConfirmationModal(
      {
        modalOptions: { title: "Initiate Deposit" },
        submitBtnText: "deposit",
        amount,
      },
      InitiateDepositModal
    )
    dispatch(depositAssetPool(amount, awaitingPromise))
  }

  return (
    <>
      <section className="tile coverage-pool__overview">
        <section className="coverage-pool__overview__tvl">
          <h2 className="h2--alt text-grey-70 mb-1">Total Value Locked</h2>
          <TokenAmount
            amount={totalValueLocked}
            amountClassName="h1 text-mint-100"
            symbolClassName="h2 text-mint-100"
            withIcon
          />
          <h3 className="tvl tvl--usd">
            {`$${totalValueLockedInUSD.toString()} USD`}
          </h3>
        </section>
        <div className="coverage-pool__overview__metrics">
          <section className="metrics__apy">
            <h4 className="text-grey-70 mb-1">Rewards Rate</h4>

            <MetricsTile className="bg-mint-10 mr-2">
              <APY
                apy={apy}
                isFetching={isApyFetching}
                className="text-mint-100"
              />
              <h5 className="text-grey-60">annual</h5>
            </MetricsTile>
          </section>
          <section className="metrics__total-rewards">
            <h4 className="text-grey-70 mb-1">Total Rewards</h4>

            <MetricsTile className="bg-mint-10">
              {isTotalValueLockedFetching ? (
                <Skeleton tag="h2" shining color="grey-10" />
              ) : (
                <TokenAmount
                  amount={totalAllocatedRewards}
                  withIcon
                  withSymbol={false}
                  withMetricSuffix
                />
              )}
              <h5 className="text-grey-60">pool lifetime</h5>
            </MetricsTile>
          </section>
        </div>

        {/* TODO add more metrics according to the Figma vies */}
      </section>
      <section className="coverage-pool__deposit-wrapper">
        <section className="tile coverage-pool__deposit-form">
          <h3>Deposit</h3>
          <DepositForm
            onSubmit={onSubmitDepositForm}
            tokenAmount={keepTokenBalance.value}
            apy={apy}
          />
        </section>

        <section className="tile coverage-pool__share-of-pool">
          <h4 className="text-grey-70 mb-3">Your Share of Pool</h4>

          <OnlyIf condition={shareOfPool <= 0}>
            <div className="text-grey-30 text-center">
              You have no balance yet.&nbsp;
              <br />
              <u>Deposit KEEP</u>&nbsp;to see balance.
            </div>
          </OnlyIf>
          <OnlyIf condition={shareOfPool > 0}>
            <div className="flex column center">
              <TokenAmount amount={estimatedKeepBalance} withSymbol={false} />
              <h4 className="text-mint-100">{KEEP.symbol}</h4>
              <div className="text-grey-40 mt-2">
                <b>{displayPercentageValue(shareOfPool * 100, false)}</b>
                &nbsp;of Pool
              </div>
            </div>
          </OnlyIf>
        </section>

        <section className="tile coverage-pool__rewards">
          <h4 className="text-grey-70 mb-3">Your Rewards</h4>
          <OnlyIf condition={lte(estimatedRewards, 0) && shareOfPool <= 0}>
            <div className="text-grey-30 text-center">
              You have no rewards yet.&nbsp;
              <br />
              <u>Deposit KEEP</u>&nbsp;to see rewards.
            </div>
          </OnlyIf>
          <OnlyIf condition={shareOfPool > 0}>
            <div className="flex column center">
              <TokenAmount amount={estimatedRewards} withSymbol={false} />
              <h4 className="text-mint-100">{KEEP.symbol}</h4>
            </div>
          </OnlyIf>
        </section>
      </section>
    </>
  )
}

CoveragePoolPage.route = {
  title: "Deposit",
  path: "/coverage-pools/deposit",
  exact: true,
}

export default CoveragePoolPage
