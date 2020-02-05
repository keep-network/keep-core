import React from 'react'
import { useField } from 'formik'

const FromInput = ({ label, ...props }) => {
  const [field, meta] = useField(props.name, props.type)

  return (
    <div className="form-input flex flex-column">
      <label>
        {label}
      </label>
      <input {...field} {...props} />
      {meta.touched && meta.error ? (
        <div className='form-error'>{meta.error}</div>
      ) : null}
    </div>
  )
}

export default React.memo(FromInput)
