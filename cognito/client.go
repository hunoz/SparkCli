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
	Client cognitoidentityprovider.Client
}

func New() CognitoClient {
	cognitoIdpClient := cognitoidentityprovider.New(cognitoidentityprovider.Options{
		Region: *aws.String(PoolRegion),
	})
	return CognitoClient{
		Client: *cognitoIdpClient,
	}
}

func (c *CognitoClient) InitiateAuth(username string, password string, force bool) {
	initiateAuth(c.Client, username, password, force)
}

func (c *CognitoClient) ChangePassword() {
	changePassword(c.Client)
}

func (c *CognitoClient) RegisterMfaDevice() {
	if config, e := config.GetCognitoConfig(); e != nil {
		fmt.Printf("Error getting current session: %s\n", e.Error())
		os.Exit(1)
	} else {
		registerTotp(c.Client, config.Session)
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
	forgotPassword(c.Client, username)
}

func (c *CognitoClient) PerformFirstSignIn() {
	performFirstSignIn(c.Client)
}
