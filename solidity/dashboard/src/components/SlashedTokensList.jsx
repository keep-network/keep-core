import React from 'react'
import moment from 'moment'
import { formatDate } from '../utils/general.utils'

const SlashedTokensList = (props) => {
  return (
    <div className="slashed-tokens-list">
      <div className="flex flex-1">
        <span className="text-label flex-2">SLASH EXPLANATION</span>
        <span className="text-label flex flex-1">AMOUNT (KEEP)</span>
        <span className="text-label flex flex-1">MIN STAKE (KEEP)</span>
      </div>
      <ul className="flex column flex-1">
        <li className="flex row flex-1" >
          <div className="details text-big flex-2">
            <p className="text-big">
              Group 12305162340123 was selected to do work and not enough members participated.
            </p>
            <span className="text-small text-grey">{formatDate(moment())}</span>
          </div>
          <span className="text-big text-dark-red flex-1">
            - 25
          </span>
          <span className="text-big flex-1">
            50
          </span>
        </li>
        <li className="flex row flex-1" >
          <div className="details text-big flex-2">
            <p className="text-big">
              Group 12305162340123 key was leaked. Private key was published outside of the members of the signing group.
            </p>
            <span className="text-small text-grey">{formatDate(moment())}</span>
          </div>
          <span className="text-big text-dark-red flex-1">
            - 25
          </span>
          <span className="text-big flex-1">
            50
          </span>
        </li>
      </ul>
    </div>
  )
}

export default SlashedTokensList
