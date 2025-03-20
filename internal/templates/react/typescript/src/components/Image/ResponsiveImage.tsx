import React from 'react';

interface ImageProps {
  src: string;
  alt: string;
  sizes?: string;
  className?: string;
  loading?: 'lazy' | 'eager';
}

export const ResponsiveImage: React.FC<ImageProps> = ({
  src,
  alt,
  sizes = '100vw',
  className,
  loading = 'lazy',
}) => {
  const image = require(`../../assets/images/${src}?sizes[]=300,sizes[]=600,sizes[]=1200,sizes[]=2000`);

  return (
    <img
      src={image.src}
      srcSet={image.srcSet}
      sizes={sizes}
      alt={alt}
      className={className}
      loading={loading}
      width={image.width}
      height={image.height}
    />
  );
};