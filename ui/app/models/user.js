import DS from 'ember-data';

export default DS.Model.extend({
  username: DS.attr('string'),
  is_admin: DS.attr('boolean')
});
