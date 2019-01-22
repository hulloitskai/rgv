<template>
  <div class="not-found flex col center"><code>404 page not found</code></div>
</template>

<style lang="scss" scoped>
@import url("https://fonts.googleapis.com/css?family=Inconsolata:700");
$gradient: linear-gradient(to bottom right, #7b4397, #dc2430);

// prettier-ignore
.not-found {
  top: 0; left: 0; right: 0; bottom: 0;
  position: absolute;
  justify-content: center;
  background-image: $gradient;
}

code {
  font-family: Inconsolata, "Courier New", Courier, monospace;
  font-weight: bold;
  font-size: 24pt;
  color: white;
}
</style>
