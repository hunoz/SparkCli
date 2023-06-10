package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hunoz/spark/aws"
	"github.com/hunoz/spark/homedir"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type CognitoConfig struct {
	Region             string `json:",omitempty" type:"string"`
	ClientId           string `json:",omitempty" type:"string"`
	PoolId             string `json:",omitempty" type:"string"`
	AccessToken        string `json:",omitempty" type:"string"`
	IdToken            string `json:",omitempty" type:"string"`
	RefreshToken       string `json:",omitempty" type:"string"`
	RefreshTokenExpiry int64  `json:",omitempty" type:"integer"`
	Expires            int64  `json:",omitempty" type:"integer"`
	Session            string `json:",omitempty" type:"string"`
}

type Config struct {
	Cognito CognitoConfig `json:",omitempty"`
}

func GetSparkConfigFile() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "unable to find home folder.")
	}

	return filepath.Join(home, ".config", "spark", "config.json"), nil
}

// OpenReadConfigFile opens the session config file with read only permissions
func OpenReadConfigFile() (*os.File, error) {
	configPath, err := GetSparkConfigFile()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get config path")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return nil, errors.Wrap(err, "Cannot create config file")
		}
	}

	return os.OpenFile(configPath, os.O_RDONLY|os.O_CREATE, 0600)
}

// OpenWriteConfigFile opens the session config file with write only permissions
func OpenWriteConfigFile() (*os.File, error) {
	configPath, err := GetSparkConfigFile()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get config path")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return nil, errors.Wrap(err, "Cannot create config file")
		}
	}

	return os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}

// writeSparkConfig overwrites the spark session config file with an updated config
func writeSparkConfig(config *Config) error {
	file, err := OpenWriteConfigFile()
	if err != nil {
		return errors.Wrap(err, "writeSparkConfig failed to open the spark session config file")
	}
	defer file.Close()

	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Failed to marshal the new spark session config")
	}

	_, err = file.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "writeSparkConfig failed to write to the spark session config file")
	}

	return nil
}

// readSparkConfig reads the spark session config and returns a struct containing the data from the file
func readSparkConfig() (*Config, error) {
	file, err := OpenReadConfigFile()
	if err != nil {
		return nil, errors.Wrap(err, "readSparkConfig failed to open the Spark session config file")
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "readSparkConfig unable to retrieve info about Spark session config file")
	}

	configBytes := make([]byte, stat.Size())
	var config Config

	count, err := file.Read(configBytes)
	if err != nil || count < 0 {
		return nil, errors.Wrap(err, "readSparkConfig failed to read the Spark session config file")
	} else if count == 0 {
		return &config, nil
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal Spark session config")
	}

	return &config, nil
}

func GetCognitoConfig() (*CognitoConfig, error) {
	config, err := readSparkConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Could not read Spark session config")
	}

	if err := CognitoConfigIsValid(&config.Cognito); err != nil {
		return nil, err
	}

	return &config.Cognito, nil
}

// UpdateCognitoConfig takes a session as an argument and adds it to the spark session config file
func UpdateCognitoConfig(newConfig CognitoConfig) error {
	config, err := readSparkConfig()
	if err != nil {
		return errors.Wrap(err, "Could not read Spark session config")
	}

	config.Cognito = newConfig

	err = writeSparkConfig(config)
	if err != nil {
		return errors.Wrap(err, "Could not write to Spark session config")
	}

	return nil
}

func CognitoIsInitialized() (*CognitoConfig, error) {
	config, err := readSparkConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Could not read Spark session config")
	}

	if config.Cognito.AccessToken != "" || config.Cognito.Expires == 0 || config.Cognito.IdToken != "" {
		return &config.Cognito, nil
	}

	if config.Cognito.ClientId == "" || config.Cognito.Region == "" || config.Cognito.PoolId == "" {
		return nil, errors.New("Spark has not been initialized. Please run 'spark init' to initialize Spark.")
	}

	return &config.Cognito, nil
}

func CheckIfCognitoIsInitialized() {
	if _, err := CognitoIsInitialized(); err != nil {
		color.Red("%v", err.Error())
		os.Exit(1)
	}
}

func CognitoConfigIsValid(config *CognitoConfig) error {
	if err := IsValidAwsRegion(config.Region); err != nil {
		return err
	}
	return nil
}

func IsValidAwsRegion(region string) error {
	if regions, err := aws.GetAwsRegions(); err != nil {
		return errors.Wrap(err, "Error fetching AWS regions")
	} else {
		if slices.Contains(regions, region) {
			return nil
		}
		return errors.New(fmt.Sprintf("Invalid region: '%v'", region))
	}
}
