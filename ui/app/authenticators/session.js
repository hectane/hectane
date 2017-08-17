import Ember from 'ember';
import Base from 'ember-simple-auth/authenticators/base';

/**
 * Session authenticator
 */
export default Base.extend({

  /**
   * Send a POST request to the API with login credentials
   */
  authenticate(username, password) {
    return new Ember.RSVP.Promise(function(resolve, reject) {
      Ember.$.post({
        url: '/api/auth/login',
        contentType: 'application/json;charset=utf-8',
        dataType: 'json',
        data: JSON.stringify({
          username: username,
          password: password
        })
      })
      .then(function(response) {
        Ember.run(null, resolve, response);
      }, function(xhr, status, error) {
        Ember.run(null, reject, xhr.responseText);
      });
    });
  },

  /**
   * Invalidate the session by logging out
   */
  invalidate() {
    return new Ember.RSVP.Promise(function(resolve, reject) {
      Ember.$.post('/api/auth/logout')
      .then(function(response) {
        Ember.run(null, resolve, response);
      }, function(xhr, status, error) {
        Ember.run(null, reject, xhr.responseText);
      });
    });
  },

  /**
   * Attempt to restore the session
   */
  restore() {
    return Ember.RSVP.reject();
  }
});
