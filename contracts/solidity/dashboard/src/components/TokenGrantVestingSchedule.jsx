import React from 'react'
import Timeline from './Timeline'

const TokenGrantVestingSchedule = ({ }) => {
  return (
    <div>
      <div className="text-big text-darker-grey">
        Grant ID 12345
      </div>
      <div className="flex flex-row-space-between text-small text-grey">
        <div>
            Start Date 11/01/2019
          <div className="text-smaller">
            2 month cliff
          </div>
        </div>
        <div>
            Fully vested 11/01/2020
        </div>
      </div>
      <div className="mt-1">
        <div className="text-title text-darker-grey">total</div>
        <div>100 KEEP</div>
        <div>progress bar here</div>
      </div>
      <div className="mt-1">
        <Timeline
          title='schedule'
          breakePoints={[
            { dotColorClassName: 'grey', date: '11/11/2020', label: '10 Vested' },
            { dotColorClassName: 'grey', date: '12/11/2020', label: '10 Vested' },
            { dotColorClassName: 'grey', date: '13/11/2020', label: '10 Vested' },
            { dotColorClassName: '', date: '14/11/2020', label: '10 Released' },
          ]}
        />
        <div className="text-smaller text-darker-grey">
            Vesting will continue until completion on 11/01/2020
        </div>
      </div>
    </div>
  )
}

export default TokenGrantVestingSchedule
