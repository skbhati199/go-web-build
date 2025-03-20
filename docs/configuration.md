# Configuration Reference

## Overview

Go Web Build uses a YAML-based configuration system that supports multiple environments (development, production, staging). This guide covers all available configuration options and their usage.

## Configuration File

The default configuration file is `gobuild.yaml` in your project root. Environment-specific configurations can be placed in the `config/` directory.

## Basic Structure

```yaml
environment: development # or production, staging
```

## Project Settings

### Template Configuration
```yaml
template:
  version: "1.0"
  variables:
    # Custom template variables
    projectName: "my-app"
    author: "Your Name"
```

### Development Server
```yaml
dev:
  port: 3000
  host: "localhost"
  hot_reload: true
  proxy:
    "/api": "http://localhost:8080"
```

### Build Configuration
```yaml
build:
  output_dir: "dist"
  source_maps: true
  optimization:
    minify: true
    split_chunks: true
    tree_shaking: true
  assets:
    images:
      optimize: true
      responsive: true
    styles:
      postcss: true
      css_modules: true
```

### Source Maps
```yaml
source_maps:
  development:
    type: "inline"
    include_content: true
  production:
    type: "external"
    url_prefix: "https://cdn.example.com/sourcemaps/"
```

### Error Tracking
```yaml
error_tracking:
  sentry:
    dsn: "your-sentry-dsn"
    upload_source_maps: true
  rollbar:
    access_token: "your-rollbar-token"
```

### Docker Configuration
```yaml
docker:
  base_image: "node:16-alpine"
  expose_ports:
    - "3000:3000"
  volumes:
    - ".:/app"
  health_check:
    path: "/health"
    interval: "30s"
```

### Static File Serving
```yaml
static:
  dir: "public"
  cdn:
    enabled: true
    url: "https://cdn.example.com"
  cache_control: "public, max-age=31536000"
  compression: true
```

### Testing Configuration
```yaml
test:
  runner: "jest"
  coverage:
    enabled: true
    threshold: 80
  e2e:
    browser: "chrome"
    headless: true
```

## Environment Variables

Configuration values can be overridden using environment variables. Use the prefix `GOBUILD_` followed by the uppercase configuration path with underscores.

Example:
```bash
GOBUILD_DEV_PORT=4000
GOBUILD_BUILD_OPTIMIZATION_MINIFY=false
```

## Example Configuration

Here's a complete example configuration file:

```yaml
environment: development

template:
  version: "1.0"
  variables:
    projectName: "my-web-app"

dev:
  port: 3000
  hot_reload: true

build:
  output_dir: "dist"
  source_maps: true
  optimization:
    minify: true
    split_chunks: true

static:
  dir: "public"
  compression: true

test:
  runner: "jest"
  coverage:
    enabled: true
```