import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider'
import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import { BrowserRouter, Route } from 'react-router-dom'
import injectTapEventPlugin from 'react-tap-event-plugin'
import { createStore } from 'redux'

import Login from './components/Login'
import rootReducer from './reducers'

import './index.css'

// Needed for onTouchTap
// http://stackoverflow.com/a/34015469/988941
injectTapEventPlugin()

const store = createStore(rootReducer)

ReactDOM.render((
  <Provider store={store}>
    <MuiThemeProvider>
      <BrowserRouter>
        <div>
          <Login />
          <Route path="/" />
          <Route path="/login" />
        </div>
      </BrowserRouter>
    </MuiThemeProvider>
  </Provider>
), document.getElementById('root'))
