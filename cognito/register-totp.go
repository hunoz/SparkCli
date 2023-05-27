package cognito

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func getOtp() string {
	otpPrompt := promptui.Prompt{
		Label: "OTP",
	}
	otp, err := otpPrompt.Run()
	if err != nil {
		color.Red("Error reading OTP: %v", err.Error())
		os.Exit(1)
	}

	return otp
}

func verifyTotp(client cognitoidentityprovider.Client, input cognitoidentityprovider.VerifySoftwareTokenInput) {
	_, err := client.VerifySoftwareToken(context.TODO(), &input)
	if err != nil {
		color.Red("Error verifying token: %v", err.Error())
		os.Exit(1)
	}
}

func registerTotp(client cognitoidentityprovider.Client, session string) {
	response, err := client.AssociateSoftwareToken(context.TODO(), &cognitoidentityprovider.AssociateSoftwareTokenInput{
		Session: &session,
	})
	if err != nil {
		color.Red("Error registering Totp: %v", err.Error())
		os.Exit(1)
	}

	color.Cyan("Please use the following code code to register your TOTP device: %v", response.SecretCode)
	otp := getOtp()

	verifyTotp(client, cognitoidentityprovider.VerifySoftwareTokenInput{
		Session:  response.Session,
		UserCode: aws.String(otp),
	})
}
