import DS from 'ember-data';

import ajax from '../util/ajax';

export default DS.Adapter.extend({
  createRecord(store, type, snapshot) {
    return ajax('/api/admin/users/new', this.serialize(snapshot));
  },

  deleteRecord(store, type, snapshot) {
    return ajax(`/api/admin/users/${snapshot.id}/delete`, {});
  },

  findAll() {
    return ajax('/api/admin/users');
  }
});
