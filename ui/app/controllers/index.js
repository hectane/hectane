import Ember from 'ember';

export default Ember.Controller.extend({
  session: Ember.inject.service('session'),
  store: Ember.inject.service('store'),

  // Make the user properties a little bit easier to reference
  user: Ember.computed('session', function() {
    return this.get('session.data.authenticated');
  }),

  actions: {
    logout() {
      this.get('session').invalidate();
    },

    showNewDialog() {
      this.setProperties({
        newFolderName: '',
        newVisible: true
      });
    },

    getNewPromise() {
      let newFolderName = this.get('newFolderName');
      let record = this.get('store').createRecord('folder', {
        name: newFolderName
      });
      return record.save();
    }
  }
});
