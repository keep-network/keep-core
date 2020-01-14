import React, { useState, useEffect } from 'react'
import { Tabs, Tab } from 'react-bootstrap'
import { useHistory } from 'react-router-dom'

export const RoutingTabs = (props) => {
  const history = useHistory()
  const [activeKey, setActiveKey] = useState(history.location.pathname.split('/')[1])

  useEffect(() => {
    setActiveKey(history.location.pathname.split('/')[1])
  }, [])

  useEffect(() => {
    history.push(activeKey)
  }, [activeKey])

  const onSelect = (k) => {
    setActiveKey(k)
  }

  const isOwner = !props.isOperator && props.isTokenHolder

  return (

    <>
      <Tabs activeKey={activeKey} onSelect={onSelect} id='dashboard-tabs' >
        <Tab eventKey='overview' title='Overview' />
        { isOwner && <Tab eventKey='stake' title='Stake' /> }
        { isOwner && <Tab eventKey='token-grants' title='Token Grants' /> }
        { isOwner && <Tab eventKey='create-token-grants' title='Create Token Grant' /> }
      </Tabs>
      {props.children}
    </>

  )
}
