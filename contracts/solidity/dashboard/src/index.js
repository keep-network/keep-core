import 'bootstrap/dist/css/bootstrap.css'
import React from 'react'
import ReactDOM from 'react-dom'
import App from './App'
import registerServiceWorker from './registerServiceWorker'
import './app.css'
import { getWeb3 } from './utils'

window.addEventListener('load', async () => {
  await getWeb3()
  ReactDOM.render(<App />, document.getElementById('root'))
  registerServiceWorker()
})
