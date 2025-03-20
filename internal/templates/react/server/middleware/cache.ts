import { Request, Response, NextFunction } from 'express';

interface CacheOptions {
  maxAge?: number;
  private?: boolean;
  immutable?: boolean;
}

export const cacheControl = (options: CacheOptions = {}) => {
  return (req: Request, res: Response, next: NextFunction) => {
    const directives = ['public'];
    
    if (options.maxAge) {
      directives.push(`max-age=${options.maxAge}`);
    }
    
    if (options.private) {
      directives.push('private');
    }
    
    if (options.immutable) {
      directives.push('immutable');
    }
    
    res.setHeader('Cache-Control', directives.join(', '));
    next();
  };
};

export const staticAssetCache = cacheControl({
  maxAge: 31536000,
  immutable: true,
});

export const dynamicContentCache = cacheControl({
  maxAge: 300,
  private: true,
});