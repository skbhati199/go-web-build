type ImportCallback<T> = () => Promise<{ default: T }>;

export const lazyLoad = <T>(importCallback: ImportCallback<T>) => {
  return new Promise<T>((resolve) => {
    importCallback().then((module) => {
      resolve(module.default);
    });
  });
};

export const withPreload = <T>(importCallback: ImportCallback<T>) => {
  const Component = React.lazy(importCallback);
  Component.preload = importCallback;
  return Component;
};

export const preloadComponents = (components: Array<{ preload?: () => Promise<any> }>) => {
  return Promise.all(components.map((component) => component.preload?.()));
};