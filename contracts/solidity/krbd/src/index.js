import React from 'react';
import ReactDOM from 'react-dom'
import App from './components/App'
import './app.css'
import { Drizzle, generateStore } from "drizzle"
import Contracts from './contracts.js'

const options = { contracts: [] }
const drizzleStore = generateStore(options)
const drizzle = new Drizzle(options, drizzleStore)

// Async add contracts
Contracts.addContracts(drizzle)

ReactDOM.render(<App drizzle={drizzle} />, document.getElementById('root'))
