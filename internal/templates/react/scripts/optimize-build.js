const webpack = require('webpack');
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin;
const SpeedMeasurePlugin = require('speed-measure-webpack-plugin');

class BuildOptimizer {
  constructor(config) {
    this.config = config;
    this.smp = new SpeedMeasurePlugin();
  }

  async analyze() {
    const config = {
      ...this.config,
      plugins: [
        ...this.config.plugins,
        new BundleAnalyzerPlugin({
          analyzerMode: 'static',
          reportFilename: 'bundle-analysis.html',
        }),
      ],
    };

    return new Promise((resolve, reject) => {
      webpack(this.smp.wrap(config), (err, stats) => {
        if (err || stats.hasErrors()) reject(err || stats.toString());
        resolve(stats);
      });
    });
  }
}

module.exports = BuildOptimizer;