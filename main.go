package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	ctx = context.TODO()
)

func main() {
	//sourceProfile := flag.String("profile", "", "AWS Profile to use")
	secret := flag.String("secret-name", "", "Secret To Fetch")
	region := flag.String("aws-region", "", "AWS Region where to send requests to")
	version := flag.String("secret-version", "AWSCURRENT", "Version of secret To Fetch")
	//credFile := flag.String("credentials-file", "", "Full path to credentials file")
	flag.Parse()

	if *secret == "" {
		fmt.Printf("You must specify a secret name to fetch by setting the --secret-name CLI flag\n\nHelp:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg, err := config.LoadDefaultConfig(ctx)

	if *region != "" {
		cfg.Region = *region
	}

	secretValue, err := getSecret(cfg, secret, version)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s\n", secretValue)
}

// getCredentialPath returns the users home directory path as a string
func getCredentialPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// Retrieves a secret value from AWS secrets manager
func getSecret(config aws.Config, secretName *string, secretVersionStage *string) (string, error) {
	conn := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(*secretName),
		VersionStage: aws.String(*secretVersionStage),
	}

	result, err := conn.GetSecretValue(ctx, input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}
