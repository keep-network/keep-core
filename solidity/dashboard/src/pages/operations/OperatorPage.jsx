import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import DelegatedTokens from "../../components/DelegatedTokens"
import PendingUndelegation from "../../components/PendingUndelegation"
import Tile from "../../components/Tile"
import SlashedTokensList from "../../components/SlashedTokensList"
import PageWrapper from "../../components/PageWrapper"
import { LoadingOverlay } from "../../components/Loadable"
import DelegatedTokensSkeleton from "../../components/skeletons/DelegatedTokensSkeleton"
import { useWeb3Address } from "../../components/WithWeb3Context"
import { DataTableSkeleton } from "../../components/skeletons"
import {
  FETCH_OPERATOR_DELEGATIONS_RERQUEST,
  FETCH_OPERATOR_SLASHED_TOKENS_RERQUEST,
  OPERATR_DELEGATION_CANCELED,
} from "../../actions"

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

  return (
    <PageWrapper title={title}>
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DelegatedTokensSkeleton />}
      >
        <DelegatedTokens
          data={data}
          cancelSuccessCallback={() => {
            dispatch({ type: OPERATR_DELEGATION_CANCELED })
          }}
        />
      </LoadingOverlay>
      <PendingUndelegation data={data} />
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
