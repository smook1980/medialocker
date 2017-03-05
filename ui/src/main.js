/* eslint-disable no-unused-vars */
/* eslint-disable no-multi-spaces */
/* eslint-disable no-console */

// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue';

import UIKit      from 'uikit/src/js/uikit';
import UIKitIcons from 'uikit/src/js/icons';

import App from './App';
import router from './router';

UIKit.use(UIKitIcons);
// window.UIKit = UIKit;

// components can be called from the imported UIkit reference
UIKit.notification('UIKit Workin\'');

Vue.config.productionTip = false;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App },
});
