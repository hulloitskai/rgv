<template>
  <div id="app"><router-view /></div>
</template>

<style lang="scss">
body {
  // Normalize spacing in body and fonts.
  padding: 0;
  margin: 0;

  // prettier-ignore
  p, h1, h2, h3, h4, h5, h6 {
    margin: 0;
    padding: 0;
  }

  // Configure smoothing.
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
