import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    reply() {
      alert("Reply clicked.");
    },

    delete() {
      alert("Delete clicked.");
    }
  }
});
