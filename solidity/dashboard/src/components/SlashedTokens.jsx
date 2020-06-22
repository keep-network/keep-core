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
      skeletonComponent={<DataTableSkeleton />}
    >
      <Tile
        title="Slashed Tokens"
        id="slashed-tokens"
        withTooltip
        tooltipProps={{
          text: (
            <>
              A &nbsp;<span className="text-bold">slash</span>&nbsp; is a
              penalty for signing group misbehavior. It results in a removal of
              a portion of your delegated KEEP tokens.
            </>
          ),
        }}
      >
        <SlashedTokensList slashedTokens={data} />
      </Tile>
    </LoadingOverlay>
  )
}

export default SlashedTokens
