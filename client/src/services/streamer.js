import Logger from "@/utils/logger";
import { getOpt } from "@/utils/utils";

/**
 * A Streamer is able to stream Reddit activity from the rgv API server.
 */
class Streamer {
  /** @type WebSocket */
  ws = null;
  /** @type Logger */
  l = undefined;

  constructor(logger = new Logger("Streamer")) {
    this.l = logger;
  }

  /**
   * Open a connection to the API socket server, and request for a particular
   * subreddit.
   * @param {string} subreddit
   */
  load(subreddit) {
    // Validate arguments.
    if (!subreddit) {
      this.l.error("Subreddit must be a non-empty string.");
      return;
    }

    // Create and connect websocket, configure API server.
    const API_URL = getOpt("API_URL");
    if (!API_URL) {
      this.l.error("Env var API_URL was not set.");
      return;
    }

    const ws = new WebSocket(API_URL);
    ws.addEventListener("open", event => {
      this.l.debug("Websocket opened:", event);

      const config = { subreddit };
      ws.send(JSON.stringify(config));
    });
    this.ws = ws;
    this.configureWS();
  }

  configureWS() {
    const { ws, l } = this;
    ws.addEventListener("error", event => l.error("Websocket error:", event));
    ws.addEventListener("close", event => l.debug("Websocket closed:", event));
  }

  /**
   * @param {string} type
   * @param {(this: WebSocket, ev: Event) => any} listener
   */
  addEventListener(type, listener) {
    this.ws.addEventListener(type, listener);
  }

  /**
   * @param {string} type
   * @param {(this: WebSocket, ev: Event) => any} listener
   */
  removeEventListener(type, listener) {
    this.ws.removeEventListener(type, listener);
  }
}

export default Streamer;
