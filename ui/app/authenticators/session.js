import Ember from 'ember';
import Base from 'ember-simple-auth/authenticators/base';

import ajax from '../util/ajax';

export default Base.extend({
  authenticate(username, password) {
    return ajax('/api/auth/login', {
      username: username,
      password: password
    });
  },

  invalidate() {
    return ajax('/api/auth/logout', {});
  },

  restore() {
    return Ember.RSVP.reject();
  }
});
