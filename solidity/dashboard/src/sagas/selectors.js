const getUserAddress = (state) => state.app.address

const getCoveragePool = (state) => state.coveragePool

const getModalData = (state) => state.modal

const selectors = {
  getUserAddress,
  getCoveragePool,
  getModalData,
}

export default selectors
