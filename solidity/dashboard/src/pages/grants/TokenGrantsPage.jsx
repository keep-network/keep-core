import React, { useEffect, useMemo } from "react"
import Web3 from "web3"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket"
import ProviderEngine from "web3-provider-engine"
import { useSelector, useDispatch } from "react-redux"
import { useParams } from "react-router-dom"
import { Link } from "react-router-dom"
import PageWrapper from "../../components/PageWrapper"
import { TokenGrantDetails } from "../../components/TokenGrantOverview"
import TokenAmount from "../../components/TokenAmount"
import { LoadingOverlay } from "../../components/Loadable"
import { TokenGrantSkeletonOverview } from "../../components/skeletons/TokenOverviewSkeleton"
import { CircularProgressBars } from "../../components/CircularProgressBar"
import { SubmitButton } from "../../components/Button"
import { add, gt } from "../../utils/arithmetics.utils"
import {
  displayAmount,
  displayAmountWithMetricSuffix,
} from "../../utils/token.utils"
import { colors } from "../../constants/colors"
import useReleaseTokens from "../../hooks/useReleaseTokens"
import { useFetchData } from "../../hooks/useFetchData"
import { tokenGrantsService } from "../../services/token-grants.service"
import {
  getContracts,
  resolveWeb3Deferred,
  resovleContractsDeferred,
} from "../../contracts"
import { getWsUrl } from "../../connectors/utils"
import { Web3Context } from "../../components/WithWeb3Context"

const TokenGrantsPage = (props) => {
  const dispatch = useDispatch()
  const { grants, isFetching } = useSelector((state) => state.tokenGrants)

  useEffect(() => {
    dispatch({ type: "token-grant/fetch_grants_request" })
  }, [dispatch])

  const totalGrantAmount = useMemo(() => {
    return grants.map(({ amount }) => amount).reduce(add, "0")
  }, [grants])

  return (
    <PageWrapper {...props}>
      <TokenAmount
        wrapperClassName="mb-2"
        amount={totalGrantAmount}
        amountClassName="h2 text-grey-40"
        currencyIconProps={{ className: "keep-outline grey-40" }}
        displayWithMetricSuffix={false}
      />

      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<TokenGrantSkeletonOverview />}
      >
        {grants.map(renderTokenGrantOverview)}
      </LoadingOverlay>
    </PageWrapper>
  )
}

const renderTokenGrantOverview = (tokenGrant) => (
  <TokenGrantOverview key={tokenGrant.id} tokenGrant={tokenGrant} />
)

const TokenGrantOverview = React.memo(({ tokenGrant }) => {
  return (
    <section
      key={tokenGrant.id}
      className="tile token-grant-overview"
      style={{ marginBottom: "1.2rem" }}
    >
      <div className="grant-amount">
        <header className="flex row center space-between mb-1">
          <h3 className="text-grey-70">Grant Amount</h3>
          <Link
            to={{
              pathname: "/delegation/grant",
              hash: `${tokenGrant.id}`,
            }}
            className="btn btn-secondary btn-sm"
          >
            delegate
          </Link>
        </header>
        <TokenGrantDetails
          selectedGrant={tokenGrant}
          availableAmount={tokenGrant.availableToStake}
        />
      </div>
      <div className="unlocking-details">
        <TokenGrantUnlockingdDetails selectedGrant={tokenGrant} />
      </div>
      <div className="staked-details">
        <TokenGrantStakedDetails
          selectedGrant={tokenGrant}
          stakedAmount={tokenGrant.stakedAmount}
        />
      </div>
    </section>
  )
})

export const TokenGrantStakedDetails = ({ selectedGrant, stakedAmount }) => {
  return (
    <>
      <div className="flex-1 self-center">
        <CircularProgressBars
          total={selectedGrant.amount}
          items={[
            {
              value: stakedAmount,
              color: colors.grey70,
              backgroundStroke: colors.grey10,
              label: "Staked",
            },
          ]}
          withLegend
        />
      </div>
      <div className="ml-2 mt-1 self-start flex-1">
        <h5 className="text-grey-70">staked</h5>
        <h4 className="text-grey-70">{displayAmount(stakedAmount)}</h4>
        <div className="text-smaller text-grey-40">
          of {displayAmountWithMetricSuffix(selectedGrant.amount)} Total
        </div>
      </div>
    </>
  )
}

const TokenGrantUnlockingdDetails = ({
  selectedGrant,
  hideReleaseTokensBtn = false,
}) => {
  const releaseTokens = useReleaseTokens()

  const onReleaseTokens = async (awaitingPromise) => {
    releaseTokens(selectedGrant, awaitingPromise)
  }

  return (
    <>
      <div className="flex-1">
        <CircularProgressBars
          total={selectedGrant.amount}
          items={[
            {
              value: selectedGrant.unlocked,
              backgroundStroke: colors.bgSuccess,
              color: colors.primary,
              label: "Unlocked",
            },
            {
              value: selectedGrant.released,
              color: colors.secondary,
              backgroundStroke: colors.bgSecondary,
              radius: 48,
              label: "Released",
            },
          ]}
          withLegend
        />
      </div>
      <div
        className={`ml-2 mt-1 flex-1${
          selectedGrant.readyToRelease === "0" ? " self-start" : ""
        }`}
      >
        <h5 className="text-grey-70">unlocked</h5>
        <h4 className="text-grey-70">
          {displayAmount(selectedGrant.unlocked)}
        </h4>
        <div className="text-smaller text-grey-40">
          of {displayAmountWithMetricSuffix(selectedGrant.amount)} Total
        </div>
        {gt(selectedGrant.readyToRelease || 0, 0) && (
          <div className="mt-2">
            <div className="text-secondary text-small flex wrap">
              <span className="mr-1">
                {`${displayAmountWithMetricSuffix(
                  selectedGrant.readyToRelease
                )} Available`}
              </span>
            </div>
            {!hideReleaseTokensBtn && (
              <SubmitButton
                className="btn btn-sm btn-secondary"
                onSubmitAction={onReleaseTokens}
                withMessageActionIsPending={false}
              >
                release tokens
              </SubmitButton>
            )}
          </div>
        )}
      </div>
    </>
  )
}

const TokenGrantPreview = (props) => {
  const { grantId } = useParams()
  const [state] = useFetchData(tokenGrantsService.fetchGrantById, {}, grantId)

  return (
    <PageWrapper {...props}>
      <h2 className="h2--alt mb-2 text-grey-60">Grant ID {grantId}</h2>
      <LoadingOverlay
        isFetching={state.isFetching}
        skeletonComponent={<TokenGrantSkeletonOverview />}
      >
        {state.isError ? (
          <section className="tile flex full-center">
            <h3 className="text-validation">{state.error.message}</h3>
          </section>
        ) : (
          <section className="tile token-grant-overview">
            <div className="grant-amount">
              <h3 className="text-grey-70">Grant Amount</h3>
              <TokenGrantDetails
                selectedGrant={state.data}
                availableAmount={state.data.availableToStake}
              />
            </div>
            <div className="unlocking-details">
              <TokenGrantUnlockingdDetails
                selectedGrant={state.data}
                hideReleaseTokensBtn
              />
            </div>
            <div className="staked-details">
              <TokenGrantStakedDetails
                selectedGrant={state.data}
                stakedAmount={state.data.staked}
              />
            </div>
          </section>
        )}
      </LoadingOverlay>
    </PageWrapper>
  )
}

const getCustomWeb3Context = () => {
  const engine = new ProviderEngine()
  const web3 = new Web3(engine)
  engine.addProvider(
    new WebsocketSubprovider({ rpcUrl: getWsUrl(), debug: false })
  )
  engine.start()
  const contracts = getContracts(web3)
  resolveWeb3Deferred(web3)
  resovleContractsDeferred(contracts)

  return {
    web3,
    ...contracts,
  }
}

const TokenGrantPreviewPage = (props) => {
  return (
    <Web3Context.Provider value={getCustomWeb3Context()}>
      <TokenGrantPreview {...props} />
    </Web3Context.Provider>
  )
}

export { TokenGrantPreviewPage }

export default TokenGrantsPage
