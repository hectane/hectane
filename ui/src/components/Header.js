import $ from 'jquery'
import AppBar from 'material-ui/AppBar';
import IconButton from 'material-ui/IconButton';
import IconMenu from 'material-ui/IconMenu';
import MenuItem from 'material-ui/MenuItem';
import MoreVertIcon from 'material-ui/svg-icons/navigation/more-vert';
import React from 'react'
import { connect } from 'react-redux'
import { Redirect } from 'react-router-dom'

import { setUser } from '../actions/auth'

class Header extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      redirect: false
    }
  }

  // TODO: handle logout error

  handleLogout = () => {
    $.post('/api/auth/logout')
    .then(() => {
      this.props.dispatch(setUser(undefined))
      this.setState({redirect: true})
    })
  }

  render() {
    if (this.state.redirect) {
      return (
        <Redirect to="/login" />
      )
    }
    return (
      <AppBar title="Hectane"
        iconElementRight={
          <IconMenu
            iconButtonElement={<IconButton><MoreVertIcon /></IconButton>}
          >
            <MenuItem
              primaryText="Logout"
              onTouchTap={this.handleLogout}
            />
          </IconMenu>
        }
      />
    )
  }
}

export default connect()(Header)
