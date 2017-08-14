import $ from 'jquery'
import { Card, CardActions, CardTitle, CardText } from 'material-ui/Card'
import RaisedButton from 'material-ui/RaisedButton'
import TextField from 'material-ui/TextField'
import React from 'react'
import { connect } from 'react-redux'
import { Redirect } from 'react-router-dom'

import { setUser } from '../actions/auth'

/**
 * Present a login dialog for authentication
 *
 * Upon successful login, the current user is added to the store and the user
 * is redirected to the home page.
 */
class Login extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      username: '',
      password: '',
      redirect: false,
      error: ''
    }
  }

  handleUsernameChange = (e) => {
    this.setState({username: e.target.value})
  }

  handlePasswordChange = (e) => {
    this.setState({password: e.target.value})
  }

  handleLogin = () => {
    $.post({
      url: '/api/auth/login',
      data: JSON.stringify({
        username: this.state.username,
        password: this.state.password
      }),
      contentType: 'application/json; charset=utf-8',
    })
    .then((d) => {
      this.props.dispatch(setUser(d.user))
      this.setState({redirect: true})
    })
    .fail((j, s, e) => this.setState({error: j.responseText}))
  }

  render() {
    if (this.state.redirect) {
      return (
        <Redirect to="/" />
      )
    }
    return (
      <Card
        style={{margin: 'auto', maxWidth: 400}}
      >
        <CardTitle
          title="Login"
          subtitle="Enter your credentials to login."
        />
        <CardText>
          {this.state.error && <p className="error">Error: {this.state.error}</p>}
          <TextField
            hintText="Username"
            fullWidth={true}
            onChange={this.handleUsernameChange}
          />
          <br />
          <TextField
            hintText="Password"
            fullWidth={true}
            type="password"
            onChange={this.handlePasswordChange}
          />
        </CardText>
        <CardActions>
          <RaisedButton
            label="Login"
            primary={true}
            onTouchTap={this.handleLogin}
          />
        </CardActions>
      </Card>
    )
  }
}

export default connect()(Login)
