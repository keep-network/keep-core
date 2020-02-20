import React from 'react'
import { formatDate } from '../utils'
import moment from 'moment'

const SelectedGrantDropdown = ({ grant }) => {
  return (
    <div className="flex flex-column">
      <div className="text-smaller">
        Grant issued {grant.start && formatDate(moment.unix(grant.start).add(grant.duration, 'seconds'))}
      </div>
      <div className="text-smaller text-grey">
        Grant ID {grant.id}
      </div>
    </div>
  )
}

export default React.memo(SelectedGrantDropdown)
