import { ModuleFederationConfig } from 'webpack/lib/container/ModuleFederationPlugin';

interface RemoteConfig {
  name: string;
  url: string;
}

export const createFederationConfig = (remotes: RemoteConfig[]): ModuleFederationConfig => ({
  name: 'host',
  filename: 'remoteEntry.js',
  remotes: remotes.reduce((acc, remote) => ({
    ...acc,
    [remote.name]: `${remote.name}@${remote.url}/remoteEntry.js`,
  }), {}),
  exposes: {
    './App': './src/App',
    './components/Button': './src/components/Button',
    './components/Card': './src/components/Card',
  },
  shared: {
    react: { singleton: true, requiredVersion: '^18.0.0' },
    'react-dom': { singleton: true, requiredVersion: '^18.0.0' },
    '@emotion/react': { singleton: true },
    '@emotion/styled': { singleton: true },
  },
});