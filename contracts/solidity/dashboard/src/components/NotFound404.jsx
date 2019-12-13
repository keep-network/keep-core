import React from 'react'
import { useLocation } from 'react-router-dom'

export const NotFound404 = () => {
    const location = useLocation()
  
    return (
      <div>
        <h3>
          No match for <code>{location.pathname}</code>
        </h3>
      </div>
    );
  }