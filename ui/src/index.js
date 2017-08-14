import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider'
import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import { BrowserRouter, Redirect, Route } from 'react-router-dom'
import injectTapEventPlugin from 'react-tap-event-plugin'
import { createStore } from 'redux'

import Index from './components/Index'
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
          <Route
            path="/"
            render={() => (
              store.getState().auth.user ?
              <Index /> :
              <Redirect to="/login" />
            )}
          />
          <Route path="/login" component={Login} />
        </div>
      </BrowserRouter>
    </MuiThemeProvider>
  </Provider>
), document.getElementById('root'))
