import Vue from "vue";
import Router from "vue-router";

import Home from "@/views/Home.vue";
import Visualizer from "@/views/Visualizer.vue";

// Configure router.
Vue.use(Router);
export default new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [
    { path: "/", name: "home", component: Home },
    { path: "/:subreddit", name: "visualizer", component: Visualizer },
  ],
});
