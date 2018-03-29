import React, { Component, PropTypes } from 'react'
import Web3 from 'web3'
import { Link } from 'react-scroll'
import { prefixLink } from 'gatsby-helpers'
import { config } from 'config'
import { Menu, Container, Label, Segment, Grid, Icon } from 'semantic-ui-react'
import '../css/style.scss'
import sortBy from 'sort-by'
import * as Icons from '../components/icons';

export default class Index extends Component {
  constructor(props) {
    super(props)
    if (config.interaction && config.interaction.providerUrl) {
      this.web3 = new Web3()
      this.web3.setProvider(new this.web3.providers.HttpProvider(config.interaction.providerUrl))
    }
  }
  getChildContext() {
    return { web3: this.web3 }
  }
  render() {
    const onIndex = prefixLink('/') === this.props.location.pathname
    const childRoutes = this.props.route && this.props.route.childRoutes
    const docsPath = childRoutes && childRoutes[0] && childRoutes[0].path
    const menu = childRoutes.map((child) => {
      if (child.path === this.props.location.pathname) {
        return (
          <Menu text vertical key={child.page.data.name} className="sidebarMenu">
            {child.page.data.abiDocs.sort(sortBy('type', 'name')).map(method => {
              if (method.name) return (
                <Menu.Item as={Link} 
                  name={method.name}
                  key={`${child.page.data.name}${method.name}`}
                  to={`${child.page.data.name}${method.name}`}
                  isDynamic={false}
                  duration={500}
                  smooth={true}
                  containerId="mainscroll">
                  {method.name}
                </Menu.Item>
              )
            })}
          </Menu>
        )
      }
    })
    return (
      <div style={{ paddingTop: '120px' }} className="pusher">
        <Menu borderless fixed="top" className="navbar">
          <Container>
            <Menu.Item as={'a'} href={`${config.baseUrl}`} className="brand">
              <div>
                <Icons.Keep width="160px" height="42px" />
                <p>Contracts Documentation</p>
              </div>
            </Menu.Item>
            <Menu.Menu position="right" className="">
              {childRoutes.map((child) => {
                return (
                  <Menu.Item key={child.page.data.name} as={'a'} href={'./docs/'+child.page.data.name+'/'}>
                    {child.page.data.name}
                  </Menu.Item>
                )
              })}
            </Menu.Menu>
          </Container>
        </Menu>
        <Container>
          <Grid>
            <Grid.Row>
              {!onIndex && <Grid.Column width={3}><div className="scrollable">{menu}</div></Grid.Column>}
              <Grid.Column width={13}>
                <div className="scrollable" id="mainscroll">{this.props.children}</div>
              </Grid.Column>
            </Grid.Row>
          </Grid>
        </Container>
        <Container className="footer">
          <Segment secondary size="small" attached="top" compact>
            <Grid stackable>
              <Grid.Row>
                <Grid.Column width={6}>
                  <b>&copy; {config.author}</b> - {config.license}, {new Date(config.buildTime).getFullYear()}
                </Grid.Column>
                <Grid.Column width={10} textAlign="right">
                  Docs built using <b>Solidity {config.compiler}</b> on <b>{new Date(config.buildTime).toLocaleDateString()}</b>
                </Grid.Column>
              </Grid.Row>
            </Grid>
          </Segment>
        </Container>
      </div>
    )
  }
}
Index.childContextTypes = {
  web3: PropTypes.object,
}
Index.propTypes = {
  children: PropTypes.object,
  location: PropTypes.object,
  route: PropTypes.object,
}
