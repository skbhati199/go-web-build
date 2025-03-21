# Serverless Deployment

The go-web-build framework provides built-in support for serverless deployments, allowing you to easily deploy your applications to serverless platforms like AWS Lambda.

## Overview

The serverless module in go-web-build enables you to:

- Deploy functions to AWS Lambda
- Create API Gateway endpoints for your functions
- Manage function configurations and environment variables
- Monitor the status of your deployed serverless resources

## Getting Started

### Configuration

To use the serverless features, you need to configure your serverless deployment in your `gobuild.yaml` file:

```yaml
serverless:
  provider:
    type: aws
    region: us-west-2
    credentials:
      accessKey: ${AWS_ACCESS_KEY_ID}
      secretKey: ${AWS_SECRET_ACCESS_KEY}
  
  functions:
    - name: my-function
      runtime: nodejs16.x
      memory: 256
      timeout: 30
      handler: index.handler
      codePath: ./functions/my-function
      environment:
        NODE_ENV: production
        DB_HOST: ${DB_HOST}
      template: node-lambda
    
    - name: api-function
      runtime: go1.x
      memory: 128
      timeout: 10
      handler: main
      codePath: ./functions/api-function
      environment:
        STAGE: production
  
  triggers:
    - type: http
      function: api-function
      properties:
        path: users
        method: GET