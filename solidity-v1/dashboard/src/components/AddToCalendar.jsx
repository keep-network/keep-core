import React from "react"
import * as moment from "moment"
import { CompoundDropdown as Dropdown } from "./Dropdown"
import * as Icons from "./Icons"
import { ADD_TO_CALENDAR_OPTIONS } from "../constants/constants"

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

const options = [
  {
    label: "Google Calendar",
    value: ADD_TO_CALENDAR_OPTIONS.GOOGLE_CALENDER,
  },
  {
    label: "Apple Calendar",
    value: ADD_TO_CALENDAR_OPTIONS.APPLE_CALENDAR,
  },
]

const AddToCalendar = ({
  name = "New Event",
  details = "Event details",
  startsAt = moment().format("YYYY-MM-DD HH:mm:ss"), // date in YYY-MM-DD HH:mm:ss format
  endsAt = moment().add(2, "hours").format("YYYY-MM-DD HH:mm:ss"), // date in YYY-MM-DD HH:mm:ss format
  className = "",
}) => {
  const event = {
    name,
    details,
    startsAt,
    endsAt,
  }

  const onCalendarSelect = (selectedCalendar) => {
    switch (selectedCalendar) {
      case ADD_TO_CALENDAR_OPTIONS.GOOGLE_CALENDER: {
        const win = window.open(makeGoogleCalendarUrl(event))
        win.focus()
        break
      }
      case ADD_TO_CALENDAR_OPTIONS.APPLE_CALENDAR: {
        const win = window.open(makeICSCalendarUrl(event))
        win.focus()
        break
      }
      default:
        return
    }
  }

  return (
    <Dropdown
      selectedItem={options[0]}
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
        {options.map((option) => {
          return (
            <Dropdown.Option
              key={`dropdown-${option.value}`}
              value={option.value}
            >
              {option.label}
            </Dropdown.Option>
          )
        })}
      </Dropdown.Options>
    </Dropdown>
  )
}

export default AddToCalendar
