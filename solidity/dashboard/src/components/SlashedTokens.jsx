import React from "react"
import SlashedTokensList from "./SlashedTokensList"
import { LoadingOverlay } from "./Loadable"
import { useFetchData } from "../hooks/useFetchData"
import { slashedTokensService } from "../services/slashed-tokens.service"
import Tile from "./Tile"
import DataTableSkeleton from "./skeletons/DataTableSkeleton"

const SlashedTokens = (props) => {
  const [state] = useFetchData(slashedTokensService.fetchSlashedTokens, [])
  const { data, isFetching } = state

  return (
    <LoadingOverlay
      isFetching={isFetching}
      skeletonComponent={<DataTableSkeleton columns={2} />}
    >
      <Tile id="slashed-tokens">
        <SlashedTokensList slashedTokens={data} />
      </Tile>
    </LoadingOverlay>
  )
}

export default SlashedTokens
