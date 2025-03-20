export const getImageUrl = (path: string): string => {
  const image = require(`../assets/images/${path}`);
  return image.src;
};

export const getResponsiveImage = (path: string) => {
  const image = require(`../assets/images/${path}?sizes[]=300,sizes[]=600,sizes[]=1200,sizes[]=2000`);
  return {
    src: image.src,
    srcSet: image.srcSet,
    width: image.width,
    height: image.height,
    placeholder: image.placeholder,
  };
};

export const preloadImage = (path: string): void => {
  const link = document.createElement('link');
  link.rel = 'preload';
  link.as = 'image';
  link.href = getImageUrl(path);
  document.head.appendChild(link);
};