import { RemoteConfig } from '../remotes';

export async function loadRemoteComponent(config: RemoteConfig) {
  await loadRemoteEntry(config.url);
  const container = window[config.scope];
  await container.init(__webpack_share_scopes__.default);
  const factory = await container.get(config.module);
  return factory();
}

async function loadRemoteEntry(url: string) {
  return new Promise<void>((resolve, reject) => {
    const element = document.createElement('script');
    element.src = url;
    element.type = 'text/javascript';
    element.async = true;
    element.onload = () => resolve();
    element.onerror = reject;
    document.head.appendChild(element);
  });
}