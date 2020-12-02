import React, { useMemo, useCallback } from "react"
import TBTCRewardsDataTable from "../../components/TBTCRewardsDataTable"
import { LoadingOverlay } from "../../components/Loadable"
import { useWeb3Context } from "../../components/WithWeb3Context"
import {
  DataTableSkeleton,
  TokenAmountSkeleton,
} from "../../components/skeletons"
import TokenAmount from "../../components/TokenAmount"
import * as Icons from "../../components/Icons"
import { tbtcRewardsService } from "../../services/tbtc-rewards.service"
import { useFetchData } from "../../hooks/useFetchData"
import { add } from "../../utils/arithmetics.utils"
import { toTokenUnit } from "../../utils/token.utils"
import { findIndexAndObject } from "../../utils/array.utils"
import EmptyStatePage from "./EmptyStatePage"

const TBTCRewardsPage = () => {
  const web3Context = useWeb3Context()
  const { yourAddress } = web3Context
  const [{ data, isFetching }, updateRewardsData] = useFetchData(
    tbtcRewardsService.fetchTBTCRewards,
    [],
    yourAddress
  )

  const totalAmount = useMemo(() => {
    return data
      .map((reward) => reward.amount)
      .reduce(add, 0)
      .toString()
  }, [data])

  const fetchOperatorByDepositId = useCallback(
    async (depositId) => {
      const operators = await tbtcRewardsService.fetchBeneficiaryOperatorsFromDeposit(
        web3Context,
        yourAddress,
        depositId
      )
      const updatedData = [...data]
      operators.forEach((operator) => {
        const { indexInArray, obj } = findIndexAndObject(
          "depositTokenId",
          depositId,
          updatedData
        )

        if (indexInArray !== null) {
          updatedData[indexInArray] = { ...obj, operatorAddress: operator }
        }
      })

      updateRewardsData(updatedData)
    },
    [data, updateRewardsData, yourAddress, web3Context]
  )

  return (
    <section className="tile">
      <h2 className="text-grey-70">Total Amount</h2>
      {isFetching ? (
        <TokenAmountSkeleton textStyles={{ width: "35%" }} />
      ) : (
        <TokenAmount
          currencyIcon={Icons.TBTC}
          currencyIconProps={{
            className: "tbtc-icon--mint-80",
            width: 32,
            height: 32,
          }}
          amount={totalAmount}
          currencySymbol="tBTC"
          displayWithMetricSuffix={false}
          displayAmountFunction={(amount) => toTokenUnit(amount).toString()}
        />
      )}

      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DataTableSkeleton />}
      >
        <TBTCRewardsDataTable
          rewards={data}
          fetchOperatorByDepositId={fetchOperatorByDepositId}
        />
      </LoadingOverlay>
    </section>
  )
}

TBTCRewardsPage.route = {
  title: "tBTC",
  path: "/earnings/tbtc",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

export default TBTCRewardsPage
