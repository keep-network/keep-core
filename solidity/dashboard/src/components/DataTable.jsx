import React from "react"
import { isEmptyArray } from "../utils/array.utils"
import SpeechBubbleTooltip from "./SpeechBubbleTooltip"
import Dropdown from "./Dropdown"

export class DataTable extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      headers: [],
    }
  }

  componentDidMount() {
    this.initializeDataTable()
  }

  initializeDataTable = () => {
    const headers = []
    React.Children.forEach(this.props.children, (children) => {
      headers.push({
        title: children.props.header,
        headerStyle: children.props.headerStyle,
      })
    })
    this.setState({ headers })
  }

  renderItemRow = (item, index) => {
    const { itemFieldId, cellStyle } = this.props
    return (
      <tr key={`${item[itemFieldId]}-${index}`}>
        {React.Children.map(this.props.children, (column, index) => {
          const {
            props: { field },
          } = column
          const cellKey = `${item[itemFieldId]}-${field}-${item[field]}-${index}`
          return (
            <td key={cellKey} style={cellStyle}>
              <span className="responsive-header">{column.props.header}</span>
              {this.renderColumnContent(column, item)}
            </td>
          )
        })}
      </tr>
    )
  }
  renderColumnContent = (column, item) => {
    if (!column.props.renderContent) {
      return item[column.props.field]
    }
    return column.props.renderContent(item)
  }

  renderHeader = ({ title, headerStyle }) => (
    <th key={title} style={headerStyle}>
      {title}
    </th>
  )

  render() {
    const {
      title,
      titleClassName,
      titleStyle,
      cellStyle,
      withTooltip,
      tooltipProps,
      subtitle,
    } = this.props

    return (
      <>
        <header className="table-header-wrapper">
          <div>
            <div className="flex center">
              <h4 className={titleClassName} style={titleStyle}>
                {title}
              </h4>
              {withTooltip && <SpeechBubbleTooltip {...tooltipProps} />}
            </div>
            <div className="text-grey-40 text-small">{subtitle}</div>
          </div>
          {this.props.withFilterDropdown && (
            <Dropdown
              withLabel={false}
              isFilterDropdow
              {...this.props.filterDropdownProps}
            />
          )}
        </header>
        <table>
          <thead>
            <tr>{this.state.headers.map(this.renderHeader)}</tr>
          </thead>
          <tbody>
            {isEmptyArray(this.props.data) ? (
              <tr className="text-center">
                <td colSpan={this.state.headers.length}>
                  <h4 className="text-grey-30">{this.props.noDataMessage}</h4>
                </td>
              </tr>
            ) : (
              this.props.data.map(this.renderItemRow)
            )}
          </tbody>
        </table>
      </>
    )
  }
}

DataTable.defaultProps = {
  noDataMessage: "No data.",
  titleClassName: "mr-1 text-grey-70",
  titleStyle: {},
  cellStyle: {},
  withTooltip: false,
  withFilterDropdown: false,
}

export const Column = ({ header, headerStyle, field, renderContent }) => null
Column.defaultProps = {
  headerStyle: {},
}
