package cloud

type AWSConfig struct {
	Profile     string
	Region      string
	Credentials AWSCredentials
	Services    AWSServices
}

type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	RoleARN        string
	ExternalID     string
}

type AWSServices struct {
	Lambda      AWSLambdaConfig
	APIGateway  AWSAPIGatewayConfig
	CloudWatch  AWSCloudWatchConfig
	S3          AWSS3Config
	DynamoDB    AWSDynamoDBConfig
}

type AWSLambdaConfig struct {
	Runtime      string
	MemorySize   int64
	Timeout      int64
	Handler      string
	Architecture string
}

type AWSAPIGatewayConfig struct {
	Stage        string
	EndpointType string
	CustomDomain string
	Cors        bool
}

type AWSCloudWatchConfig struct {
	LogRetentionDays int
	MetricsEnabled   bool
	AlertsEnabled    bool
}

type AWSS3Config struct {
	BucketPrefix    string
	Versioning      bool
	EncryptionType  string
	LifecycleRules []S3LifecycleRule
}

type AWSDynamoDBConfig struct {
	BillingMode string
	TableClass  string
	Encryption  string
}

type S3LifecycleRule struct {
	Enabled        bool
	ExpirationDays int
	Transitions    []S3Transition
}

type S3Transition struct {
	Days         int
	StorageClass string
}