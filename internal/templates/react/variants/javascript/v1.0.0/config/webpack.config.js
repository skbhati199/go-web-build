const path = require('path');
const { merge } = require('webpack-merge');
const baseConfig = require('./webpack.base');
const optimizationConfig = require('./webpack.optimization');
const assetsConfig = require('./webpack.assets');

module.exports = (env) => {
  const isProduction = env.production;
  
  return merge(
    baseConfig,
    optimizationConfig,
    assetsConfig,
    {
      mode: isProduction ? 'production' : 'development',
      devtool: isProduction ? 'source-map' : 'eval-source-map',
    }
  );
};