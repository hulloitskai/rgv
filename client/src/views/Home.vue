<template>
  <div class="home flex col center">
    <div class="card flex col center">
      <img class="logo" src="@/assets/logo.svg" />
      <div class="title flex col center">
        <h1 class="name">RGV</h1>
        <h5 class="description">Reddit Graph Visualizer</h5>
      </div>
      <form @submit.prevent="gotoVisualizer">
        <h4>Subreddit:</h4>
        <input
          v-model="subreddit"
          v-tooltip="{
            trigger: 'manual',
            show: tooltip.show,
            content: tooltip.content,
          }"
          @click="hideTooltip"
          type="text"
          placeholder="uwaterloo"
          required
        />
        <button type="submit">visualize</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  data: () => ({
    subreddit: "",
    tooltip: { show: false },
  }),
  methods: {
    gotoVisualizer() {
      let { subreddit: sr } = this;

      // Validate subreddit.
      if (!sr) {
        this.tooltip.content = "Subreddit cannot be blank.";
        this.tooltip.show = true;
        return;
      }

      const rsi = sr.indexOf("r/");
      if (rsi > 0) sr = sr.slice(rsi + 2);
      if (!sr.match("^\\w+$")) {
        this.tooltip.content = "Invalid subreddit name.";
        this.tooltip.show = true;
        return;
      }

      this.$router.push({
        name: "visualizer",
        params: { subreddit: sr },
      });
    },
    hideTooltip() {
      this.tooltip.show = false;
    },
  },
};
</script>

<style lang="scss" scoped>
$gradient: linear-gradient(to top right, #b535f6, #376bf6);

// prettier-ignore
.home {
  top: 0; left: 0; right: 0; bottom: 0;
  position: absolute;
  justify-content: center;

  background-image: $gradient;
}

.card {
  $card-text-color: rgb(55, 55, 55);

  padding: 20px;
  border-radius: 10px;

  color: $card-text-color;
  background-color: white;
  box-shadow: 3px 3px 10px 2px rgba(black, 0.2);

  .logo {
    max-width: 200px;
    width: 75vw;
  }

  .title {
    margin: 5px 0 20px 0;

    .description {
      font-weight: 600;
      font-size: 0.9em;
      color: lighten($color: $card-text-color, $amount: 25%);
    }
  }

  form {
    width: 200px;
    display: flex;
    flex-direction: column;

    // prettier-ignore
    h4 { font-weight: 600; }

    input {
      $bg-color: rgb(227, 227, 227);

      width: 100%;
      margin: 2px 0;
      padding: 4px 6px;
      box-sizing: border-box;

      outline: none;
      border: none;
      border-radius: 4px;

      color: rgb(50, 50, 50);
      background-color: $bg-color;
      font-size: 11pt;

      transition: background 75ms ease-in-out;

      // prettier-ignore
      &::placeholder { color: rgb(160, 160, 160); }

      // prettier-ignore
      &:hover, &:focus, &:valid {
        background-color: lighten($color: $bg-color, $amount: 3%);
      }
    }

    button {
      $bg-color: #b45efa;

      align-self: center;
      margin-top: 10px;
      padding: 3px 7px;

      outline: none;
      border: none;
      border-radius: 4px;

      color: white;
      background-color: $bg-color;
      cursor: pointer;

      font-size: 10pt;
      font-weight: 500;

      transition: background 75ms ease-in-out;

      &:hover {
        background-color: lighten($color: $bg-color, $amount: 3%);
      }
    }
  }
}
</style>
