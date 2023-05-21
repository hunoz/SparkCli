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
	"gtech.dev/spark/config"
)

func getEmailAttributeCode(client cognitoidentityprovider.Client, accessToken *string) {
	_, err := client.GetUserAttributeVerificationCode(context.TODO(), &cognitoidentityprovider.GetUserAttributeVerificationCodeInput{
		AccessToken:   accessToken,
		AttributeName: aws.String("email"),
	})
	if err != nil {
		fmt.Printf("Error getting email code: %v", err.Error())
		os.Exit(1)
	}
}

func verifyEmail(client cognitoidentityprovider.Client, accessToken *string, attributeCode *string) {
	_, err := client.VerifyUserAttribute(context.TODO(), &cognitoidentityprovider.VerifyUserAttributeInput{
		AccessToken:   accessToken,
		AttributeName: aws.String("email"),
		Code:          attributeCode,
	})
	if err != nil {
		fmt.Printf("Error verifying email: %v", err.Error())
		os.Exit(1)
	}
}

func performFirstSignIn(client cognitoidentityprovider.Client) {
	passwordValidator := CheckIfValidPassword

	usernamePrompt := promptui.Prompt{
		Label: "Username",
	}
	username, err := usernamePrompt.Run()
	if err != nil {
		fmt.Printf("Error reading username: %v\n", err.Error())
	}

	passwordPrompt := promptui.Prompt{
		Label:    "Temporary Password",
		Mask:     '*',
		Validate: passwordValidator,
	}
	temporaryPassword, err := passwordPrompt.Run()
	if err != nil {
		fmt.Printf("Error reading temporary password: %v\n", err.Error())
	}

	passwordPrompt = promptui.Prompt{
		Label:    "New Password",
		Mask:     '*',
		Validate: passwordValidator,
	}
	newPassword, err := passwordPrompt.Run()
	if err != nil {
		fmt.Printf("Error reading new password: %v\n", err.Error())
		os.Exit(1)
	}

	input := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": temporaryPassword,
		},
		ClientId: aws.String(ClientId),
	}

	session := callCognitoInitiateAuth(client, input, true).Session

	respondToAuthChallenge(client, cognitoidentityprovider.RespondToAuthChallengeInput{
		ClientId:      &ClientId,
		ChallengeName: types.ChallengeNameTypeNewPasswordRequired,
		ChallengeResponses: map[string]string{
			"USERNAME":     username,
			"NEW_PASSWORD": newPassword,
			"userAttributes.custom:passwordChangedAt": fmt.Sprint(time.Now().Unix()),
		},
		Session: session,
	})

	initiateAuth(client, username, newPassword, true)

	if config, e := config.GetCognitoConfig(); e != nil {
		fmt.Printf("Error getting current session: %s\n", e.Error())
		os.Exit(1)
	} else {
		fmt.Println("Performing email verification")
		getEmailAttributeCode(client, &config.AccessToken)

		tokenPrompt := promptui.Prompt{
			Label: "Token",
		}
		token, err := tokenPrompt.Run()
		if err != nil {
			fmt.Printf("Error reading new password: %v\n", err.Error())
			os.Exit(1)
		}

		verifyEmail(client, &config.AccessToken, &token)
	}
}
