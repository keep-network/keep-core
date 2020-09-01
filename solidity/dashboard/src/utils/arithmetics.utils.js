import web3Utils from "web3-utils"

const ONE_HUNDRED = web3Utils.toBN(100)

export const add = (a, b) => {
  return web3Utils.toBN(a).add(web3Utils.toBN(b))
}

export const sub = (a, b) => {
  return web3Utils.toBN(a).sub(web3Utils.toBN(b))
}

export const mul = (a, b) => web3Utils.toBN(a).mul(web3Utils.toBN(b))

export const div = (a, b) => web3Utils.toBN(a).div(web3Utils.toBN(b))

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

export const isZero = (a) => web3Utils.toBN(a).isZero()

export const percentageOf = (value, total) => {
  if (isZero(total)) {
    return 0
  }

  return web3Utils.toBN(value).mul(ONE_HUNDRED).div(total)
}
