import React from 'react'
import { SubmitButton } from './Button'

const InlineForm = ({ inputProps, onSubmit, classNames }) => {
  return (
    <form className={`inline-form ${classNames}`} onSubmit={onSubmit}>
      <input {...inputProps}/>
      <SubmitButton
        type="submit"
        className="btn btn-primary btn-large"
        onSubmitAction={onSubmit}
      >
        UNDELEGATE
      </SubmitButton>
    </form>
  )
}

export default React.memo(InlineForm)
