import Ember from 'ember';

export default Ember.Component.extend({
  store: Ember.inject.service('store'),

  actions: {
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
