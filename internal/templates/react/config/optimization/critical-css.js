const CriticalCssPlugin = require('critical-css-webpack-plugin');

module.exports = {
  criticalCss: {
    plugin: new CriticalCssPlugin({
      base: 'dist/',
      src: 'index.html',
      target: 'index.html',
      inline: true,
      extract: true,
      width: 1200,
      height: 900,
      penthouse: {
        blockJSRequests: false,
      },
    }),
  },
};