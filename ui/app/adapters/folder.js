import Ember from 'ember';
import DS from 'ember-data';

import ajax from '../util/ajax';

export default DS.Adapter.extend({
  createRecord(store, type, snapshot) {
    return ajax('/api/folders/new', this.serialize(snapshot));
  },

  findAll() {
    return ajax('/api/folders');
  }
});
