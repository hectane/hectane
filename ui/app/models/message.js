import DS from 'ember-data';

export default DS.Model.extend({
  time: DS.attr('date'),
  from: DS.attr('string'),
  to: DS.attr('string'),
  subject: DS.attr('string'),
  isSeen: DS.attr('boolean'),
  isAnswered: DS.attr('boolean'),
  isFlagged: DS.attr('boolean'),
  isDeleted: DS.attr('boolean'),
  isDraft: DS.attr('boolean'),
  isRecent: DS.attr('boolean'),
  hasAttachments: DS.attr('boolean')
});
