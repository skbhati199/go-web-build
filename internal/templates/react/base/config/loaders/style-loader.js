const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const postcssNormalize = require('postcss-normalize');

module.exports = {
  styleLoader: (isProduction) => ({
    test: /\.(scss|sass|css)$/,
    use: [
      isProduction ? MiniCssExtractPlugin.loader : 'style-loader',
      {
        loader: 'css-loader',
        options: {
          modules: {
            auto: true,
            localIdentName: isProduction
              ? '[hash:base64:8]'
              : '[name]__[local]--[hash:base64:5]',
          },
          importLoaders: 3,
          sourceMap: !isProduction,
        },
      },
      {
        loader: 'postcss-loader',
        options: {
          postcssOptions: {
            plugins: [
              'postcss-flexbugs-fixes',
              ['postcss-preset-env', { autoprefixer: { flexbox: 'no-2009' } }],
              postcssNormalize(),
            ],
          },
          sourceMap: !isProduction,
        },
      },
      {
        loader: 'sass-loader',
        options: {
          sourceMap: !isProduction,
        },
      },
    ],
  }),
};