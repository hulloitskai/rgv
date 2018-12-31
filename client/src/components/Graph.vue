<template>
  <div class="graph"><canvas ref="canvas" /></div>
</template>

<script>
import { Renderer } from "@/utils/arbor";
const { arbor } = window;

export default {
  data() {
    return {
      system: null,
    };
  },
  mounted() {
    const sys = arbor.ParticleSystem(1000, 600, 0.5);
    sys.parameters({ gravity: true });
    sys.renderer = new Renderer(this.$refs.canvas);
    this.system = sys;

    // TEMPORARY:
    sys.addEdge("a", "b");
    sys.addEdge("a", "c");
    sys.addEdge("a", "d");
    sys.addEdge("a", "e");
    sys.addNode("f", { alone: true, mass: 0.25 });
  },
};
</script>

<style lang="scss" scoped>
.graph {
  flex: 1;
  align-self: stretch;
}

canvas {
  width: 100%;
  height: 100%;
}
</style>
