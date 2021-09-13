const getUserAddress = (state) => state.app.address

const getCoveragePool = (state) => state.coveragePool

const getTBTCV2Migration = (state) => state.tbtcV2Migration
const getModalData = (state) => state.modal

const selectors = {
  getUserAddress,
  getCoveragePool,
  getTBTCV2Migration,
  getModalData,
}

export default selectors
