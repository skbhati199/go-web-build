/**
 * Development server configuration
 */
module.exports = {
  port: 3000,
  hot: true,
  historyApiFallback: true,
  compress: true,
  client: {
    overlay: {
      errors: true,
      warnings: false,
    },
    progress: true,
  },
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
  },
  static: {
    directory: './public',
    publicPath: '/',
  },
  devMiddleware: {
    publicPath: '/',
  },
};