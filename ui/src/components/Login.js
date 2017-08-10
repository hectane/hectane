import React from 'react';

/**
 * Present a login dialog for authentication
 */
export default class Login extends React.Component {
  handleSubmit(event) {
    event.preventDefault();
    alert('form was submitted!');
  }

  render() {
    return (
      <div className="container">
        <div className="card">
          <div className="card-block">
            <h4 className="card-title">Login</h4>
            <p className="card-text">
              Please enter your username and password to login.
            </p>
            <form onSubmit={this.handleSubmit}>
              <div className="form-group">
                <label>Username</label>
                <input type="text" className="form-control" />
              </div>
              <div className="form-group">
                <label>Password</label>
                <input type="password" className="form-control" />
              </div>
              <button type="submit" className="btn btn-primary">Login</button>
            </form>
          </div>
        </div>
      </div>
    );
  }
}
