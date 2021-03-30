import React from "react"
import DelegatedTokens from "../../components/DelegatedTokens"
import PendingUndelegation from "../../components/PendingUndelegation"
import SlashedTokens from "../../components/SlashedTokens"
import { useSubscribeToContractEvent } from "../../hooks/useSubscribeToContractEvent"
import { TOKEN_STAKING_CONTRACT_NAME } from "../../constants/constants"
import PageWrapper from "../../components/PageWrapper"
import { LoadingOverlay } from "../../components/Loadable"
import DelegatedTokensSkeleton from "../../components/skeletons/DelegatedTokensSkeleton"
import { useDispatch, useSelector } from "react-redux"
import { useEffect } from "react"
import { FETCH_OPERATOR_DELEGATIONS_RERQUEST } from "../../actions"
import { useWeb3Address } from "../../components/WithWeb3Context"

const OperatorPage = ({ title }) => {
  const dispatch = useDispatch()
  const address = useWeb3Address()

  useEffect(() => {
    dispatch({
      type: FETCH_OPERATOR_DELEGATIONS_RERQUEST,
      payload: { address },
    })
  }, [dispatch, address])

  const { isFetching, ...data } = useSelector((state) => state.operator)

  const { latestEvent } = useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "Undelegated"
  )

  return (
    <PageWrapper title={title}>
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DelegatedTokensSkeleton />}
      >
        <DelegatedTokens
          isFetching={isFetching}
          data={data}
          // setData={setData}
        />
      </LoadingOverlay>
      <PendingUndelegation
        latestUnstakeEvent={latestEvent}
        data={data}
        // setData={setData}
      />
      <SlashedTokens />
    </PageWrapper>
  )
}

export default OperatorPage
