import React from 'react'
import { RewardsGroupItem } from './RewardsGroupItem'

export const RewardsGroups = ({ groups }) => {
  return (
    <div className='group-items'>
      {groups.map(renderGroupItem)}
    </div>
  )
}

const renderGroupItem = (group, index) => <RewardsGroupItem key={group.groupIndex} {...group} />
