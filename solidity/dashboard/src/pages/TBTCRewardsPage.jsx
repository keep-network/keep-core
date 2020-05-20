import React from "react"
import PageWrapper from "../components/PageWrapper"
import TBTCRewardsDataTable from "../components/TBTCRewardsDataTable"
import { tbtcRewardsService } from "../services/tbtc-rewards.service"
import { useWeb3Context } from "../components/WithWeb3Context"
import { useFetchData } from "../hooks/useFetchData"

const TBTCRewardsPage = () => {
  const { yourAddress } = useWeb3Context()
  const [{ data }] = useFetchData(
    tbtcRewardsService.fetchTBTCReawrds,
    [],
    yourAddress
  )

  return (
    <PageWrapper title="tBTC Rewards">
      <section className="tile">
        <h2 className="text-grey-70">Total Amount</h2>
        <h1 className="text-primary">
          4,000&nbsp;<span className="h3">TBTC</span>
        </h1>
        <TBTCRewardsDataTable rewards={data} />
      </section>
    </PageWrapper>
  )
}

export default React.memo(TBTCRewardsPage)
