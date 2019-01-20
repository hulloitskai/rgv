/**
 * @param {string} key
 * @returns {string}
 */
const getOpt = key => process.env[`VUE_APP_${key}`];

// Set API_URL.
let API_URL = getOpt("API_URL");
if (!API_URL) {
  const { hostname } = window.location;
  API_URL = `${hostname}/api`;
}

export { getOpt, API_URL };
