import React, { useEffect, useMemo, useCallback } from "react"
import TBTCEarningsDataTable from "../../components/TBTCEarningsDataTable"
import { LoadingOverlay } from "../../components/Loadable"
import {
  DataTableSkeleton,
  TokenAmountSkeleton,
} from "../../components/skeletons"
import TokenAmount from "../../components/TokenAmount"
import { tbtcRewardsService } from "../../services/tbtc-rewards.service"
import { useFetchData } from "../../hooks/useFetchData"
import { add } from "../../utils/arithmetics.utils"
import { TBTC } from "../../utils/token.utils"
import { findIndexAndObject } from "../../utils/array.utils"
import EmptyStatePage from "./EmptyStatePage"
import { useWeb3Address } from "../../components/WithWeb3Context"

const TBTCRewardsPage = () => {
  const address = useWeb3Address()
  const [
    { data, isFetching },
    updateRewardsData,
    ,
    setServiceArgs,
  ] = useFetchData(tbtcRewardsService.fetchTBTCRewards, [], address)

  useEffect(() => {
    if (address) {
      setServiceArgs([address])
    }
  }, [setServiceArgs, address])

  const totalAmount = useMemo(() => {
    return data
      .map((reward) => reward.amount)
      .reduce(add, 0)
      .toString()
  }, [data])

  const fetchOperatorByDepositId = useCallback(
    async (depositId) => {
      const operators = await tbtcRewardsService.fetchBeneficiaryOperatorsFromDeposit(
        address,
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
    [data, updateRewardsData, address]
  )

  return (
    <section className="tile">
      <h2 className="text-grey-70 mb-1">Total Amount</h2>
      {isFetching ? (
        <TokenAmountSkeleton textStyles={{ width: "35%" }} />
      ) : (
        <TokenAmount
          token={TBTC}
          amount={totalAmount}
          withIcon
          iconProps={{
            className: "tbtc-icon--mint-80",
          }}
        />
      )}

      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DataTableSkeleton />}
      >
        <TBTCEarningsDataTable
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
