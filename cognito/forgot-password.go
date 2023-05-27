package cognito

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"gtech.dev/spark/config"
)

func forgotPassword(client cognitoidentityprovider.Client, configuration *config.CognitoConfig, username string) string {
	_, err := client.ForgotPassword(context.TODO(), &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: &configuration.ClientId,
		Username: &username,
	})
	if err != nil {
		color.Red("Error resetting password: %v", err.Error())
	}

	confirmationCodePrompt := promptui.Prompt{
		Label: "Confirmation Code",
	}
	confirmationCode, err := confirmationCodePrompt.Run()
	if err != nil {
		color.Red("Error reading confirmation code: %v", err.Error())
	}

	passwordValidator := CheckIfValidPassword

	passwordPrompt := promptui.Prompt{
		Label:    "New Password",
		Mask:     '*',
		Validate: passwordValidator,
	}
	newPassword, err := passwordPrompt.Run()
	if err != nil {
		color.Red("Error reading new password: %v", err.Error())
		os.Exit(1)
	}

	if _, err := client.ConfirmForgotPassword(context.TODO(), &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         &configuration.ClientId,
		ConfirmationCode: aws.String(confirmationCode),
		Password:         aws.String(newPassword),
		Username:         aws.String(username),
	}); err != nil {
		color.Red("Error confirming new password: %v", err.Error())
		os.Exit(1)
	}

	_, err = client.UpdateUserAttributes(context.TODO(), &cognitoidentityprovider.UpdateUserAttributesInput{
		AccessToken: callCognitoInitiateAuth(client, cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: types.AuthFlowTypeUserPasswordAuth,
			AuthParameters: map[string]string{
				"USERNAME": username,
				"PASSWORD": newPassword,
			},
			ClientId: aws.String(configuration.ClientId),
		}, false).AuthenticationResult.AccessToken,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("custom:passwordChangedAt"),
				Value: aws.String(fmt.Sprint(time.Now().Unix())),
			},
		},
	})
	if err != nil {
		color.Red("Error updating user attributes: %v", err.Error())
		os.Exit(1)
	}

	return newPassword
}
