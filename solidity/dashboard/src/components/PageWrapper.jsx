import React from "react"
import { Switch, Redirect } from "react-router-dom"
import Header from "./Header"
import { renderPage } from "./Routing"
import { isEmptyArray } from "../utils/array.utils"

const PageWrapper = ({ title, routes, children, headerClassName = "" }) => {
  const hasRoutes = !isEmptyArray(routes)

  return (
    <>
      <Header
        title={title}
        subLinks={hasRoutes ? routes.map((_) => _.route) : []}
        className={headerClassName}
      />
      <main>
        {children}
        {hasRoutes && (
          <Switch>
            {routes.map(renderPage)}
            <Redirect to={routes[0].route.path} />
          </Switch>
        )}
      </main>
    </>
  )
}
export default PageWrapper
