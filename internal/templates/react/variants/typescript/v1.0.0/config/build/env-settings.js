const buildSettings = {
  development: {
    optimization: {
      minimize: false,
      splitChunks: false,
    },
    performance: {
      hints: false,
    },
  },
  staging: {
    optimization: {
      minimize: true,
      splitChunks: {
        chunks: 'all',
      },
    },
    performance: {
      hints: 'warning',
    },
  },
  production: {
    optimization: {
      minimize: true,
      splitChunks: {
        chunks: 'all',
        minSize: 20000,
        maxSize: 244000,
      },
    },
    performance: {
      hints: 'error',
      maxEntrypointSize: 512000,
      maxAssetSize: 512000,
    },
  },
};

module.exports = (env) => buildSettings[env] || buildSettings.production;