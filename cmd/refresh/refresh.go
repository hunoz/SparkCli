package refresh

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hunoz/spark/cognito"
	"github.com/hunoz/spark/config"
	"github.com/spf13/cobra"
)

var RefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh Cognito access and ID tokens",
	Run: func(cmd *cobra.Command, args []string) {
		var configuration *config.CognitoConfig
		config.CheckIfCognitoIsInitialized()
		if config, e := config.GetCognitoConfig(); e != nil {
			if strings.Contains(e.Error(), "Invalid region") {
				color.Red("Spark has not been initialized. Please run 'spark init' to initialize Spark.")
			} else {
				color.Red("Error getting config: %v", e.Error())
			}
			os.Exit(1)
		} else {
			configuration = config
		}

		cognitoClient := cognito.New(configuration)

		cognitoClient.RefreshTokens()
	},
}
