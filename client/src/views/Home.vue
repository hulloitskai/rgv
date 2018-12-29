<template>
  <div class="home">
    <h1>Home</h1>
    <form @submit.prevent="gotoVisualizer">
      <label for="subreddit">Subreddit:</label>
      <input
        id="subreddit"
        type="text"
        v-model="subreddit"
        placeholder="uwaterloo"
      />
      <button type="submit">Visualize</button>
    </form>
  </div>
</template>

<script>
import _ from "lodash";

export default {
  data: () => ({ subreddit: "" }),
  methods: {
    gotoVisualizer() {
      // Validate subreddit.
      if (!this.subreddit) alert("No subreddit name given.");
      if (_.startsWith(this.subreddit, "r/")) {
        this.subreddit = this.subreddit.slice(2);
      }

      this.$router.push({
        name: "visualizer",
        params: { subreddit: this.subreddit },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
form {
  margin-top: 10px;
  width: 200px;
  display: flex;
  flex-direction: column;

  input {
    width: 100%;
    box-sizing: border-box;
  }
}
</style>
