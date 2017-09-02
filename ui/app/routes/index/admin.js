import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
  renderTemplate() {
    this.render();
    this.render('index.admin.sidebar', {
      into: 'index',
      outlet: 'sidebar'
    });
  }
});
