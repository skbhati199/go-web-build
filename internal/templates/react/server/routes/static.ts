import express from 'express';
import { staticAssetCache, dynamicContentCache } from '../middleware/cache';
import { compressionMiddleware } from '../middleware/compression';
import { securityHeaders } from '../middleware/security';

const router = express.Router();

router.use('/static',
  securityHeaders,
  compressionMiddleware,
  staticAssetCache,
  express.static('build/static', {
    index: false,
    etag: true,
    lastModified: true,
  })
);

router.use('/',
  securityHeaders,
  compressionMiddleware,
  dynamicContentCache,
  express.static('build', {
    index: false,
    etag: true,
  })
);

export default router;