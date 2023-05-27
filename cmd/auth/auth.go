package auth

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hunoz/spark/cognito"
	"github.com/hunoz/spark/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate to Cognito and cache credentials",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag(FlagKey.Force, cmd.Flags().Lookup(FlagKey.Force))
	},
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

		cognitoClient.InitiateAuth(viper.GetBool(FlagKey.Force))
	},
}

func init() {
	AuthCmd.Flags().BoolP(FlagKey.Force, string(FlagKey.Force[0]), false, "Force update session, even if expiration time is > 6 hours")
}
