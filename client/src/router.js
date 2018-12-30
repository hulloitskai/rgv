import Vue from "vue";
import Router from "vue-router";

import Home from "@/views/Home.vue";
import Visualizer from "@/views/Visualizer.vue";
import E404 from "@/views/E404.vue";

// Configure router.
Vue.use(Router);
export default new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [
    { path: "/", name: "home", component: Home },
    { path: "/:subreddit", name: "visualizer", component: Visualizer },
    { path: "*", name: "404", component: E404 },
  ],
});
