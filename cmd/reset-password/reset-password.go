package resetpassword

import (
	"os"

	"github.com/fatih/color"
	"github.com/hunoz/spark/cognito"
	"github.com/hunoz/spark/config"
	"github.com/spf13/cobra"
)

var ResetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset your password",
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckIfCognitoIsInitialized()
		var configuration *config.CognitoConfig
		config.CheckIfCognitoIsInitialized()
		if config, e := config.GetCognitoConfig(); e != nil {
			color.Red("Error getting config: %v", e.Error())
			os.Exit(1)
		} else {
			configuration = config
		}
		cognitoClient := cognito.New(configuration)

		cognitoClient.ResetPassword()
	},
}
