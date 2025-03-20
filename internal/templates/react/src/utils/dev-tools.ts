export const enableDevTools = () => {
  if (process.env.NODE_ENV === 'development') {
    const whyDidYouRender = require('@welldone-software/why-did-you-render');
    whyDidYouRender(React, {
      trackAllPureComponents: true,
      logOnDifferentValues: true,
    });

    if (module.hot) {
      module.hot.accept();
    }
  }
};

export const setupDevErrorBoundary = () => {
  if (process.env.NODE_ENV === 'development') {
    const { ErrorBoundary } = require('react-error-boundary');
    return ErrorBoundary;
  }
  return null;
};