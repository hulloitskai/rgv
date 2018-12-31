const getOpt = key => process.env[`VUE_APP_${key}`];

export { getOpt };
