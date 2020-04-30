import React from "react"

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
      headers.push(children.props.header)
    })
    this.setState({ headers })
  }

  renderItemRow = (item, index) => {
    const { itemFieldId } = this.props
    return (
      <tr key={`${item[itemFieldId]}-${index}`}>
        {React.Children.map(this.props.children, (column, index) => {
          const {
            props: { field },
          } = column
          const cellKey = `${item[itemFieldId]}-${field}-${item[field]}-${index}`
          return (
            <td key={cellKey}>
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

  renderHeader = (header) => <th key={header}>{header}</th>

  render() {
    return (
      <table>
        <thead>
          <tr>{this.state.headers.map(this.renderHeader)}</tr>
        </thead>
        <tbody>{this.props.data.map(this.renderItemRow)}</tbody>
      </table>
    )
  }
}

export const Column = ({ header, field, renderContent }) => null
