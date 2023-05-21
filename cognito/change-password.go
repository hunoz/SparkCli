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

func changePassword(client cognitoidentityprovider.Client) *cognitoidentityprovider.ChangePasswordOutput {
	var response *cognitoidentityprovider.ChangePasswordOutput
	passwordValidator := CheckIfValidPassword
	if config, e := config.GetCognitoConfig(); e != nil {
		fmt.Printf("Error getting current session: %s\n", e.Error())
		os.Exit(1)
	} else {
		passwordPrompt := promptui.Prompt{
			Label:    "Old Password",
			Mask:     '*',
			Validate: passwordValidator,
		}
		oldPassword, err := passwordPrompt.Run()
		if err != nil {
			fmt.Printf("Error reading old password: %v\n", err.Error())
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
		response, err = client.ChangePassword(context.TODO(), &cognitoidentityprovider.ChangePasswordInput{
			AccessToken:      aws.String(config.AccessToken),
			PreviousPassword: aws.String(oldPassword),
			ProposedPassword: aws.String(newPassword),
		})
		if err != nil {
			fmt.Printf("Error changing password: %v", err.Error())
			os.Exit(1)
		}

		_, err = client.UpdateUserAttributes(context.TODO(), &cognitoidentityprovider.UpdateUserAttributesInput{
			AccessToken: &config.AccessToken,
			UserAttributes: []types.AttributeType{
				{
					Name:  aws.String("custom:passwordChangedAt"),
					Value: aws.String(fmt.Sprint(time.Now().Unix())),
				},
			},
		})
		if err != nil {
			fmt.Printf("Error updating user attributes: %v\n", err.Error())
		}
	}
	return response
}
