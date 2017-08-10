import React from 'react';

/**
 * Display the folders for the current user
 */
export default class Sidebar extends React.Component {
  render() {
    return (
      <ul className="nav flex-column">
        <li className="nav-item">
          <a className="nav-link" href="#">Inbox</a>
        </li>
      </ul>
    );
  }
}
