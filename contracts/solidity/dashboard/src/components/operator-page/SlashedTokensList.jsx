import React from 'react'

const SlashedTokensList = (props) => {
  return (
    <>
        <div className="flex flex-1">
          <span className="text-label flex-2">DETAILS</span>
          <span className="text-label flex flex-1">SLASH AMOUNT</span>
          <span className="text-label flex flex-1">MIN STAKE</span>
        </div>
        <ul className="slashed-tokens-list flex flex-column flex-1">
          <li className="flex flex-row flex-1" >
            <div className="text-big flex-2">
              <span className="text-big">January 12, 2020</span>
              <p className="text-small text-grey">
                Group 12305162340123 was selected to do work and not enough members participated.
              </p>
            </div>
            <span className="text-big flex-1">1,000 K</span>
            <span className="text-big flex-1">50 K</span>
          </li>
          <li className="flex flex-row flex-1" >
            <div className="details text-big flex-2">
              <span className="text-big">January 12, 2020</span>
              <p className="text-small text-grey">
                Group 12305162340123 key was leaked. Private key was published outside of the members of the signing group.
              </p>
            </div>
            <span className="text-big flex-1">1,000 K</span>
            <span className="text-big flex-1">20 K</span>
          </li>
        </ul>
    </>

  )
}

export default SlashedTokensList
