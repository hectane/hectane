import DS from 'ember-data';

export default DS.Model.extend({
  time: DS.attr('date'),
  from: DS.attr('string'),
  to: DS.attr('string'),
  subject: DS.attr('string'),
  is_unread: DS.attr('boolean'),
  has_attachments: DS.attr('boolean')
});
