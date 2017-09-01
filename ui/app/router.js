import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType,
  rootURL: config.rootURL
});

Router.map(function() {
  this.route('login');

  this.route('index', {path: '/'}, function() {
    this.route('folder', {path: 'folder/:folder_id'});
    this.route('message', {path: 'message/:message_id'});
  });
  this.route('admin', function() {
    this.route('users');
    this.route('log');
    this.route('domains');
  });
});

Router.reopen({
  location: 'hash'
});

export default Router;
