import { put } from "redux-saga/effects"
import { tbtcV2Migration } from "../actions"
import { takeOnlyOnce } from "./effects"
import { identifyTaskByAddress, logErrorAndThrow } from "./utils"

function* fetchData() {
  try {
    yield put(tbtcV2Migration.fetchDataStart())
    // TODO: Fetch data
    const data = {}
    yield put(tbtcV2Migration.fetchDataSuccess(data))
  } catch (error) {
    yield* logErrorAndThrow(
      tbtcV2Migration.TBTCV2_MIGRATION_FETCH_DATA_ERROR,
      error
    )
  }
}

export function* watchFetchTvl() {
  yield takeOnlyOnce(
    tbtcV2Migration.TBTCV2_MIGRATION_FETCH_DATA_REQUEST,
    identifyTaskByAddress,
    fetchData
  )
}
