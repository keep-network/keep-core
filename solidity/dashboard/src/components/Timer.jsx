import React, { useEffect, useState } from "react"
import moment from "moment"

const Timer = ({ targetInUnix }) => {
  const [remainingDays, setRemainingDays] = useState(0)
  const [remainingHours, setRemainingHours] = useState(0)
  const [remainingMinutes, setRemainingMinutes] = useState(0)
  const [remainingSeconds, setRemainingSeconds] = useState(0)

  useEffect(() => {
    const timerDuration = calculateTimerDuration(targetInUnix)
    setRemainingDays(timerDuration.days)
    setRemainingHours(timerDuration.hours)
    setRemainingMinutes(timerDuration.minutes)
    setRemainingSeconds(timerDuration.seconds)
  }, [targetInUnix])

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
            setRemainingMinutes(remainingMinutes - 1)
          } else if (remainingMinutes === 0) {
            setRemainingMinutes(59)
            if (remainingHours > 0) {
              setRemainingHours(remainingHours - 1)
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

  const zeroPad = (num, places) => String(num).padStart(places, "0")

  return (
    <>
      {zeroPad(remainingDays, 2)}:{zeroPad(remainingHours, 2)}:
      {zeroPad(remainingMinutes, 2)}:{zeroPad(remainingSeconds, 2)}
    </>
  )
}

const calculateTimerDuration = (targetInUnix) => {
  const currentDate = moment()
  const target = moment.unix(targetInUnix)
  const remainingTimeMs = target.diff(currentDate)
  const duration = moment.duration(remainingTimeMs)

  const days = Math.floor(duration.asDays())
  const hours = Math.floor(duration.asHours() % 24)
  const minutes = Math.floor(duration.asMinutes() % 60)
  const seconds = Math.floor(duration.asSeconds() % 60)

  return {
    days: days > 0 ? days : 0,
    hours: hours > 0 ? hours : 0,
    minutes: minutes > 0 ? minutes : 0,
    seconds: seconds > 0 ? seconds : 0,
  }
}

export default Timer
