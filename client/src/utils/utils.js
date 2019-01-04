const getOpt = key => process.env[`VUE_APP_${key}`];

/** @returns {string} */
function getAPIURL() {
  const apiURL = getOpt("API_URL");
  if (apiURL) return apiURL;

  const { hostname } = window.location;
  return `${hostname}/api`;
}

export { getOpt, getAPIURL };
