import React, { useState, useEffect }from 'react'
import { Tabs, Tab } from 'react-bootstrap';
import { useHistory } from 'react-router-dom'


export const RoutingTabs = (props) => {
    const [activeKey, setActiveKey] = useState('overview')
    const history = useHistory();
    useEffect(() => {
       setActiveKey(history.location.pathname.split('/')[1])
    }, []);

    useEffect(() => {
       history.push(activeKey);
    }, [activeKey]);

    const onSelect = (k) => {
        setActiveKey(k);
    }

    return (
        <>
            <Tabs activeKey={activeKey} onSelect={onSelect} id='dashboard-tabs' >
                <Tab eventKey='overview' title='Overview' />
                <Tab eventKey='stake' title='Stake' />
                <Tab eventKey='token-grants' title='Grant Token' />
                <Tab eventKey='create-grant-token' title='Create Token' />
            </Tabs>
            {props.children}
        </>
        
    );
}