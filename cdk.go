package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	domainName    = "yourdomainname.com"
	projectPrefix = "yourprojectname"
	stackName     = projectPrefix + "-SpaStack"
	spaBucketName = projectPrefix + "-spa-bucket"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	hostedZone := awsroute53.NewPublicHostedZone(stack, jsii.String("HostedZone"), &awsroute53.PublicHostedZoneProps{
		ZoneName: jsii.String(domainName),
	})

	cert := awscertificatemanager.NewCertificate(stack, jsii.String("Certificate"), &awscertificatemanager.CertificateProps{
		DomainName: jsii.String(domainName),
		Validation: awscertificatemanager.CertificateValidation_FromDns(hostedZone),
	})

	spaBucket := awss3.NewBucket(stack, jsii.String("SpaBucket"), &awss3.BucketProps{
		BucketName:        jsii.String(spaBucketName),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
	})

	oac := awscloudfront.NewS3OriginAccessControl(stack, jsii.String("OAC"), &awscloudfront.S3OriginAccessControlProps{})

	distribution := awscloudfront.NewDistribution(stack, jsii.String("Distribution"), &awscloudfront.DistributionProps{
		DefaultRootObject: jsii.String("index.html"),
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin: awscloudfrontorigins.S3BucketOrigin_WithOriginAccessControl(spaBucket, &awscloudfrontorigins.S3BucketOriginWithOACProps{
				OriginAccessControl: oac,
			}),
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
		DomainNames: &[]*string{jsii.String(domainName)},
		Certificate: cert,
	})

	// aRecord
	_ = awsroute53.NewARecord(stack, jsii.String("ARecord"), &awsroute53.ARecordProps{
		Zone:   hostedZone,
		Target: awsroute53.RecordTarget_FromAlias(awsroute53targets.NewCloudFrontTarget(distribution)),
	})

	return stack
}

func main() {
	defer jsii.Close()
	app := awscdk.NewApp(nil)
	NewCdkStack(app, stackName, &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})
	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
