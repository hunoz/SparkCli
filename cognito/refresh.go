package cognito

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/fatih/color"
	"github.com/hunoz/spark/config"
)

func refreshTokens(client cognitoidentityprovider.Client, configuration *config.CognitoConfig) {
	if configuration.RefreshToken == "" {
		color.Red("Refresh token is missing, please run 'spark auth' to authenticate and obtain a refresh token")
		return
	} else if time.Until(time.Unix(configuration.RefreshTokenExpiry, 0)) <= 0 {
		color.Red("Refresh token is expired, please run 'spark auth' to obtain a new refresh token")
		return
	}

	input := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": configuration.RefreshToken,
		},
		ClientId: aws.String(configuration.ClientId),
	}

	response := callCognitoInitiateAuth(client, input, false)

	now := time.Now()

	cognitoConfig := config.CognitoConfig{
		ClientId:           configuration.ClientId,
		Region:             configuration.Region,
		PoolId:             configuration.PoolId,
		AccessToken:        *response.AuthenticationResult.AccessToken,
		IdToken:            *response.AuthenticationResult.IdToken,
		RefreshToken:       configuration.RefreshToken,
		RefreshTokenExpiry: configuration.RefreshTokenExpiry,
		Expires:            now.Add(time.Second * time.Duration(response.AuthenticationResult.ExpiresIn)).Unix(),
		Session:            configuration.Session,
	}

	if err := config.UpdateCognitoConfig(cognitoConfig); err != nil {
		color.Red("Error updating tokens: %v", err.Error())
		os.Exit(1)
	}

	color.Green("Successfully updated tokens")
}
