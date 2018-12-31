import Vue from "vue";
import Router from "vue-router";

import Home from "@/views/Home.vue";
import Visualizer from "@/views/Visualizer.vue";
import E404 from "@/views/E404.vue";

const { BASE_URL } = process.env;

// Configure router.
Vue.use(Router);
export default new Router({
  mode: "history",
  base: BASE_URL,
  routes: [
    { path: "/", name: "home", component: Home },
    {
      path: "/r/:subreddit/",
      name: "visualizer",
      component: Visualizer,
      pathToRegexpOptions: { strict: true },
    },

    // Redirect trailing slashes.
    {
      path: "/r/:subreddit",
      redirect: "/r/:subreddit/",
      pathToRegexpOptions: { strict: true },
    },

    // Catch-all fallback to 404 page.
    { path: "*", name: "404", component: E404 },
  ],
});
