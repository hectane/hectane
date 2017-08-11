import React from 'react';
import {
  Card, CardActions, CardTitle, CardText
} from 'material-ui/Card';
import RaisedButton from 'material-ui/RaisedButton';
import TextField from 'material-ui/TextField';

/**
 * Present a login dialog for authentication
 */
export default class Login extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      username: '',
      password: ''
    };
  }

  handleUsernameChange = (event) => {
    this.setState({
      username: event.target.value
    });
  }

  handlePasswordChange = (event) => {
    this.setState({
      password: event.target.value
    });
  }

  handleTouchTap = (event) => {
    alert(this.state.username);
  }

  render() {
    return (
      <Card
        style={{margin: 'auto', maxWidth: 400}}
      >
        <CardTitle
          title="Login"
          subtitle="Enter your credentials to login."
        />
        <CardText>
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
            onTouchTap={this.handleTouchTap}
          />
        </CardActions>
      </Card>
    );
  }
}
