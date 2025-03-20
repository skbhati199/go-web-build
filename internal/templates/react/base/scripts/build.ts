import webpack from 'webpack';
import prodConfig from '../config/webpack.prod';
import { measureBuildPerformance } from './performance';

const build = async () => {
  console.log('Creating production build...');
  const startTime = Date.now();

  try {
    const stats = await new Promise<webpack.Stats>((resolve, reject) => {
      webpack(prodConfig, (err, stats) => {
        if (err || stats?.hasErrors()) {
          reject(err || stats?.compilation.errors);
        }
        resolve(stats!);
      });
    });

    const endTime = Date.now();
    console.log(`Build completed in ${(endTime - startTime) / 1000}s`);
    
    measureBuildPerformance(stats);
  } catch (error) {
    console.error('Build failed:', error);
    process.exit(1);
  }
};

build();