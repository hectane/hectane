import Ember from 'ember';

/**
 * Modal dialog for async operations
 *
 * The component is expected to receive these properties:
 *
 *  - title - used for the UI
 *  - visible - used to show and hide the modal
 *  - loading - used for determining the state of the request
 *  - getPromise - action for retrieving a promise
 */
export default Ember.Component.extend({
  classNames: ['ui', 'mini', 'modal'],

  init() {
    this._super(...arguments);
    this.addObserver('visible', () => {
      this.$().modal(this.get('visible') ? 'show' : 'hide');
    });
  },

  didInsertElement() {
    this.$().modal({
      closable: false,
      duration: 0,
      onShow: () => {
        this.set('errorMessage', null);
      },
      onApprove: () => {
        this.set('loading', true);
        this.get('getPromise')()
        .then(() => {
          this.set('visible', false);
        }, (message) => {
          this.set('errorMessage', message);
        })
        .finally(() => {
          this.set('loading', false);
        });
        return false;
      },
      onDeny: () => {
        this.set('visible', false);
        return false;
      }
    });
  }
});
