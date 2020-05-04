import React from "react"
import { formatDate } from "../utils/general.utils"
import moment from "moment"

const SelectedGrantDropdown = ({ grant }) => {
  return (
    <div>
      <div className="text-big text-black">Grant ID {grant.id}</div>
      <div
        className="text-instruction text-grey-50"
        style={{ fontSize: "0.7rem" }}
      >
        Issued on{" "}
        {grant.start &&
          formatDate(moment.unix(grant.start).add(grant.duration, "seconds"))}
      </div>
    </div>
  )
}

export default React.memo(SelectedGrantDropdown)
