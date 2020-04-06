import web3Utils from 'web3-utils'

export const add = (a, b) => {
  return web3Utils.toBN(a).add(web3Utils.toBN(b))
}

export const sub = (a, b) => {
  return web3Utils.toBN(a).sub(web3Utils.toBN(b))
}

export const mul = (a, b) => web3Utils.toBN(a).mul(web3Utils.toBN(b))

export const gt = (a, b) => {
  return web3Utils.toBN(a).gt(web3Utils.toBN(b))
}

export const gte = (a, b) => {
  return web3Utils.toBN(a).gte(web3Utils.toBN(b))
}

export const lt = (a, b) => {
  return web3Utils.toBN(a).lt(web3Utils.toBN(b))
}

export const lte = (a, b) => {
  return web3Utils.toBN(a).lte(web3Utils.toBN(b))
}
