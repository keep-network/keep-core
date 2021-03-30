import React from "react"
import DelegatedTokens from "../../components/DelegatedTokens"
import PendingUndelegation from "../../components/PendingUndelegation"
import Tile from "../../components/Tile"
import SlashedTokensList from "../../components/SlashedTokensList"
import { useSubscribeToContractEvent } from "../../hooks/useSubscribeToContractEvent"
import { TOKEN_STAKING_CONTRACT_NAME } from "../../constants/constants"
import PageWrapper from "../../components/PageWrapper"
import { LoadingOverlay } from "../../components/Loadable"
import DelegatedTokensSkeleton from "../../components/skeletons/DelegatedTokensSkeleton"
import { useDispatch, useSelector } from "react-redux"
import { useEffect } from "react"
import {
  FETCH_OPERATOR_DELEGATIONS_RERQUEST,
  FETCH_OPERATOR_SLASHED_TOKENS_RERQUEST,
} from "../../actions"
import { useWeb3Address } from "../../components/WithWeb3Context"
import { DataTableSkeleton } from "../../components/skeletons"

const OperatorPage = ({ title }) => {
  const dispatch = useDispatch()
  const address = useWeb3Address()

  useEffect(() => {
    dispatch({
      type: FETCH_OPERATOR_DELEGATIONS_RERQUEST,
      payload: { address },
    })
    dispatch({
      type: FETCH_OPERATOR_SLASHED_TOKENS_RERQUEST,
      payload: { address },
    })
  }, [dispatch, address])

  const {
    isFetching,
    areSlashedTokensFetching,
    slashedTokens,
    ...data
  } = useSelector((state) => state.operator)

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
      <LoadingOverlay
        isFetching={areSlashedTokensFetching}
        skeletonComponent={<DataTableSkeleton columns={2} />}
      >
        <Tile id="slashed-tokens">
          <SlashedTokensList slashedTokens={slashedTokens} />
        </Tile>
      </LoadingOverlay>
    </PageWrapper>
  )
}

export default OperatorPage
