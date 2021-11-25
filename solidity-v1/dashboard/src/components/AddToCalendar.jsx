import React from "react"
import * as moment from "moment"

const makeTime = (time) => {
  return moment(time)
    .toISOString()
    .replace(/[-:]|\.\d{3}/g, "")
}

const makeUrl = (url, event) => {
  let finalUrl = url + "?"
  Object.entries(event).forEach(([key, value], index) => {
    if (!!value) {
      finalUrl = finalUrl + `${key}=${encodeURIComponent(value)}`
      if (index !== Object.entries(event).length - 1) {
        finalUrl = finalUrl + "&"
      }
    }
  })
  return finalUrl
}

const makeGoogleCalendarUrl = (event) => {
  return makeUrl("https://calendar.google.com/calendar/render", {
    action: "TEMPLATE",
    dates: `${makeTime(event.startsAt)}/${makeTime(event.endsAt)}`,
    location: event.location,
    text: event.name,
    details: event.details,
  })
}

const makeICSCalendarUrl = (event) => {
  const components = ["BEGIN:VCALENDAR", "VERSION:2.0", "BEGIN:VEVENT"]

  if (typeof document !== "undefined") {
    components.push(`URL:${document.URL}`)
  }

  components.push(
    `DTSTART:${makeTime(event.startsAt)}`,
    `DTEND:${makeTime(event.endsAt)}`,
    `SUMMARY:${event.name}`,
    `DESCRIPTION:${event.details}`,
    "END:VEVENT",
    "END:VCALENDAR"
  )

  return encodeURI(`data:text/calendar;charset=utf8,${components.join("\n")}`)
}

const AddToCalendar = ({ ...reactAddToCalendarProps }) => {
  const event = {
    name: "Coverage Pools - Tokens Ready To Claim",
    details: "You have 48 hours to claim your tokens!",
    startsAt: moment().format("YYYY-MM-DD HH:mm:ss"),
    endsAt: moment().add(2, "days").format("YYYY-MM-DD HH:mm:ss"),
  }
  console.log("google calnedar url", makeGoogleCalendarUrl(event))
  console.log("ics calendar api", makeICSCalendarUrl(event))
  return <></>
}

export default AddToCalendar
