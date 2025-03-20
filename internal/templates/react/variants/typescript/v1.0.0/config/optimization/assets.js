const CompressionPlugin = require('compression-webpack-plugin');
const zlib = require('zlib');

module.exports = {
  optimization: {
    splitChunks: {
      cacheGroups: {
        images: {
          test: /\.(png|jpe?g|gif|svg|webp)$/,
          chunks: 'all',
          enforce: true,
          name: 'images'
        },
        fonts: {
          test: /\.(woff2?|eot|ttf|otf)$/,
          chunks: 'all',
          enforce: true,
          name: 'fonts'
        }
      }
    }
  },
  plugins: [
    new CompressionPlugin({
      filename: '[path][base].gz',
      algorithm: 'gzip',
      test: /\.(js|css|html|svg)$/,
      threshold: 10240,
      minRatio: 0.8
    }),
    new CompressionPlugin({
      filename: '[path][base].br',
      algorithm: 'brotliCompress',
      test: /\.(js|css|html|svg)$/,
      compressionOptions: {
        params: {
          [zlib.constants.BROTLI_PARAM_QUALITY]: 11
        }
      },
      threshold: 10240,
      minRatio: 0.8
    })
  ]
};