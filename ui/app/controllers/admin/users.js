import Ember from 'ember';

export default Ember.Controller.extend({
  store: Ember.inject.service('store'),

  actions: {
    showDeleteDialog(user) {
      this.setProperties({
        currentUser: user,
        deleteVisible: true
      });
    },

    getDeletePromise() {
      return this.get('currentUser').destroyRecord();
    },

    showNewDialog() {
      this.setProperties({
        username: '',
        password: '',
        isAdmin: false,
        newVisible: true
      });
    },

    getNewPromise() {
      let { username, password, isAdmin } = this.getProperties('username', 'password', 'isAdmin');
      let record = this.get('store').createRecord('user', {
        username: username,
        password: password,
        is_admin: isAdmin
      });
      return record.save();
    }
  }
});
