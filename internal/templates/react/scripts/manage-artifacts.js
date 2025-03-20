const fs = require('fs').promises;
const path = require('path');

class ArtifactManager {
  constructor(artifactsDir = 'artifacts') {
    this.artifactsDir = artifactsDir;
  }

  async archiveBuild(version) {
    const buildDir = path.resolve(process.cwd(), 'build');
    const artifactPath = path.join(
      this.artifactsDir,
      `build-${version}-${Date.now()}`
    );

    await fs.mkdir(artifactPath, { recursive: true });
    await fs.cp(buildDir, artifactPath, { recursive: true });
    
    return artifactPath;
  }

  async cleanOldArtifacts(keepCount = 5) {
    const files = await fs.readdir(this.artifactsDir);
    const artifacts = files
      .filter(f => f.startsWith('build-'))
      .sort()
      .reverse();

    if (artifacts.length > keepCount) {
      const toRemove = artifacts.slice(keepCount);
      await Promise.all(
        toRemove.map(artifact =>
          fs.rm(path.join(this.artifactsDir, artifact), { recursive: true })
        )
      );
    }
  }
}

module.exports = ArtifactManager;