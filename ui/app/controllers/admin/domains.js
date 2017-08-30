import Ember from 'ember';

export default Ember.Controller.extend({
  store: Ember.inject.service('store'),

  actions: {
    showDeleteDialog(domain) {
      this.setProperties({
        currentDomain: domain,
        deleteVisible: true
      });
    },

    getDeletePromise() {
      return this.get('currentDomain').destroyRecord();
    },

    showNewDialog() {
      this.setProperties({
        name: '',
        newVisible: true
      });
    },

    getNewPromise() {
      let record = this.get('store').createRecord('domain', {
        name: this.get('name')
      });
      return record.save();
    }
  }
});
