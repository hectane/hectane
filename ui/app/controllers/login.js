import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),

  actions: {
    login() {
      this.set('loading', true);
      let { username, password } = this.getProperties('username', 'password');
      this.get('session').authenticate('authenticator:session', username, password)
      .catch((message) => this.set('errorMessage', message))
      .finally(() => this.set('loading', false));
    }
  }
});
