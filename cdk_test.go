package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestCdkStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkStack(app, "MyStack", nil)

	// THEN
	template := assertions.Template_FromStack(stack, nil)
	out, err := json.Marshal(template.ToJSON())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Template:")
	fmt.Println(string(out))

	template.HasResourceProperties(jsii.String("AWS::Route53::HostedZone"), map[string]interface{}{
		"Name": "yourdomainname.com.",
	})
	template.HasResourceProperties(jsii.String("AWS::Route53::RecordSet"), map[string]interface{}{
		"Name": "yourdomainname.com.",
		"Type": "A",
	})
	template.HasResourceProperties(jsii.String("AWS::CertificateManager::Certificate"), map[string]interface{}{
		"DomainName": "yourdomainname.com",
	})
	template.HasResourceProperties(jsii.String("AWS::S3::Bucket"), map[string]interface{}{
		"BucketName": "yourprojectname-spa-bucket",
	})
	template.HasResourceProperties(jsii.String("AWS::CloudFront::OriginAccessControl"), map[string]interface{}{
		"OriginAccessControlConfig": map[string]interface{}{
			"OriginAccessControlOriginType": "s3",
			"SigningBehavior":               "always",
			"SigningProtocol":               "sigv4",
		},
	})
	template.HasResourceProperties(jsii.String("AWS::CloudFront::Distribution"), map[string]interface{}{
		"DistributionConfig": map[string]interface{}{
			"Aliases": []interface{}{
				"yourdomainname.com",
			},
			"DefaultRootObject": "index.html",
		},
	})
}
