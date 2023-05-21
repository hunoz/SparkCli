package cognito

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/manifoldco/promptui"
)

func forgotPassword(client cognitoidentityprovider.Client, username string) string {
	_, err := client.ForgotPassword(context.TODO(), &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: &ClientId,
		Username: &username,
	})
	if err != nil {
		fmt.Printf("Error resetting password: %v\n", err.Error())
	}

	confirmationCodePrompt := promptui.Prompt{
		Label: "Confirmation Code",
	}
	confirmationCode, err := confirmationCodePrompt.Run()
	if err != nil {
		fmt.Printf("Error reading confirmation code: %v\n", err.Error())
	}

	passwordValidator := CheckIfValidPassword

	passwordPrompt := promptui.Prompt{
		Label:    "New Password",
		Mask:     '*',
		Validate: passwordValidator,
	}
	newPassword, err := passwordPrompt.Run()
	if err != nil {
		fmt.Printf("Error reading new password: %v\n", err.Error())
		os.Exit(1)
	}

	if _, err := client.ConfirmForgotPassword(context.TODO(), &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         &ClientId,
		ConfirmationCode: aws.String(confirmationCode),
		Password:         aws.String(newPassword),
		Username:         aws.String(username),
	}); err != nil {
		fmt.Printf("Error confirming new password: %v\n", err.Error())
		os.Exit(1)
	}

	_, err = client.UpdateUserAttributes(context.TODO(), &cognitoidentityprovider.UpdateUserAttributesInput{
		AccessToken: callCognitoInitiateAuth(client, cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: types.AuthFlowTypeUserPasswordAuth,
			AuthParameters: map[string]string{
				"USERNAME": username,
				"PASSWORD": newPassword,
			},
			ClientId: aws.String(ClientId),
		}, false).AuthenticationResult.AccessToken,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("custom:passwordChangedAt"),
				Value: aws.String(fmt.Sprint(time.Now().Unix())),
			},
		},
	})
	if err != nil {
		fmt.Printf("Error updating user attributes: %v\n", err.Error())
		os.Exit(1)
	}

	return newPassword
}
