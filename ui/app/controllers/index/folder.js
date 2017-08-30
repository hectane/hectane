import Ember from 'ember';

export default Ember.Controller.extend({
  store: Ember.inject.service('store'),

  actions: {
    test() {
      this.get('store').createRecord('message', {
        subject: "testing"
      }).save();
    }
  }
});
