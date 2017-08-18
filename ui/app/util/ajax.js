import Ember from 'ember';

/**
 * Create a promise for the specified API call
 * @param {String} url
 * @param {Object} [data]
 * @return {Ember.RSVP}
 */
export default function(url) {
  if (typeof url !== 'string') {
    throw new Error('first parameter must be a string');
  }
  let args = {
    url: url,
    dataType: 'json'
  };
  if (typeof arguments[1] === 'object') {
    args.type = 'POST';
    args.contentType = 'application/json;charset=utf-8';
    args.data = JSON.stringify(arguments[1]);
  }
  return new Ember.RSVP.Promise(function(resolve, reject) {
    return Ember.$.ajax(args)
    .then(function(response) {
      Ember.run(null, resolve, response);
    }, function(xhr, status, error) {
      Ember.run(null, reject, xhr.responseText);
    });
  });
};
