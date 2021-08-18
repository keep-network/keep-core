const getUserAddress = (state) => state.app.address

const getCoveragePool = (state) => state.coveragePool

const getTBTCV2Migration = (state) => state.tbtcV2Migration

const selectors = {
  getUserAddress,
  getCoveragePool,
  getTBTCV2Migration,
}

export default selectors
