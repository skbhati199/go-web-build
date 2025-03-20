const path = require('path');

module.exports = {
  images: {
    test: /\.(png|jpg|jpeg|gif|webp|avif)$/i,
    type: 'asset',
    parser: {
      dataUrlCondition: {
        maxSize: 8 * 1024, // 8kb
      },
    },
    generator: {
      filename: 'static/media/[name].[hash:8][ext]',
    },
  },
  
  svgs: {
    test: /\.svg$/,
    use: ['@svgr/webpack', 'url-loader'],
    generator: {
      filename: 'static/media/[name].[hash:8][ext]',
    },
  },

  fonts: {
    test: /\.(woff|woff2|eot|ttf|otf)$/i,
    type: 'asset/resource',
    generator: {
      filename: 'static/fonts/[name].[hash:8][ext]',
    },
  },
};