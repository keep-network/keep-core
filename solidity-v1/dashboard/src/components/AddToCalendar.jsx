import React, { useMemo } from "react"
import moment from "moment"
import { CompoundDropdown as Dropdown } from "./Dropdown"
import * as Icons from "./Icons"
import { ADD_TO_CALENDAR_OPTIONS } from "../constants/constants"
import { useWeb3Context } from "./WithWeb3Context"

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
    `LOCATION:${event.details}`,
    "END:VEVENT",
    "END:VCALENDAR"
  )

  return encodeURI(`data:text/calendar;charset=utf8,${components.join("\n")}`)
}

const calendars = [
  {
    label: "Google Calendar",
    value: ADD_TO_CALENDAR_OPTIONS.GOOGLE_CALENDER,
    calendarEventBuilder: makeGoogleCalendarUrl,
  },
  {
    label: "Apple Calendar",
    value: ADD_TO_CALENDAR_OPTIONS.APPLE_CALENDAR,
    calendarEventBuilder: makeICSCalendarUrl,
  },
]

const AddToCalendar = ({
  name = "New Event",
  details = "Event details",
  location = "http://localhost:3000/overview",
  startsAt = moment().format("YYYY-MM-DD HH:mm:ss"), // date in YYY-MM-DD HH:mm:ss format
  endsAt = moment().add(2, "hours").format("YYYY-MM-DD HH:mm:ss"), // date in YYY-MM-DD HH:mm:ss format
  className = "",
}) => {
  const { yourAddress } = useWeb3Context()
  const locationWithInsertedAddress = useMemo(() => {
    if (!yourAddress) return location
    return location.replace("${address}", yourAddress)
  }, [yourAddress, location])

  const event = {
    name,
    details,
    location: locationWithInsertedAddress,
    startsAt,
    endsAt,
  }

  const onCalendarSelect = (selectedCalendar) => {
    const win = window.open(selectedCalendar.calendarEventBuilder(event))
    win.focus()
  }

  return (
    <Dropdown
      selectedItem={{}}
      onSelect={onCalendarSelect}
      comparePropertyName="label"
      className={`add-to-calendar-dropdown ${className}`}
      rounded
    >
      <Dropdown.Trigger
        className={"add-to-calendar-dropdown__trigger"}
        withTriggerArrow={false}
      >
        <div className="flex row center">
          <Icons.Calendar width={14} height={14} />
          <span className="text-label text-label--without-hover text-black ml-1">
            ADD TO CALENDAR
          </span>
        </div>
      </Dropdown.Trigger>
      <Dropdown.Options className="add-to-calendar-dropdown__options">
        {calendars.map((calendar) => {
          return (
            <Dropdown.Option
              key={`dropdown-${calendar.value}`}
              value={calendar}
            >
              {calendar.label}
            </Dropdown.Option>
          )
        })}
      </Dropdown.Options>
    </Dropdown>
  )
}

export default AddToCalendar
