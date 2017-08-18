import Ember from 'ember';
import DS from 'ember-data';

export default DS.Adapter.extend({
  createRecord(store, type, snapshot) {
    let data = this.serialize(snapshot);
    return new Ember.RSVP.Promise(function(resolve, reject) {
      Ember.$.post({
        url: '/api/folders/new',
        contentType: 'application/json;charset=utf-8',
        dataType: 'json',
        data: JSON.stringify(data)
      })
      .then(function(response) {
        Ember.run(null, resolve, response);
      }, function(xhr, status, error) {
        Ember.run(null, reject, xhr.responseText);
      });
    });
  },

  findAll() {
    return new Ember.RSVP.Promise(function(resolve, reject) {
      Ember.$.getJSON('/api/folders')
      .then(function(response) {
        Ember.run(null, resolve, response.folders);
      }, function(xhr, status, error) {
        Ember.run(null, reject, xhr.responseText);
      });
    });
  }
});
