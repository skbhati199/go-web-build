import webpack from 'webpack';
import WebpackDevServer from 'webpack-dev-server';
import config from '../config/webpack.dev';

const compiler = webpack(config);
const devServerOptions = {
  ...config.devServer,
  open: true,
};

const server = new WebpackDevServer(devServerOptions, compiler);

const runServer = async () => {
  console.log('Starting development server...');
  await server.start();
};

runServer().catch((err) => {
  console.error('Failed to start development server:', err);
  process.exit(1);
});