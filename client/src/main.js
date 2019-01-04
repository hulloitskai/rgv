import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

// Install plugins:
import Tooltip from "v-tooltip";
Vue.use(Tooltip);

import Meta from "vue-meta";
Vue.use(Meta);

// Uncomment to enable PWA capabilities:
//import "./registerServiceWorker";

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount("#app");
