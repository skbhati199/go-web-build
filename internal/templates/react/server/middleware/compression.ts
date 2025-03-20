import compression from 'compression';
import { Request, Response } from 'express';

export const compressionMiddleware = compression({
  level: 6,
  threshold: 1024,
  filter: (req: Request, res: Response) => {
    if (req.headers['x-no-compression']) {
      return false;
    }
    
    return compression.filter(req, res);
  },
});

export const staticCompression = {
  gzip: true,
  brotli: true,
  extensions: ['.js', '.css', '.html', '.svg', '.json'],
  threshold: '1kb',
};