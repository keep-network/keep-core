import React, { useEffect, useState } from "react"

const Timer = ({ days, hours, minutes, seconds }) => {
  const [remainingDays, setRemainingDays] = useState(days || 0)
  const [remainingHours, setRemainingHours] = useState(hours || 0)
  const [remainingMinutes, setRemainingMinutes] = useState(minutes || 0)
  const [remainingSeconds, setRemainingSeconds] = useState(seconds || 0)

  const zeroPad = (num, places) => String(num).padStart(places, "0")

  useEffect(() => {
    const myInterval = setInterval(() => {
      if (remainingSeconds > 0) {
        setRemainingSeconds(remainingSeconds - 1)
      }
      if (remainingSeconds === 0) {
        if (
          remainingDays === 0 &&
          remainingHours === 0 &&
          remainingMinutes === 0
        ) {
          clearInterval(myInterval)
        } else {
          setRemainingSeconds(59)
          if (remainingMinutes > 0) {
            setRemainingMinutes(minutes - 1)
          } else if (remainingMinutes === 0) {
            setRemainingMinutes(59)
            if (remainingHours > 0) {
              setRemainingHours(hours - 1)
            } else if (remainingHours === 0) {
              setRemainingHours(23)
              if (remainingDays > 0) {
                setRemainingDays(remainingDays - 1)
              }
            }
          }
        }
      }
    }, 1000)
    return () => {
      clearInterval(myInterval)
    }
  })

  return (
    <>
      {zeroPad(remainingDays, 2)}:{zeroPad(remainingHours, 2)}:
      {zeroPad(remainingMinutes, 2)}:{zeroPad(remainingSeconds, 2)}
    </>
  )
}

export default Timer
