const path = require('path');

module.exports = {
  cdn: {
    enabled: process.env.USE_CDN === 'true',
    provider: process.env.CDN_PROVIDER || 'cloudfront',
    domain: process.env.CDN_DOMAIN,
    bucket: process.env.CDN_BUCKET,
    paths: {
      images: '/static/media',
      js: '/static/js',
      css: '/static/css',
      fonts: '/static/fonts',
    },
    options: {
      maxAge: 31536000,
      immutable: true,
    },
  },
};