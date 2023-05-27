package cognito

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/fatih/color"
)

func respondToAuthChallenge(client cognitoidentityprovider.Client, input cognitoidentityprovider.RespondToAuthChallengeInput) cognitoidentityprovider.RespondToAuthChallengeOutput {
	var response *cognitoidentityprovider.RespondToAuthChallengeOutput
	response, err := client.RespondToAuthChallenge(context.TODO(), &input)
	for {
		if err == nil {
			break
		}

		if strings.Contains(err.Error(), "Your software token has already been used once") {
			color.Red("Your OTP has already been used before, please enter a different OTP")
			otp := getOtp()
			input.ChallengeResponses["SOFTWARE_TOKEN_MFA_CODE"] = otp
			response, err = client.RespondToAuthChallenge(context.TODO(), &input)
		} else {
			color.Red("Error responding to auth challenge: %v", err.Error())
			os.Exit(1)
		}

	}

	return *response
}
