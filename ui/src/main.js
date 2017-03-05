/* eslint-disable no-unused-vars */
/* eslint-disable no-multi-spaces */
/* eslint-disable no-console */
import UIKitStyle from 'uikit/dist/css/uikit.css';
// import theme   from 'uikit/dist/css/uikit.theme';

import JQuery from 'jquery';
import UIKit      from 'uikit/dist/js/uikit';
import UIKitIcons from 'uikit/dist/js/uikit-icons';

// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue';
import App from './App';
import router from './router';

UIKit.use(UIKitIcons);
window.JQuery = JQuery;
window.UIKit = UIKit;
console.dir(UIKit);
console.dir(UIKitIcons);

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
