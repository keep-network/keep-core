import React from 'react'
import { formatDate, displayAmount } from '../utils'
import { SubmitButton } from './Button'
import { colors } from '../constants/colors'
import { CircularProgressBars } from './CircularProgressBar'
import moment from 'moment'

const TokenGrantOverview = ({ selectedGrant }) => {
  const cliffPeriod = moment
    .unix(selectedGrant.cliff)
    .from(moment.unix(selectedGrant.start), true)
  const fullyUnlockedDate = moment
    .unix(selectedGrant.start)
    .add(selectedGrant.duration, 'seconds')

  return (
    <div className="token-grant-overview">
      <div>
        <h2 className="balance">{displayAmount(selectedGrant.amount)}&nbsp;KEEP</h2>
        <div className="text-small text-grey-40">
          Issued: 01/01/2020
          <span className="text-smaller text-grey-30">&nbsp;Cliff: {cliffPeriod}</span>
          <br/>
          Fully Unlocked: {formatDate(fullyUnlockedDate)}
        </div>
      </div>
      <hr/>
      <div className="flex row full-center">
        <div>
          <CircularProgressBars
            total={selectedGrant.amount}
            items={[
              {
                value: selectedGrant.vested,
                backgroundStroke: '#D7F6EE',
                color: colors.primary,
                label: 'Unlocked',
              },
              {
                value: selectedGrant.released,
                color: colors.secondary,
                withBackgroundStroke: false,
                radius: 48,
                label: 'Relesed',
              },
            ]}
            withLegend
          />
        </div>
        <div className={`ml-2 mt-1 ${selectedGrant.readyToRealese === '0' ? 'self-start' : '' }`}>
          <div className="text-label">
            unlocked
          </div>
          <div className="text-label">
            {displayAmount(selectedGrant.vested)}
            <div className="text-smaller text-grey-40">
              of {displayAmount(selectedGrant.amount)} total
            </div>
          </div>
          {
            selectedGrant.readyToRealese !== '0' &&
            <div className="mt-2">
              <div>
                <span className="text-secondary text-bold text-small">
                  {displayAmount(selectedGrant.readyToRealese)}
                &nbsp;
                </span>
                <span className="text-smaller text-grey-40">Ready to realese</span>
              </div>
              <SubmitButton
                className="btn btn-sm btn-secondary"
                onSubmitAction={() => {
                  console.log('on submit action here')
                }}
              >
              release tokens
              </SubmitButton>
            </div>
          }
        </div>
      </div>
      <div className="flex row full-center mt-1">
        <div>
          <CircularProgressBars
            total={selectedGrant.amount}
            items={[
              {
                value: selectedGrant.staked,
                backgroundStroke: '#F8E9D3',
                color: colors.brown,
                label: 'Staked',
              },
            ]}
            withLegend
          />
        </div>
        <div className="ml-2 mt-1 self-start">
          <div className="text-label">
            staked
          </div>
          <div className="text-label">
            {displayAmount(selectedGrant.staked)}
            <div className="text-smaller text-grey-40">
              of {displayAmount(selectedGrant.amount)} total
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default TokenGrantOverview
