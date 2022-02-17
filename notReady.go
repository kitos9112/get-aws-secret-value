// CreateSession Creates AWS Session with specified profile
func CreateSession(profileName *string) *session.Session {
	profileNameValue := *profileName
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: profileNameValue,
	}))
	return sess
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


