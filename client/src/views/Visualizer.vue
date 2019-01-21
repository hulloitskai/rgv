<template>
  <div class="visualizer">
    <div class="graph-container" ref="graph-container" />
    <div class="status flex center" :class="{ live: isLive }">
      <div class="indicator" />
      <p>live</p>
    </div>
  </div>
</template>

<script>
import Streamer from "@/services/streamer";
import * as Drawing from "jsnetworkx/node/drawing";
import * as d3 from "d3"; // cannot be modularized due to event mechanism

// TODO: Add zoom-in zoom-out indicator / button.
export default {
  data() {
    const { subreddit } = this.$route.params;
    return {
      isLive: false,
      isMounted: false,
      streamer: new Streamer(),
      graph: undefined,
      subreddit,
    };
  },
  methods: {
    /** @param {Event} event */
    handleMessage(event) {
      const data = JSON.parse(event.data);
      // console.log(data);

      // Early return if message is an error.
      const { error } = data;
      if (error) {
        alert(`Received error: ${error}`);
        return;
      }

      // Message is a data point, add it to the graph.
      const { graph } = this;
      const { author, link_author: linkAuthor } = data;

      // Add author and linkAuthor to graph.
      [author, linkAuthor].forEach(user => {
        if (!user) return;

        // Determine new node weight.
        let weight = 1;
        const node = graph.node.get(user);
        if (node) weight = node.weight + 1;

        graph.addNode(user, { weight });
      });

      // Create edges between link authors.
      if (!linkAuthor) return;

      const edgeData = graph.getEdgeData(author, linkAuthor, { weight: 0 });
      graph.addEdge(author, linkAuthor, { weight: edgeData.weight + 1 });
    },
    // TODO: Handle bot-in-creation status (yellow).
    handleStatus(event) {
      this.isLive = event.type === "open";
    },
    initDrawing() {
      // Only start drawing if both 'isMounted' and 'graph !== undefined'.
      const { isMounted, graph } = this;
      if (!(isMounted && graph)) return;

      // Initialize drawing.
      const { "graph-container": container } = this.$refs;
      Drawing.draw(
        this.graph,
        {
          withLabels: true,
          weighted: true,
          layoutAttr: {
            gravity: 0.035,
            linkStrength: 0.4,
          },
          nodeAttr: {
            r: ({ data }) => data.weight * 5,
            class: "node node-shape",
          },
          labelAttr: {
            class: "node-label",
            transform: ({ data }) => `translate(0, ${data.weight * 5 + 10})`,
          },
          edgeAttr: { class: "line edge" },
          element: container,
          d3,
        },
        true
      );

      // Correct svg boundaries.
      const svg = container.firstChild;
      ["width", "height"].forEach(name => svg.setAttribute(name, "100%"));

      // Start streamer.
      const { streamer, subreddit } = this;
      streamer.load(subreddit);
      streamer.addEventListener("message", this.handleMessage);
      ["open", "close"].forEach(type =>
        streamer.addEventListener(type, this.handleStatus)
      );
    },
  },
  async created() {
    const {
      default: DiGraph,
    } = await import("jsnetworkx/node/classes/DiGraph");
    this.graph = new DiGraph();
    this.initDrawing();
  },
  mounted() {
    this.isMounted = true;
    this.initDrawing();
  },
  beforeDestroy() {
    const { streamer } = this;
    streamer.removeEventListener("message", this.handleMessage);
    ["open", "close"].forEach(type =>
      streamer.removeEventListener(type, this.handleStatus)
    );
  },
  metaInfo: {
    title: "RGV",
    titleTemplate(name) {
      const { subreddit } = this.$route.params;
      return `${name}: ${subreddit}`;
    },
  },
};
</script>

<style lang="scss" scoped>
@import "@/styles/mixins.scss";

// prettier-ignore
.visualizer {
  position: absolute;
  left: 0; right: 0; top: 0; bottom: 0;
}

.graph-container {
  width: 100%;
  height: 100%;

  /deep/ .node {
    fill: #ff8300 !important;
    stroke-width: 0 !important;
  }

  /deep/ .node-label {
    fill: #666666 !important;
    font-size: 11pt;
  }

  /deep/ .edge {
    fill: #ffc180 !important;
  }
}

.status {
  $offset: 0.9em;
  $phablet-offset: 1.25em;

  position: absolute;
  top: $offset;
  left: $offset;
  padding: 0.2em 0.5em;
  border-radius: 1em;

  background-color: rgb(240, 240, 240);
  color: rgb(180, 180, 180);
  font-size: 11pt;
  font-weight: 500;

  .indicator {
    $radius: 0.75em;

    width: $radius;
    height: $radius;
    margin-right: 0.3em;
    border-radius: 100%;

    background-color: rgb(243, 139, 139);
    transition: background 200ms ease-in-out;
  }

  &.live {
    color: grey;

    // prettier-ignore
    .indicator { background-color: rgb(71, 224, 148); }
  }

  @include breakpoint(phablet) {
    top: $phablet-offset;
    left: $phablet-offset;
  }
}
</style>
