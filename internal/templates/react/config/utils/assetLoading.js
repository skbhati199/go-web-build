const path = require('path');

const createLoadingStrategy = (env) => ({
  images: {
    inlineLimit: env.production ? 4096 : 10000,
    formats: ['webp', 'avif'],
    minimumCacheTTL: 60,
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
    domains: [],
    path: '/_next/image',
    loader: 'default',
    disableStaticImages: false
  },
  
  fonts: {
    preload: true,
    preloadFonts: [
      '/fonts/inter-var-latin.woff2'
    ],
    displaySwap: true,
    inlineLimit: 0
  },
  
  static: {
    directory: path.join(process.cwd(), 'public'),
    publicPath: '/',
    serveIndex: false,
    watch: true
  }
});

module.exports = {
  createLoadingStrategy
};