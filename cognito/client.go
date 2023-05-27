package cognito

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/manifoldco/promptui"
	"gtech.dev/spark/config"
)

type CognitoClient struct {
	Client *cognitoidentityprovider.Client
	Config *config.CognitoConfig
}

func New(configuration *config.CognitoConfig) CognitoClient {
	cognitoIdpClient := cognitoidentityprovider.New(cognitoidentityprovider.Options{
		Region: *aws.String(configuration.Region),
	})
	return CognitoClient{
		Client: cognitoIdpClient,
		Config: configuration,
	}
}

func (c *CognitoClient) InitiateAuth(username string, password string, force bool) {
	initiateAuth(*c.Client, c.Config, username, password, force)
}

func (c *CognitoClient) ChangePassword() {
	changePassword(*c.Client, c.Config)
}

func (c *CognitoClient) RegisterMfaDevice() {
	if config, e := config.GetCognitoConfig(); e != nil {
		fmt.Printf("Error getting config: %s\n", e.Error())
		os.Exit(1)
	} else {
		registerTotp(*c.Client, config.Session)
	}
}

func (c *CognitoClient) ResetPassword() {
	usernamePrompt := promptui.Prompt{
		Label: "Username",
		Mask:  '*',
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		fmt.Printf("Error reading new password: %v\n", err.Error())
	}
	forgotPassword(*c.Client, c.Config, username)
}

func (c *CognitoClient) PerformFirstSignIn() {
	performFirstSignIn(*c.Client, c.Config)
}
