import React from 'react'
import Timeline from './Timeline'
import { useFetchData } from '../hooks/useFetchData'
import { tokenGrantsService } from '../services/token-grants.service'
import { displayAmount, formatDate } from '../utils'
import moment from 'moment'
import ProgressBar from './ProgressBar'
import { colors } from '../constants/colors'
import web3Utils from 'web3-utils'

const TokenGrantVestingSchedule = ({ grant }) => {
  const [state] = useFetchData(tokenGrantsService.fetchGrantVestingSchedule, [], grant.id)
  const { data } = state

  const cliffPeriod = moment.unix(grant.cliff).from(moment.unix(grant.start), true)
  const fullyVestedDate = moment.unix(grant.start).add(grant.duration, 'seconds')
  const unvested = web3Utils.toBN(grant.amount).sub(web3Utils.toBN(grant.vested))

  return (
    <div>
      <div className="text-big text-darker-grey">
        Grant ID {grant.id}
      </div>
      <div className="flex flex-row-space-between text-small text-grey">
        <div>
            Start Date {formatDate(moment.unix(grant.start))}
          <div className="text-smaller">
            {cliffPeriod} cliff
          </div>
        </div>
        <div>
            Fully vested {formatDate(fullyVestedDate)}
        </div>
      </div>
      <div className="mt-1">
        <div className="text-title text-darker-grey">total</div>
        <div className="text-big text-darker-grey">{displayAmount(grant.amount)} KEEP</div>
        <ProgressBar
          total={grant.amount}
          items={[
            { value: grant.vested, color: colors.darkGrey, label: 'Vested' },
            { value: unvested, color: colors.lightGrey, label: 'Unvested' },
          ]}
          withLegend
        />
      </div>
      <div className="mt-1">
        <Timeline
          title='schedule'
          breakpoints={data}
          footer={
            <div className="mb-3 text-smaller text-darker-grey">
              Vesting will continue until completion on {formatDate(fullyVestedDate)}
            </div>
          }
        />
      </div>
    </div>
  )
}

export default TokenGrantVestingSchedule
