const getUserAddress = (state) => state.app.address

const getCoveragePool = (state) => state.coveragePool

const selectors = {
  getUserAddress,
  getCoveragePool,
}

export default selectors
