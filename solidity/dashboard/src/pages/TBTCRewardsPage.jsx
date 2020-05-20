import React, { useMemo } from "react"
import PageWrapper from "../components/PageWrapper"
import TBTCRewardsDataTable from "../components/TBTCRewardsDataTable"
import { tbtcRewardsService } from "../services/tbtc-rewards.service"
import { useWeb3Context } from "../components/WithWeb3Context"
import { useFetchData } from "../hooks/useFetchData"
import { add } from "../utils/arithmetics.utils"
import { displayAmount } from "../utils/token.utils"

const TBTCRewardsPage = () => {
  const { yourAddress } = useWeb3Context()
  const [{ data }] = useFetchData(
    tbtcRewardsService.fetchTBTCReawrds,
    [],
    yourAddress
  )

  const totalAmount = useMemo(() => {
    return displayAmount(data.map((reward) => reward.amount).reduce(add, 0))
  }, [data])

  return (
    <PageWrapper title="tBTC Rewards">
      <section className="tile">
        <h2 className="text-grey-70">Total Amount</h2>
        <h1 className="text-primary">
          {totalAmount}&nbsp;<span className="h3">TBTC</span>
        </h1>
        <TBTCRewardsDataTable rewards={data} />
      </section>
    </PageWrapper>
  )
}

export default React.memo(TBTCRewardsPage)
