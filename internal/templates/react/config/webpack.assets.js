const ImageMinimizerPlugin = require('image-minimizer-webpack-plugin');
const { extendDefaultPlugins } = require('svgo');

module.exports = {
  module: {
    rules: [
      {
        test: /\.(png|jpg|gif|webp)$/i,
        type: 'asset',
        parser: {
          dataUrlCondition: {
            maxSize: 8 * 1024,
          },
        },
      },
      {
        test: /\.svg$/,
        use: ['@svgr/webpack', 'url-loader'],
      },
    ],
  },
  plugins: [
    new ImageMinimizerPlugin({
      minimizerOptions: {
        plugins: [
          ['gifsicle', { interlaced: true }],
          ['mozjpeg', { quality: 80 }],
          ['optipng', { optimizationLevel: 5 }],
          ['svgo', {
            plugins: extendDefaultPlugins([
              { name: 'removeViewBox', active: false },
            ]),
          }],
        ],
      },
    }),
  ],
};