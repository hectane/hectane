import React from 'react';
import ReactDOM from 'react-dom';
import {
  BrowserRouter,
  Route
} from 'react-router-dom';
import injectTapEventPlugin from 'react-tap-event-plugin';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';

import Login from './components/Login';

// Needed for onTouchTap
// http://stackoverflow.com/a/34015469/988941
injectTapEventPlugin();

ReactDOM.render((
  <MuiThemeProvider>
    <BrowserRouter>
      <div>
        <Login />
        <Route path="/" />
        <Route path="/login" />
      </div>
    </BrowserRouter>
  </MuiThemeProvider>
), document.getElementById('root'));
