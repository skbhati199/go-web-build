module.exports = {
  plugins: {
    'postcss-import': {},
    'postcss-preset-env': {
      stage: 1,
      features: {
        'nesting-rules': true,
        'custom-properties': true,
        'custom-media-queries': true,
      },
    },
    'postcss-custom-media': {},
    'postcss-nested': {},
    'postcss-flexbugs-fixes': {},
    autoprefixer: {
      flexbox: 'no-2009',
      grid: 'autoplace',
    },
    cssnano: process.env.NODE_ENV === 'production' ? {
      preset: ['default', {
        discardComments: { removeAll: true },
        normalizeWhitespace: false,
      }],
    } : false,
  },
};