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

func changePassword(client cognitoidentityprovider.Client, configuration *config.CognitoConfig) *cognitoidentityprovider.ChangePasswordOutput {
	var response *cognitoidentityprovider.ChangePasswordOutput
	passwordValidator := CheckIfValidPassword
	passwordPrompt := promptui.Prompt{
		Label:    "Old Password",
		Mask:     '*',
		Validate: passwordValidator,
	}
	oldPassword, err := passwordPrompt.Run()
	if err != nil {
		color.Red("Error reading old password: %v", err.Error())
	}

	passwordPrompt = promptui.Prompt{
		Label:    "New Password",
		Mask:     '*',
		Validate: passwordValidator,
	}
	newPassword, err := passwordPrompt.Run()
	if err != nil {
		color.Red("Error reading new password: %v", err.Error())
		os.Exit(1)
	}
	response, err = client.ChangePassword(context.TODO(), &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(configuration.AccessToken),
		PreviousPassword: aws.String(oldPassword),
		ProposedPassword: aws.String(newPassword),
	})
	if err != nil {
		color.Red("Error changing password: %v", err.Error())
		os.Exit(1)
	}

	_, err = client.UpdateUserAttributes(context.TODO(), &cognitoidentityprovider.UpdateUserAttributesInput{
		AccessToken: &configuration.AccessToken,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("custom:passwordChangedAt"),
				Value: aws.String(fmt.Sprint(time.Now().Unix())),
			},
		},
	})
	if err != nil {
		color.Red("Error updating user attributes: %v", err.Error())
	}
	return response
}
