const path = require('path');
const fs = require('fs');

async function validateBuild() {
  const buildPath = path.resolve(process.cwd(), 'build');
  const requiredFiles = ['index.html', 'asset-manifest.json'];
  
  for (const file of requiredFiles) {
    if (!fs.existsSync(path.join(buildPath, file))) {
      throw new Error(`Required file ${file} is missing from build`);
    }
  }

  // Validate bundle sizes
  const stats = require(path.join(buildPath, 'bundle-stats.json'));
  const maxSize = 1024 * 1024; // 1MB

  for (const asset of stats.assets) {
    if (asset.size > maxSize) {
      console.warn(`Warning: ${asset.name} exceeds recommended size`);
    }
  }
}

module.exports = validateBuild;