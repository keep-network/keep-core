export const findIndexAndObject = (propertyName, value, array) => {
  let indexInArray = null
  let obj = null
  for (let index = 0; index < array.length; index++) {
    const object = array[index]
    if (object[propertyName] === value) {
      obj = object
      indexInArray = index
      break
    }
  }

  return { indexInArray, obj }
}
