package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkProjStackProps struct {
	awscdk.StackProps
}

func NewWebServiceGinTutorialStack(scope constructs.Construct, id string, props *CdkProjStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("Table"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: aws.String("pk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: aws.String("sk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	// awssqs.NewQueue(stack, jsii.String("user-notifier-queue"), &awssqs.QueueProps{})

	// TODO: Break out to make less disgusting?
	albumLambdaRole := awsiam.NewRole(stack, jsii.String("album-lambda-role"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaBasicExecutionRole")),
		},
		InlinePolicies: &map[string]awsiam.PolicyDocument{
			"tableActions": awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
				Statements: &[]awsiam.PolicyStatement{
					awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
						Effect: awsiam.Effect_ALLOW,
						Actions: &[]*string{
							jsii.String("dynamodb:Scan"),
							jsii.String("dynamodb:GetItem"),
							jsii.String("dynamodb:PutItem"),
							jsii.String("dynamodb:DeleteItem"),
						},
						Resources: &[]*string{
							table.TableArn(),
						},
					}),
				},
			}),
		},
	})

	albumLambda := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("albums-handler"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime:     awslambda.Runtime_PROVIDED_AL2023(),
		Environment: &map[string]*string{"TABLE_NAME": table.TableName()},
		Entry:       jsii.String("./web-service"),
		Role:        albumLambdaRole,
	})

	api := awsapigateway.NewLambdaRestApi(stack, jsii.String("web-service-api"), &awsapigateway.LambdaRestApiProps{
		Handler: albumLambda,
	})

	albumsResource := api.Root().AddResource(jsii.String("albums"), &awsapigateway.ResourceOptions{})
	albumsResource.AddProxy(&awsapigateway.ProxyResourceOptions{
		DefaultIntegration: awsapigateway.NewLambdaIntegration(albumLambda, &awsapigateway.LambdaIntegrationOptions{}),
	})

	awscdk.NewCfnOutput(stack, jsii.String("api-gateway-endpoint"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("API-Gateway-Endpoint"),
		Value:      api.Url(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewWebServiceGinTutorialStack(app, "WebServiceGinTutorialStack", &CdkProjStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
