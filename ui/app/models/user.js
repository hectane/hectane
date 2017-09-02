import DS from 'ember-data';

export default DS.Model.extend({
  username: DS.attr('string'),
  password: DS.attr('string'),
  isAdmin: DS.attr('boolean')
});
