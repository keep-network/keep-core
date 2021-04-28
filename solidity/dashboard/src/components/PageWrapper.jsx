import React from "react"
import { Switch } from "react-router-dom"
import Header from "./Header"
import { renderExplorerModePage, renderPage } from "./Routing"
import { isEmptyArray } from "../utils/array.utils"
import Redirect from "./Redirect"

const PageWrapper = ({
  title,
  routes,
  children,
  headerClassName = "",
  newPage = false,
}) => {
  const hasRoutes = !isEmptyArray(routes)

  return (
    <>
      <Header
        title={title}
        subLinks={hasRoutes ? routes.map((_) => _.route) : []}
        className={headerClassName}
        newPage={newPage}
      />
      <main>
        {children}
        {hasRoutes && (
          <Switch>
            {routes.map(renderPage)}
            {routes.map(renderExplorerModePage)}
            <Redirect to={routes[0].route.path} />
          </Switch>
        )}
      </main>
    </>
  )
}
export default PageWrapper
