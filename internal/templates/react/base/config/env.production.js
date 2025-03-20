module.exports = {
  NODE_ENV: 'production',
  PUBLIC_URL: '',
  API_URL: 'https://api.production.com',
  ENABLE_SOURCEMAP: false,
  BUNDLE_ANALYZER: true,
  COMPRESSION_ENABLED: true,
  ASSET_PREFIX: 'https://cdn.production.com',
  SENTRY_DSN: process.env.SENTRY_DSN,
  PERFORMANCE_BUDGET: {
    maxJs: 512000,
    maxCss: 100000,
    maxMedia: 1000000,
  },
};