import React from 'react'
import ReactDOM from 'react-dom'
import App from './App'
import registerServiceWorker from './registerServiceWorker'
import './css/app.css'

window.addEventListener('load', () => {
  ReactDOM.render(<App />, document.getElementById('root'))
  registerServiceWorker()
})
