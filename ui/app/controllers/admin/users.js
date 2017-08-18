import Ember from 'ember';

export default Ember.Controller.extend({
  store: Ember.inject.service('store'),

  actions: {
    delete(user) {
      user.destroyRecord();
    },

    showNewUser() {
      this.setProperties({
        errorMessage: null,
        username: '',
        password: '',
        isAdmin: false
      });
      Ember.$('#new-user-modal')
      .modal({
        closable: false,
        onApprove: () => {
          let { username, password, isAdmin } = this.getProperties('username', 'password', 'isAdmin');
          let record = this.get('store').createRecord('user', {
            username: username,
            password: password,
            is_admin: isAdmin
          });
          this.set('loading', true);
          record.save()
          .then(function() {
            Ember.$('#new-user-modal').modal('hide');
          }, (message) => {
            this.set('errorMessage', message);
            record.deleteRecord();
          })
          .finally(() => {
            this.set('loading', false);
          });
          return false;
        }
      })
      .modal('show');
    }
  }
});
