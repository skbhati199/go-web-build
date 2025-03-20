const path = require('path');
const ReactRefreshWebpackPlugin = require('@pmmmwh/react-refresh-webpack-plugin');
const { merge } = require('webpack-merge');
const baseConfig = require('./webpack.common');

module.exports = merge(baseConfig, {
  mode: 'development',
  devtool: 'eval-source-map',
  
  devServer: {
    hot: true,
    port: 3000,
    historyApiFallback: true,
    client: {
      overlay: {
        errors: true,
        warnings: false,
      },
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
    },
  },

  plugins: [
    new ReactRefreshWebpackPlugin(),
  ],

  optimization: {
    minimize: false,
    splitChunks: {
      chunks: 'all',
    },
  },
});