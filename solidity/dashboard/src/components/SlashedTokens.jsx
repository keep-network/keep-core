import React from "react"
import SlashedTokensList from "./SlashedTokensList"
import { LoadingOverlay } from "./Loadable"
import { useFetchData } from "../hooks/useFetchData"
import { slashedTokensService } from "../services/slashed-tokens.service"
import Tile from "./Tile"

const SlashedTokens = (props) => {
  const [state] = useFetchData(slashedTokensService.fetchSlashedTokens, [])
  const { data, isFetching } = state

  return (
    <LoadingOverlay isFetching={isFetching}>
      <Tile id="slashed-tokens">
        <SlashedTokensList slashedTokens={data} />
      </Tile>
    </LoadingOverlay>
  )
}

export default SlashedTokens
