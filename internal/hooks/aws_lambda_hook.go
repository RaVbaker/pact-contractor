package hooks

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/viper"
)

type AWSLambdaHook struct {
	AssumeRole     string `mapstructure:"assume_role"`
	FunctionName   string `mapstructure:"function_name"`
	InvocationType string `mapstructure:"invocation_type"` // RequestResponse, Event, DryRun
	Region         string
	Payload        string
}

func (l *AWSLambdaHook) Run(path string) error {
	region := templateString(path, l.Region)
	functionName := templateString(path, l.FunctionName)
	payload := templateString(path, l.Payload)
	invocationType := templateString(path, l.InvocationType)

	if len(functionName) == 0 {
		return fmt.Errorf("cannot find 'function_name' for hook: %v", *l)
	}

	config := &aws.Config{Region: optionalAWSString(region)}
	sess := session.Must(session.NewSession(config))

	pactAssumeRolArn := viper.GetString("aws_assume_role_arn")

	if len(pactAssumeRolArn) > 0 {
		log.Printf("AWS Client with Assumed Role from Config: %q", pactAssumeRolArn)
		credentials := stscreds.NewCredentials(sess, pactAssumeRolArn)
		config = config.WithCredentials(credentials)
		sess = session.Must(session.NewSession(config))
	}

	assumeRole := templateString(path, l.AssumeRole)
	if len(assumeRole) > 0 {
		log.Printf("AWS Client with Assumed Role: %q", assumeRole)
		credentials := stscreds.NewCredentials(sess, assumeRole)
		config = config.WithCredentials(credentials)
	}

	fmt.Printf("%s AWS Lambda %s %q %s invoked\n", assumeRole, region, functionName, invocationType)

	client := lambda.New(sess, config)
	result, err := client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: optionalAWSString(invocationType),
		Payload:        []byte(payload),
	})
	if err != nil {
		return err
	}

	statusCode := *result.StatusCode
	var executedVersion string
	if result.ExecutedVersion != nil {
		executedVersion = *result.ExecutedVersion
	}

	fmt.Printf("Response (status code: %d, version: %q) received\n", statusCode, executedVersion)
	bodyString := string(result.Payload)
	if len(bodyString) > 0 {
		fmt.Println("Body:")
		fmt.Println(bodyString)
	}

	if result.FunctionError != nil && len(*result.FunctionError) > 0 {
		return fmt.Errorf("lambda AWS function error: %q", *result.FunctionError)
	}
	if result.LogResult != nil {
		fmt.Printf("Logs: \n%s\n", *result.LogResult)
	}

	return nil
}

func optionalAWSString(s string) *string {
	if len(s) == 0 {
		return nil
	}
	return aws.String(s)
}
