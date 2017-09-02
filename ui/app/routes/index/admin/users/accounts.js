import Ember from 'ember';

export default Ember.Route.extend({
  model(params) {
    return this.store.query('account', {user_id: params.user_id});
  }
});
