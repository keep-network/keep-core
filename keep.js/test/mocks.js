import sinon from "sinon"

export const ContractWrapperMock = {
  makeCall: sinon.spy(),
  sendTransaction: sinon.spy(),
  getPastEvents: sinon.spy(),
}
