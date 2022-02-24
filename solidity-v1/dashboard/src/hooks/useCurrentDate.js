import { useEffect, useState } from "react"
import moment from "moment"

/**
 * Uses current date in unix updated every second
 * @return {number} - current date in unix timestamp
 */
const useCurrentDate = () => {
  const [currentDateInUnix, setCurrentDateInUnix] = useState(moment().unix())

  useEffect(() => {
    const myInterval = setInterval(() => {
      setCurrentDateInUnix(moment().unix())
    }, 1000)
    return () => {
      clearInterval(myInterval)
    }
  })

  return currentDateInUnix
}

export default useCurrentDate
