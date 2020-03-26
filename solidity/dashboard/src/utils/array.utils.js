import { isSameEthAddress } from '../utils/general.utils'

export const findIndexAndObject = (propertyName, value, array, comparationFunction = defaultComparisonFunction ) => {
  let indexInArray = null
  let obj = null
  for (let index = 0; index < array.length; index++) {
    const object = array[index]
    if (comparationFunction(object, propertyName, value)) {
      obj = object
      indexInArray = index
      break
    }
  }

  return { indexInArray, obj }
}

const defaultComparisonFunction = (object, propertyName, value) => object[propertyName] === value

export const compareEthAddresses = (object, propertyName, value) => isSameEthAddress(object[propertyName], value)
