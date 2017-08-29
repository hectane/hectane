import Ember from 'ember';

export default Ember.Route.extend({
  model(params) {
    return this.store.query('message', {folder_id: params.folder_id});
  }
});
