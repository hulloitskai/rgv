/**
 * A Logger is an environment-aware logger, which only logs debug messages when
 * in development.
 */
class Logger {
  /** @type {string} */
  namespace;
  /** @type {boolean} */
  isDebug;

  /**
   * @param {string} namespace
   */
  constructor(namespace = "") {
    const { NODE_ENV } = process.env;
    this.isDebug = NODE_ENV === "development";
    this.namespace = namespace;
  }

  ///////////////////////
  // Logging functions
  ///////////////////////

  debug(...objs) {
    if (!this.isDebug) return;
    console.log(...this.appendNamespace(objs));
  }
  log(...objs) {
    console.log(...this.appendNamespace(objs));
  }
  info(...objs) {
    console.info(...this.appendNamespace(objs));
  }
  error(...objs) {
    console.error(...this.appendNamespace(objs));
  }

  /**
   * Build a child logger with a sub-namespace that is appended onto this
   * logger's namespace (with a '.').
   * @param {string} namespace
   * @returns {Logger}
   */
  child(namespace = "") {
    const { namespace: ns } = this.namespace;
    if (!namespace) return new Logger(ns);
    if (ns) return new Logger(`${ns}.${namespace}`);
    return new Logger(namespace);
  }

  /**
   * @param {string[]} objs
   * @returns {string[]}
   */
  appendNamespace(objs = []) {
    const { namespace } = this;
    if (namespace) return [`(${namespace})`, ...objs];
    return objs;
  }
}

export default Logger;
