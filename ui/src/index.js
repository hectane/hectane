import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';

import Sidebar from './Sidebar';

global.jQuery = require('jquery');
global.Tether = require('tether');

require('bootstrap');
require('bootstrap/dist/css/bootstrap.css');

ReactDOM.render((
  <BrowserRouter>
    <Sidebar />
  </BrowserRouter>
), document.getElementById('root'));
