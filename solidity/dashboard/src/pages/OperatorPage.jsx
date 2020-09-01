import React from "react"
import DelegatedTokens from "../components/DelegatedTokens"
import PendingUndelegation from "../components/PendingUndelegation"
import SlashedTokens from "../components/SlashedTokens"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent"
import { TOKEN_STAKING_CONTRACT_NAME } from "../constants/constants"
import PageWrapper from "../components/PageWrapper"
import { operatorService } from "../services/token-staking.service"
import { useFetchData } from "../hooks/useFetchData"
import { LoadingOverlay } from "../components/Loadable"
import DelegatedTokensSkeleton from "../components/skeletons/DelegatedTokensSkeleton"

const OperatorPage = (props) => {
  const [state, setData] = useFetchData(
    operatorService.fetchDelegatedTokensData,
    {}
  )
  const { isFetching, data } = state

  const { latestEvent } = useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "Undelegated"
  )

  return (
    <PageWrapper title="Operations">
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DelegatedTokensSkeleton />}
      >
        <DelegatedTokens
          isFetching={isFetching}
          data={data}
          setData={setData}
        />
      </LoadingOverlay>
      <PendingUndelegation
        latestUnstakeEvent={latestEvent}
        data={data}
        setData={setData}
      />
      <SlashedTokens />
    </PageWrapper>
  )
}

export default OperatorPage
