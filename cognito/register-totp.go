package cognito

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/manifoldco/promptui"
)

func getOtp() string {
	otpPrompt := promptui.Prompt{
		Label: "OTP",
	}
	otp, err := otpPrompt.Run()
	if err != nil {
		fmt.Printf("Error reading OTP: %s", err)
		os.Exit(1)
	}

	return otp
}

func verifyTotp(client cognitoidentityprovider.Client, input cognitoidentityprovider.VerifySoftwareTokenInput) {
	_, err := client.VerifySoftwareToken(context.TODO(), &input)
	if err != nil {
		fmt.Printf("Error verifying token: %v\n", err.Error())
		os.Exit(1)
	}
}

func registerTotp(client cognitoidentityprovider.Client, session string) {
	response, err := client.AssociateSoftwareToken(context.TODO(), &cognitoidentityprovider.AssociateSoftwareTokenInput{
		Session: &session,
	})
	if err != nil {
		fmt.Printf("Error registering Totp: %v\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("TOTP Verification code: %v\n", *response.SecretCode)
	otp := getOtp()

	verifyTotp(client, cognitoidentityprovider.VerifySoftwareTokenInput{
		Session:  response.Session,
		UserCode: aws.String(otp),
	})
}
