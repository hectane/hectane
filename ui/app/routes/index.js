import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
  model() {
    return this.store.findAll('folder');
  },

  renderTemplate() {
    this.render();
    this.render('index.sidebar', {
      into: 'index',
      outlet: 'sidebar'
    });
  },

  actions: {

    // Show the dimmer and spinning circle for all transitions
    loading(transition) {
      let loader = Ember.$('#loader').addClass('active');
      transition.promise.finally(function() {
        loader.removeClass('active');
      });
    }
  }
});
