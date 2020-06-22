import React from "react"

const DataTableSkeleton = ({ columns = 5, rows = 3 }) => {
  return (
    <section style={{ padding: "2rem" }}>
      <div className="skeleton mb-1" style={{ width: "300px" }} />
      <div className="text-small skeleton" style={{ width: "164px" }} />
      <table className="table__skeleton">
        <thead>
          {Array.from(Array(columns)).map((_, index) => (
            <th key={index} />
          ))}
        </thead>
        <tbody>
          {Array.from(Array(rows)).map((_, rowIndex) => (
            <tr key={`data-row-${rowIndex}`} className="datagrid__row">
              {Array.from(Array(columns)).map((_, index) => (
                <td key={`data-row-cell-${index}`}>
                  <div className="datagrid__loader" style={{ width: "70%" }} />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}

export default DataTableSkeleton
