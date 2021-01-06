package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/alyu/configparser"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

var sess *session.Session
var defaultAwsRegion string = "eu-west-1"

func main() {
	sourceProfile := flag.String("profile", "default", "AWS Profile to use")
	secret := flag.String("secret-name", "secret", "Secret To Fetch")
	region := flag.String("aws-region", "default", "AWS Region where to send requests to")
	version := flag.String("secret-version", "version", "Version of secret To Fetch")
	credFile := flag.String("credentials-file", filepath.Join(getCredentialPath(), ".aws", "credentials"), "Full path to credentials file")
	flag.Parse()

	if *secret == "secret" {
		fmt.Printf("You must specify a secret name to fetch by setting the --secret-name CLI flag\n\nHelp:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	val, regionPresent := os.LookupEnv("AWS_REGION")

	if regionPresent && *region == "default" {
		*region = val
	} else if !regionPresent && *region == "default" {
		*region = defaultAwsRegion
	}

	if *sourceProfile == "default" {
		//Use Default Credentials
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(*region)}))

	} else {
		//Get Specified Credentials
		exists, err := checkProfileExists(credFile, sourceProfile)
		if err != nil || !exists {
			fmt.Println(err.Error())
			return
		}
		sess = CreateSession(sourceProfile)
	}

	getSecret(sess, secret, version)
}

// CreateSession Creates AWS Session with specified profile
func CreateSession(profileName *string) *session.Session {
	profileNameValue := *profileName
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profileNameValue,
	}))
	return sess
}

// getCredentialPath returns the users home directory path as a string
func getCredentialPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// checkProfileExists takes path to the credentials file and profile name to search for
// Returns bool and any errors
func checkProfileExists(credFile *string, profileName *string) (bool, error) {
	config, err := configparser.Read(*credFile)
	if err != nil {
		log.Fatal("Could not find credentials file")
		log.Fatal(err.Error())
		return false, err
	}
	section, err := config.Section(*profileName)
	if err != nil {
		log.Fatal("Could not find profile in credentials file")
		return false, nil
	}
	if !section.Exists("aws_access_key_id") {
		log.Fatal("Could not find access key in profile")
		return false, nil
	}

	return true, nil
}

func getSecret(sess *session.Session, secretName *string, secretVersion *string) {
	svc := secretsmanager.New(sess)
	var versionID string
	if *secretVersion == "version" {
		versionID = "AWSCURRENT"
	} else {
		versionID = *secretVersion
	}
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(*secretName),
		VersionStage: aws.String(versionID),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException)
				log.Fatal("FATAL: Secret name --> ", *secretName, " could not be found in ", *sess.Config.Region, " region")
				os.Exit(1)
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
				log.Fatal("FATAL: Secret name --> ", *secretName, " appears invalid in ", *sess.Config.Region, " region")
				os.Exit(1)
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
				log.Fatal("FATAL: Secret name --> ", *secretName, " cannot be retrieved in ", *sess.Config.Region, " region - Invalid request")
				os.Exit(1)
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
				log.Fatal("FATAL: Secret name --> ", *secretName, " cannot be decrypted in ", *sess.Config.Region, " region")
				os.Exit(1)
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
				log.Fatal("FATAL: Secret name --> ", *secretName, " cannot be retrieved in ", *sess.Config.Region, " region - Internal server error")
				os.Exit(1)
			default:
				fmt.Println(aerr.Error())
				log.Fatal("FATAL: Secret name --> ", *secretName, " cannot be retrieved in ", *sess.Config.Region, " region - Have credentials been passed?")
				os.Exit(1)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Printf("%s\n", *result.SecretString)
}
