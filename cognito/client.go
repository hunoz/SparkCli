package cognito

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/fatih/color"
	"github.com/hunoz/spark/config"
	"github.com/manifoldco/promptui"
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

func (c *CognitoClient) InitiateAuth(force bool) {
	initiateAuth(*c.Client, c.Config, "", "", force)
}

func (c *CognitoClient) RefreshTokens() {
	refreshTokens(*c.Client, c.Config)
}

func (c *CognitoClient) ChangePassword() {
	changePassword(*c.Client, c.Config)
}

func (c *CognitoClient) RegisterMfaDevice() {
	if config, e := config.GetCognitoConfig(); e != nil {
		color.Red("Error getting config: %v", e.Error())
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
		color.Red("Error reading new password: %v", err.Error())
	}
	forgotPassword(*c.Client, c.Config, username)
}

func (c *CognitoClient) PerformFirstSignIn() {
	performFirstSignIn(*c.Client, c.Config)
}
