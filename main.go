package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var sourceProfile string
var secretName string
var secretVersion string
var region string
var credFile string
var debug bool

// initialise initialises the program and sets up the AWS session
func initialiseClient() (*secretsmanager.Client, error) {
	defaultRegion := getEnv("AWS-REGION", "eu-west-1")
	defaultCredFile := getEnv("AWS_SHARED_CREDENTIALS_FILE", filepath.Join(getCredentialPath(), ".aws", "credentials"))

	flag.StringVar(&secretName, "secret-name", "secret", "AWS Secret Name or ARN To fetch.")
	flag.StringVar(&sourceProfile, "profile", "default", "AWS Profile to use.")
	flag.StringVar(&secretVersion, "secret-version", "latest", "Version of the AWS secret To Fetch.")
	flag.StringVar(&credFile, "credentials-file", defaultCredFile, "Full path to credentials file.")
	flag.StringVar(&region, "aws-region", defaultRegion, "AWS Region to use. Falls back to $AWS_REGION if not set.")
	flag.BoolVar(&debug, "v", false, "Enable debug logging")
	flag.Parse()

	if secretName == "secret" {
		return nil, errors.New("Secret name is required")
	}
	if debug {
		log.Println("DEBUG - Secret Name:", secretName)
		log.Println("DEBUG - Source Profile:", sourceProfile)
		log.Println("DEBUG - Secret Version:", secretVersion)
		log.Println("DEBUG - Credentials File:", credFile)
		log.Println("DEBUG - AWS Region:", region)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithSharedCredentialsFiles([]string{credFile}))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil, err
	}

	return secretsmanager.NewFromConfig(cfg), nil
}

// getEnv returns the value of the environment variable named by the key or the fallback value if the variable is not set.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getCredentialPath returns the users home directory path as a string
func getCredentialPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// getSecret gets the secret value from AWS Secrets Manager
func getSecret(client *secretsmanager.Client, secretName *string, secretVersion *string) (string, error) {

	var versionID string
	if *secretVersion == "latest" {
		versionID = "AWSCURRENT"
	} else {
		versionID = *secretVersion
	}
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(*secretName),
		VersionStage: aws.String(versionID),
	}

	result, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", err
		// if aerr, ok := err.(awserr.Error); ok {
		// 	switch aerr.Code() {
		// 	case secretsmanager.ErrCodeResourceNotFoundException:
		// 		fmt.Println(secretsmanager.type.ResourceNotFoundException, aerr.Error())
		// 		log.Fatal("FATAL: Secret name --> ", *secretName, " could not be found in ", *sess.Config.Region, " region")
		// 		os.Exit(1)
		// 	case secretsmanager.ErrCodeInvalidParameterException:
		// 		fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
		// 		log.Fatal("FATAL: Secret name --> ", *secretName, " appears invalid in ", *sess.Config.Region, " region")
		// 		os.Exit(1)
		// 	case secretsmanager.ErrCodeInvalidRequestException:
		// 		fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
		// 		log.Fatal("FATAL: Secret name --> ", *secretName, " cannot be retrieved in ", *sess.Config.Region, " region - Invalid request")
		// 		os.Exit(1)
		// 	case secretsmanager.ErrCodeDecryptionFailure:
		// 		fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
		// 		log.Fatal("FATAL: Secret name --> ", *secretName, " cannot be decrypted in ", *sess.Config.Region, " region")
		// 		os.Exit(1)
		// 	case secretsmanager.ErrCodeInternalServiceError:
		// 		fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
		// 		log.Fatal("FATAL: Secret name --> ", *secretName, " cannot be retrieved in ", *sess.Config.Region, " region - Internal server error")
		// 		os.Exit(1)
		// 	default:
		// 		fmt.Println(aerr.Error())
		// 		log.Fatal("FATAL: Secret name --> ", *secretName, " cannot be retrieved in ", *sess.Config.Region, " region - Have credentials been passed?")
		// 		os.Exit(1)
		// 	}
		// } else {
		// 	// Print the error, cast err to awserr.Error to get the Code and Message from an error.
		// 	fmt.Println(err.Error())
		// }
	}
	return *result.SecretString, nil
}

func main() {

	client, err := initialiseClient()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	secret, err := getSecret(client, &secretName, &secretVersion)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(secret)
}
