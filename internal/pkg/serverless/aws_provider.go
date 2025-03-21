package serverless

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	apigatewayTypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// AWSProvider implements the Provider interface for AWS
type AWSProvider struct {
	config           ProviderConfig
	lambdaClient     *lambda.Client
	s3Client         *s3.Client
	apiGWClient      *apigateway.Client
	cfnClient        *cloudformation.Client
	deploymentBucket string
}

// NewAWSProvider creates a new AWS provider
func NewAWSProvider() Provider {
	return &AWSProvider{}
}

// Configure configures the AWS provider
func (p *AWSProvider) Configure(ctx context.Context, providerConfig ProviderConfig) error {
	p.config = providerConfig

	// Create options for AWS config
	var opts []func(*config.LoadOptions) error

	// Set region
	opts = append(opts, config.WithRegion(providerConfig.Region))

	// Create AWS credentials if provided
	if accessKey, ok := providerConfig.Credentials["accessKey"]; ok {
		secretKey := providerConfig.Credentials["secretKey"]
		sessionToken := providerConfig.Credentials["sessionToken"]
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken),
		))
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return fmt.Errorf("unable to load AWS config: %w", err)
	}

	// Initialize AWS service clients
	p.lambdaClient = lambda.NewFromConfig(cfg)
	p.s3Client = s3.NewFromConfig(cfg)
	p.apiGWClient = apigateway.NewFromConfig(cfg)
	p.cfnClient = cloudformation.NewFromConfig(cfg)

	// Set deployment bucket
	p.deploymentBucket = fmt.Sprintf("serverless-deploy-%s", providerConfig.Region)

	return nil
}

// Deploy deploys serverless resources to AWS
func (p *AWSProvider) Deploy(ctx context.Context, resources []Resource) (*DeploymentResult, error) {
	// Ensure deployment bucket exists
	if err := p.ensureDeploymentBucket(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure deployment bucket: %w", err)
	}

	// Upload function code to S3
	if err := p.uploadFunctionCode(ctx, resources); err != nil {
		return nil, fmt.Errorf("failed to upload function code: %w", err)
	}

	// Deploy Lambda functions
	deployedResources, err := p.deployFunctions(ctx, resources)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy functions: %w", err)
	}

	// Deploy API Gateway if needed
	apiEndpoint, err := p.deployAPIGateway(ctx, resources)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy API Gateway: %w", err)
	}

	// Create deployment result
	result := &DeploymentResult{
		Resources: deployedResources,
		Endpoint:  apiEndpoint,
		Version:   time.Now().Format("20060102150405"),
		CreatedAt: time.Now(),
	}

	return result, nil
}

// Remove removes a deployed serverless application
func (p *AWSProvider) Remove(ctx context.Context, name string) error {
	// Delete Lambda functions
	if err := p.deleteFunctions(ctx, name); err != nil {
		return fmt.Errorf("failed to delete functions: %w", err)
	}

	// Delete API Gateway
	if err := p.deleteAPIGateway(ctx, name); err != nil {
		return fmt.Errorf("failed to delete API Gateway: %w", err)
	}

	return nil
}

// GetStatus gets the status of a deployed serverless application
func (p *AWSProvider) GetStatus(ctx context.Context, name string) (*DeploymentStatus, error) {
	status := &DeploymentStatus{
		State:     "UNKNOWN",
		Resources: []ResourceStatus{},
		UpdatedAt: time.Now(),
	}

	// Get Lambda function status
	functionStatus, err := p.getFunctionStatus(ctx, name)
	if err != nil {
		status.LastError = err.Error()
	} else {
		status.Resources = append(status.Resources, functionStatus...)
		status.State = "ACTIVE"
	}

	// Get API Gateway status
	apiStatus, err := p.getAPIGatewayStatus(ctx, name)
	if err != nil {
		status.LastError = err.Error()
	} else if apiStatus != nil {
		status.Resources = append(status.Resources, *apiStatus)
	}

	return status, nil
}

// Helper methods
func (p *AWSProvider) ensureDeploymentBucket(ctx context.Context) error {
	// Check if bucket exists
	_, err := p.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(p.deploymentBucket),
	})

	if err == nil {
		// Bucket exists
		return nil
	}

	// Create bucket
	_, err = p.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(p.deploymentBucket),
	})

	return err
}

func (p *AWSProvider) uploadFunctionCode(ctx context.Context, resources []Resource) error {
	for _, resource := range resources {
		if resource.Type == "Function" {
			codePath, ok := resource.Properties["codePath"].(string)
			if !ok {
				continue
			}

			// Create zip file
			zipPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s.zip", resource.Name))
			if err := createZipFromDir(codePath, zipPath); err != nil {
				return fmt.Errorf("failed to create zip for %s: %w", resource.Name, err)
			}

			// Upload to S3
			zipFile, err := os.Open(zipPath)
			if err != nil {
				return fmt.Errorf("failed to open zip file: %w", err)
			}
			defer zipFile.Close()

			key := fmt.Sprintf("functions/%s/%s.zip", resource.Name, time.Now().Format("20060102150405"))
			_, err = p.s3Client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(p.deploymentBucket),
				Key:    aws.String(key),
				Body:   zipFile,
			})
			if err != nil {
				return fmt.Errorf("failed to upload function code: %w", err)
			}

			// Update resource properties with S3 location
			resource.Properties["s3Bucket"] = p.deploymentBucket
			resource.Properties["s3Key"] = key
		}
	}
	return nil
}

func (p *AWSProvider) deployFunctions(ctx context.Context, resources []Resource) ([]Resource, error) {
	deployedResources := make([]Resource, 0)

	for _, resource := range resources {
		if resource.Type == "Function" {
			// Create or update Lambda function
			functionInput := &lambda.CreateFunctionInput{
				FunctionName: aws.String(resource.Name),
				Role:         aws.String(p.getIAMRole()),
				Code: &lambdaTypes.FunctionCode{
					S3Bucket: aws.String(resource.Properties["s3Bucket"].(string)),
					S3Key:    aws.String(resource.Properties["s3Key"].(string)),
				},
			}

			// Set runtime using the proper enum type
			if runtime, ok := resource.Properties["runtime"].(string); ok {
				functionInput.Runtime = p.getRuntimeEnum(runtime)
			} else {
				functionInput.Runtime = lambdaTypes.RuntimeNodejs16x
			}

			// Set handler
			if handler, ok := resource.Properties["handler"].(string); ok {
				functionInput.Handler = aws.String(handler)
			} else {
				functionInput.Handler = aws.String("index.handler")
			}

			// Set memory size
			if memory, ok := resource.Properties["memory"].(int); ok {
				functionInput.MemorySize = aws.Int32(int32(memory))
			} else {
				functionInput.MemorySize = aws.Int32(128)
			}

			// Set timeout
			if timeout, ok := resource.Properties["timeout"].(int); ok {
				functionInput.Timeout = aws.Int32(int32(timeout))
			} else {
				functionInput.Timeout = aws.Int32(30)
			}

			// Set environment variables
			if env, ok := resource.Properties["environment"].(map[string]string); ok && len(env) > 0 {
				variables := make(map[string]string)
				for k, v := range env {
					variables[k] = v
				}
				functionInput.Environment = &lambdaTypes.Environment{
					Variables: variables,
				}
			}

			// Try to get function first to see if it exists
			_, err := p.lambdaClient.GetFunction(ctx, &lambda.GetFunctionInput{
				FunctionName: aws.String(resource.Name),
			})

			if err != nil {
				// Function doesn't exist, create it
				_, err = p.lambdaClient.CreateFunction(ctx, functionInput)
				if err != nil {
					return nil, fmt.Errorf("failed to create function %s: %w", resource.Name, err)
				}
			} else {
				// Function exists, update it
				updateInput := &lambda.UpdateFunctionCodeInput{
					FunctionName: aws.String(resource.Name),
					S3Bucket:     aws.String(resource.Properties["s3Bucket"].(string)),
					S3Key:        aws.String(resource.Properties["s3Key"].(string)),
				}
				_, err = p.lambdaClient.UpdateFunctionCode(ctx, updateInput)
				if err != nil {
					return nil, fmt.Errorf("failed to update function %s: %w", resource.Name, err)
				}

				// Update configuration
				configInput := &lambda.UpdateFunctionConfigurationInput{
					FunctionName: aws.String(resource.Name),
					Handler:      functionInput.Handler,
					Runtime:      functionInput.Runtime,
					Timeout:      functionInput.Timeout,
					MemorySize:   functionInput.MemorySize,
				}
				if functionInput.Environment != nil {
					configInput.Environment = functionInput.Environment
				}
				_, err = p.lambdaClient.UpdateFunctionConfiguration(ctx, configInput)
				if err != nil {
					return nil, fmt.Errorf("failed to update function configuration %s: %w", resource.Name, err)
				}
			}

			// Add to deployed resources
			deployedResources = append(deployedResources, resource)
		}
	}

	return deployedResources, nil
}

// Add a helper function to convert string to Runtime enum
func (p *AWSProvider) getRuntimeEnum(runtime string) lambdaTypes.Runtime {
	switch runtime {
	case "nodejs":
		return lambdaTypes.RuntimeNodejs
	case "nodejs12.x":
		return lambdaTypes.RuntimeNodejs12x
	case "nodejs14.x":
		return lambdaTypes.RuntimeNodejs14x
	case "nodejs16.x":
		return lambdaTypes.RuntimeNodejs16x
	case "nodejs18.x":
		return lambdaTypes.RuntimeNodejs18x
	case "python3.7":
		return lambdaTypes.RuntimePython37
	case "python3.8":
		return lambdaTypes.RuntimePython38
	case "python3.9":
		return lambdaTypes.RuntimePython39
	case "python3.10":
		return lambdaTypes.RuntimePython310
	case "java8":
		return lambdaTypes.RuntimeJava8
	case "java11":
		return lambdaTypes.RuntimeJava11
	case "java17":
		return lambdaTypes.RuntimeJava17
	case "go1.x":
		return lambdaTypes.RuntimeGo1x
	case "dotnet6":
		return lambdaTypes.RuntimeDotnet6
	case "ruby2.7":
		return lambdaTypes.RuntimeRuby27
	default:
		return lambdaTypes.RuntimeNodejs16x // Default to Node.js 16.x
	}
}

func (p *AWSProvider) deployAPIGateway(ctx context.Context, resources []Resource) (string, error) {
	// Check if we have any HTTP triggers
	hasHttpTriggers := false
	for _, resource := range resources {
		if resource.Type == "HttpTrigger" {
			hasHttpTriggers = true
			break
		}
	}

	if !hasHttpTriggers {
		return "", nil
	}

	// Create API Gateway
	apiName := fmt.Sprintf("serverless-api-%s", time.Now().Format("20060102150405"))
	createResult, err := p.apiGWClient.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
		Name: aws.String(apiName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create API Gateway: %w", err)
	}

	apiId := *createResult.Id

	// Get root resource ID
	resourcesResult, err := p.apiGWClient.GetResources(ctx, &apigateway.GetResourcesInput{
		RestApiId: aws.String(apiId),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get API Gateway resources: %w", err)
	}

	var rootResourceId string
	for _, item := range resourcesResult.Items {
		if item.Path != nil && *item.Path == "/" {
			rootResourceId = *item.Id
			break
		}
	}

	// Create resources and methods for each HTTP trigger
	for _, resource := range resources {
		if resource.Type == "HttpTrigger" {
			functionName, ok := resource.Properties["function"].(string)
			if !ok {
				continue
			}

			// Create resource
			path := "/"
			if pathProp, ok := resource.Properties["path"].(string); ok {
				path = pathProp
			}

			var resourceId string
			if path == "/" {
				resourceId = rootResourceId
			} else {
				// Create resource
				createResourceResult, err := p.apiGWClient.CreateResource(ctx, &apigateway.CreateResourceInput{
					RestApiId: aws.String(apiId),
					ParentId:  aws.String(rootResourceId),
					PathPart:  aws.String(path),
				})
				if err != nil {
					return "", fmt.Errorf("failed to create API Gateway resource: %w", err)
				}
				resourceId = *createResourceResult.Id
			}

			// Create method
			method := "GET"
			if methodProp, ok := resource.Properties["method"].(string); ok {
				method = methodProp
			}

			_, err = p.apiGWClient.PutMethod(ctx, &apigateway.PutMethodInput{
				RestApiId:         aws.String(apiId),
				ResourceId:        aws.String(resourceId),
				HttpMethod:        aws.String(method),
				AuthorizationType: aws.String("NONE"),
			})
			if err != nil {
				return "", fmt.Errorf("failed to create API Gateway method: %w", err)
			}

			// Create integration
			lambdaArn := fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s",
				p.config.Region, p.getAccountId(), functionName)

			_, err = p.apiGWClient.PutIntegration(ctx, &apigateway.PutIntegrationInput{
				RestApiId:             aws.String(apiId),
				ResourceId:            aws.String(resourceId),
				HttpMethod:            aws.String(method),
				Type:                  apigatewayTypes.IntegrationTypeAwsProxy,
				IntegrationHttpMethod: aws.String("POST"),
				Uri:                   aws.String(fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/%s/invocations", p.config.Region, lambdaArn)),
			})
			if err != nil {
				return "", fmt.Errorf("failed to create API Gateway integration: %w", err)
			}

			// Add permission to Lambda
			_, err = p.lambdaClient.AddPermission(ctx, &lambda.AddPermissionInput{
				FunctionName: aws.String(functionName),
				StatementId:  aws.String(fmt.Sprintf("apigateway-%s", resourceId)),
				Action:       aws.String("lambda:InvokeFunction"),
				Principal:    aws.String("apigateway.amazonaws.com"),
				SourceArn:    aws.String(fmt.Sprintf("arn:aws:execute-api:%s:%s:%s/*/%s%s", p.config.Region, p.getAccountId(), apiId, method, path)),
			})
			if err != nil {
				return "", fmt.Errorf("failed to add Lambda permission: %w", err)
			}
		}
	}

	// Deploy API
	stageName := "prod"
	_, err = p.apiGWClient.CreateDeployment(ctx, &apigateway.CreateDeploymentInput{
		RestApiId: aws.String(apiId),
		StageName: aws.String(stageName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to deploy API Gateway: %w", err)
	}

	// Return API endpoint
	return fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/%s", apiId, p.config.Region, stageName), nil
}

func (p *AWSProvider) deleteFunctions(ctx context.Context, name string) error {
	// Delete Lambda function
	_, err := p.lambdaClient.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
		FunctionName: aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("failed to delete function: %w", err)
	}
	return nil
}

func (p *AWSProvider) deleteAPIGateway(ctx context.Context, name string) error {
	// List APIs
	apisResult, err := p.apiGWClient.GetRestApis(ctx, &apigateway.GetRestApisInput{})
	if err != nil {
		return fmt.Errorf("failed to list API Gateways: %w", err)
	}

	// Find API by name
	for _, item := range apisResult.Items {
		if item.Name != nil && *item.Name == name {
			// Delete API
			_, err = p.apiGWClient.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
				RestApiId: item.Id,
			})
			if err != nil {
				return fmt.Errorf("failed to delete API Gateway: %w", err)
			}
			break
		}
	}
	return nil
}

func (p *AWSProvider) getFunctionStatus(ctx context.Context, name string) ([]ResourceStatus, error) {
	result := make([]ResourceStatus, 0)

	// Get function
	functionResult, err := p.lambdaClient.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(name),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get function: %w", err)
	}

	// Create status
	status := ResourceStatus{
		Name:  name,
		Type:  "Function",
		State: string(functionResult.Configuration.State),
	}

	// LastUpdateStatus is already a string type, not a pointer
	if functionResult.Configuration.LastUpdateStatus != "" {
		status.State = string(functionResult.Configuration.LastUpdateStatus)
	}

	result = append(result, status)
	return result, nil
}

func (p *AWSProvider) getAPIGatewayStatus(ctx context.Context, name string) (*ResourceStatus, error) {
	// List APIs
	apisResult, err := p.apiGWClient.GetRestApis(ctx, &apigateway.GetRestApisInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list API Gateways: %w", err)
	}

	// Find API by name
	for _, item := range apisResult.Items {
		if item.Name != nil && *item.Name == name {
			status := &ResourceStatus{
				Name:  name,
				Type:  "APIGateway",
				State: "ACTIVE",
				URL:   fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/prod", *item.Id, p.config.Region),
			}
			return status, nil
		}
	}

	return nil, nil
}

func (p *AWSProvider) getIAMRole() string {
	// In a real implementation, this would get or create an IAM role
	return fmt.Sprintf("arn:aws:iam::%s:role/lambda-execution-role", p.getAccountId())
}

func (p *AWSProvider) getAccountId() string {
	// In a real implementation, this would get the AWS account ID
	return "123456789012"
}

// Helper function to create a zip file from a directory
func createZipFromDir(sourceDir, zipPath string) error {
	// This is a placeholder - in a real implementation, this would create a zip file
	// containing the contents of the source directory
	return nil
}
