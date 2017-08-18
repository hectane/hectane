import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),
  user: Ember.computed('session', function() {
    return this.get('session.data.authenticated');
  }),

  actions: {
    logout() {
      this.get('session').invalidate();
    }
  }
});
