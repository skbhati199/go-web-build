import { Stats } from 'webpack';
import filesize from 'filesize';
import chalk from 'chalk';

export const measureBuildPerformance = (stats: Stats) => {
  const assets = stats.toJson().assets || [];
  const totalSize = assets.reduce((size, asset) => size + asset.size, 0);
  
  console.log('\nFile sizes after build:\n');
  
  assets
    .filter(asset => !asset.name.endsWith('.map'))
    .sort((a, b) => b.size - a.size)
    .forEach(asset => {
      console.log(
        `${chalk.cyan(asset.name.padEnd(40))} ${chalk.green(filesize(asset.size))}`
      );
    });

  console.log(`\nTotal size: ${chalk.green(filesize(totalSize))}`);
};