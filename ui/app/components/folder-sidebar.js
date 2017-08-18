import Ember from 'ember';

export default Ember.Component.extend({
  store: Ember.inject.service('store'),

  actions: {
    showNewFolder() {
      this.setProperties({
        errorMessage: null,
        newFolderName: ''
      });
      Ember.$('#new-folder-modal')
      .modal({
        closable: false,
        onApprove: () => {
          let newFolderName = this.get('newFolderName');
          let record = this.get('store').createRecord('folder', {
            name: newFolderName
          });
          this.set('loading', true);
          record.save()
          .then(function() {
            Ember.$('#new-folder-modal').modal('hide');
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
