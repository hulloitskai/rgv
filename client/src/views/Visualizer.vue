<template>
  <div class="visualizer">
    <h1 class="title">messages</h1>
    <ul class="messages">
      <li class="message" v-for="msg in messages" :key="msg.id">
        {{ JSON.stringify(msg) }}
      </li>
    </ul>
  </div>
</template>

<script>
const API_URL = "ws://localhost:3000";

export default {
  data() {
    const ws = new WebSocket(API_URL);
    ws.onopen = event => {
      console.log(`Websocket opened: ${event}`);

      const { subreddit } = this.$route.params;
      if (!subreddit) {
        console.error("No subreddit found in boute params.");
        return;
      }

      const config = { subreddit: this.$route.params.subreddit };
      this.ws.send(JSON.stringify(config));
    };
    ws.onclose = event => console.log(`Websocket closed: ${event}`);
    ws.onerror = event => console.error(`Websocket emitted an error: ${event}`);
    ws.onmessage = this.appendMessage;

    return { messages: [], ws };
  },
  methods: {
    appendMessage(msg) {
      this.messages.push(JSON.parse(msg.data));
    },
  },
};
</script>
