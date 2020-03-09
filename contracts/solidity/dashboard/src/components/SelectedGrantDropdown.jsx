import React from 'react'
import { formatDate } from '../utils/general.utils'
import moment from 'moment'

const SelectedGrantDropdown = ({ grant }) => {
  return (
    <div className="flex column">
      <div className="text-smaller">
        Issued on {grant.start && formatDate(moment.unix(grant.start).add(grant.duration, 'seconds'))}
      </div>
      <div className="text-smaller text-grey-60">
        Grant ID {grant.id}
      </div>
    </div>
  )
}

export default React.memo(SelectedGrantDropdown)
