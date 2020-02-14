import React, { useEffect, useRef } from 'react'
import ReactDOM from 'react-dom'
import * as Icons from './Icons'

const modalRoot = document.getElementById('modal-root')

const Modal = React.memo(({ closeModal, ...props }) => {
  const modalOverlay = useRef(null)
  useEffect(() => {
    document.body.style.overflow = 'hidden'

    return () => {
      document.body.style.overflow = 'scroll'
    }
  }, [])

  const onOverlayClick = (event) => {
    if (modalOverlay.current === event.target) {
      closeModal()
    }
  }

  return ReactDOM.createPortal(
    <div ref={modalOverlay} className="modal-overlay" onClick={onOverlayClick}>
      <div className="modal-content">
        <div className="flex">
          <h5 className="text-darker-grey">{props.title}</h5>
          <div className="modal-close" onClick={closeModal}>
            <Icons.Cross width={18} height={18} />
          </div>
        </div>
        {props.children}
      </div>
    </div>,
    modalRoot
  )
})


export default Modal
