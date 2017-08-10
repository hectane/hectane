import React from 'react';
import ReactDOM from 'react-dom';
import {
  BrowserRouter,
  Route
} from 'react-router-dom';

import Login from './components/Login';

global.jQuery = require('jquery');
global.Tether = require('tether');

require('bootstrap');
require('bootstrap/dist/css/bootstrap.css');

ReactDOM.render((
  <BrowserRouter>
    <div>
      <Login />
      <Route path="/" />
      <Route path="/login" />
    </div>
  </BrowserRouter>
), document.getElementById('root'));
