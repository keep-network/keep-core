import { writeFileSync } from "fs"
import { logger } from "../lib/winston.js"

/**
 * Stores data to a file. If the file already exists it overwrites it.
 * @param {*} data Data to store.
 * @param {String} filePath Destination file path.
 */
export function dumpDataToFile(data, filePath) {
  logger.info(`dump data to a file: ${filePath}`)

  if (data instanceof Map) {
    data = mapToObject(data)
  }
  if (data instanceof Set) {
    data = Array.from(data)
  }

  writeFileSync(filePath, JSON.stringify(data, null, 2))
}

function mapToObject(map) {
  return Array.from(map).reduce((obj, [key, value]) => {
    obj[key] = value
    return obj
  }, {})
}
