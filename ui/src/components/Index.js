import React from 'react'
import { connect } from 'react-redux'

import Header from './Header';

/**
 * Container for the application
 */
class Index extends React.Component {
  render() {
    return (
      <div>
        <Header />
      </div>
    )
  }
}

export default connect()(Index)
