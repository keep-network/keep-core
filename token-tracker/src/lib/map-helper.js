export function mapToObject(map) {
  return Array.from(map).reduce((obj, [key, value]) => {
    obj[key] = value
    return obj
  }, {})
}
