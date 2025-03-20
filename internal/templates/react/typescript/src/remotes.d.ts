declare module 'remote/Button';
declare module 'remote/Modal';

export interface RemoteConfig {
  url: string;
  scope: string;
  module: string;
}

export const remotes: Record<string, RemoteConfig> = {
  button: {
    url: 'http://localhost:3001/remoteEntry.js',
    scope: 'remote',
    module: './Button',
  },
  modal: {
    url: 'http://localhost:3001/remoteEntry.js',
    scope: 'remote',
    module: './Modal',
  },
};